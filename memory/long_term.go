package memory

import (
	"fmt"
	"strings"
	"time"

	"github.com/pfczx/goagentai/prompt"
)

type LongTermMemory struct {
	handler  *MemoryHandler
	analyzer *TextAnalyzer
}

type MemoryChunk struct {
	Profile   string
	Summary   string
	Keywords  []string
	TF        map[string]float32
	CreatedAt time.Time
}

func NewLongTermMemory() (*LongTermMemory, error) {
	h, err := NewMemoryHandler()
	if err != nil {
		return nil, err
	}
	a := NewTextAnalyzer()
	return &LongTermMemory{
		handler:  h,
		analyzer: a,
	}, nil

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

func (m *MemoryMenager) Summarize() (string, error) {
	shortTermBuffer, err := m.LongTermMemory.handler.GetShortTerm(m.profile)
	if err != nil {
		return "", err
	}
	var parts []string
	for _, part := range shortTermBuffer {
		parts = append(parts, part.Memory)
	}
	prompt, err := prompt.BuildSummarize(m.LongTermMemoryChunkSize, parts)
	if err != nil {
		return "", err
	}
	resp, err := m.LongTermMemorySummarizationProvider.Generate(prompt)
	if err != nil {
		return "", err
	}
	if resp.Usage != nil {
		return resp.Text, nil
	}

	return "", fmt.Errorf("no token usage error : summarization")

}

func (m *MemoryMenager) UpdateLongTerm() error {
	t, err := m.TriggerSummarization()
	if err != nil {
		return err
	}
	if t {
		summarization, err := m.Summarize()
		if err != nil {
			return err
		}
		tf := m.LongTermMemory.analyzer.ComputeTF(summarization)
		keywords := m.LongTermMemory.analyzer.ExtractKeywords(tf)
		m.LongTermMemory.handler.SaveLongTerm(m.profile, summarization, tf, keywords)

	}

	return nil
}

func (m *MemoryMenager) GetLongTermContextString(input string) (string, error) {
	memory, err := m.LongTermMemory.handler.GetLongTerm(m.profile)
	if err != nil {
		return "", err
	}
	selected := m.LongTermMemory.analyzer.SelectRelevantChunks(input, memory, m.LongTermMemoryChunksToAdd)
	out := fmt.Sprint("This is long term user history for additional context: ")
	for _, memo := range selected {
		out = out + fmt.Sprintf(" %s ", memo.Summary)
	}
	return out, nil

}
