package memory

import (
	"fmt"
	"math"
	"time"

	"github.com/pfczx/goagentai/prompt"
	"github.com/pfczx/goagentai/token"
)

type LongTermMemory struct {
	handler  *MemoryHandler
	analyzer *TextAnalyzer
	idf      map[string]float32
	docCount int
}

type MemoryChunk struct {
	Profile   string
	Summary   string
	TF        map[string]float32
	CreatedAt time.Time
}

func (m *MemoryMenager) CloseDb() error {
	return m.LongTermMemory.handler.CloseDB()
}

func NewLongTermMemory(profileID string) (*LongTermMemory, error) {
	h, err := NewMemoryHandler()
	if err != nil {
		return nil, err
	}
	a := NewTextAnalyzer()
	longTermMemory := &LongTermMemory{
		handler:  h,
		analyzer: a,
	}
	allChunks, err := h.GetLongTerm(profileID)
	if err != nil {
		return nil, err
	}
	longTermMemory.buildIDF(allChunks)

	return longTermMemory, nil

}

func (m *LongTermMemory) buildIDF(chunks []MemoryChunk) {
	df := make(map[string]int)
	N := float32(len(chunks))

	for _, c := range chunks {
		seen := make(map[string]bool)
		for w := range c.TF {
			if !seen[w] {
				df[w]++
				seen[w] = true
			}
		}
	}

	idf := make(map[string]float32)
	for w, count := range df {
		idf[w] = float32(math.Log(1 + float64(N/(1+float32(count)))))
	}

	m.idf = idf
	m.docCount = int(N)
}

func (m *MemoryMenager) SaveShortTermToBuffer(parts []ShortTermPart) error {
	for num, part := range m.ShortTermMemory.Content {
		stringPart := fmt.Sprintf("%d User: %s Agent: %s ", num, part.Prompt, part.Response)
		if m.ShortTermMemoryEvaluation {
			stringPart = stringPart + fmt.Sprintf("Usefull: %t  ", part.Usefull)
		}
		err := m.LongTermMemory.handler.SaveShortTerm(m.profile, stringPart)
		if err != nil {
			return err
		}

	}
	return nil
}

func (m *MemoryMenager) TriggerSummarization() (bool, error) {
	shortTermSize, err := m.LongTermMemory.handler.CountShortTerm(m.profile)
	if err != nil {
		return false, err
	}
	if shortTermSize > m.LongTermMemoryBufferSize {
		return true, nil
	}
	return false, nil
}

func (m *MemoryMenager) Summarize() (string, int, error) {
	shortTermBuffer, err := m.LongTermMemory.handler.GetShortTerm(m.profile)
	if err != nil {
		return "", 0, err
	}
	var parts []string
	for _, part := range shortTermBuffer {
		parts = append(parts, part.Memory)
	}
	prompt, err := prompt.BuildSummarize(m.LongTermMemoryChunkSize, parts)
	if err != nil {
		return "", 0, err
	}
	resp, err := m.LongTermMemorySummarizationProvider.Generate(prompt)
	if err != nil {
		return "", 0, err
	}
	if resp.Usage != nil {
		return resp.Text, resp.Usage.TotalTokens, nil
	}

	return "", 0, fmt.Errorf("no token usage error : summarization")

}

func (m *MemoryMenager) UpdateLongTerm(tokenMenager *token.TokenMenager) error {
	t, err := m.TriggerSummarization()
	if err != nil {
		return err
	}
	if t {
		fmt.Println("Summarizing memory in background...")
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("panic in summary:", r)
				}
			}()

			summarization, used, err := m.Summarize()
			if err != nil {
				fmt.Println("summary error:", err)
				return
			}

			tf := m.LongTermMemory.analyzer.ComputeTF(summarization)

			m.LongTermMemory.handler.SaveLongTerm(m.profile, summarization, tf)

			if err := m.LongTermMemory.handler.ClearShortTerm(m.profile); err != nil {
				fmt.Println("clear short term error:", err)
				return
			}

			if err := m.LongTermMemory.handler.TrimLongTermToStorageSize(
				m.profile,
				m.LongTermMemoryStorageSize,
			); err != nil {
				fmt.Println("trim long term error:", err)
				return
			}
			chunks, err := m.LongTermMemory.handler.GetLongTerm(m.profile)
			if err != nil {
				fmt.Println(err)
				return
			}
			m.LongTermMemory.buildIDF(chunks)
			if used > 0 && tokenMenager != nil {
				tokenMenager.AddUsage(used)
			}
		}()
	}
	return nil
}

func (m *MemoryMenager) GetLongTermContextString(input string) (string, error) {
	memory, err := m.LongTermMemory.handler.GetLongTerm(m.profile)
	if err != nil {
		return "", err
	}
	selected := m.LongTermMemory.analyzer.SelectRelevantChunks(input, memory, m.LongTermMemory.idf, m.LongTermMemoryChunksToAdd)
	out := fmt.Sprint("This is long term chunks of summarized user history for additional context: ")
	if len(selected) == 0 {
		return "", nil
	}
	for num, memo := range selected {
		out = out + fmt.Sprintf("%d %s ", num, memo.Summary)
	}
	return out, nil

}

func (m *MemoryMenager) ClearLongTermMemory(profileID string) error {
	err := m.LongTermMemory.handler.ClearShortTerm(profileID)
	if err != nil {
		return err
	}
	return m.LongTermMemory.handler.ClearLongTerm(profileID)
}
