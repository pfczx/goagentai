package memory

import (
	"fmt"
	"time"

	"github.com/pfczx/goagentai/prompt"
)

type LongTermMemory struct {
	handler *MemoryHandler
}

type MemoryChunk struct {
	Profile   string
	Summary   string
	Embedding []float32
	Keywords  []string
	TFIDF     []float32
	CreatedAt time.Time
}

func NewLongTermMemory() (*LongTermMemory, error) {
	handler, err := NewMemoryHandler()
	if err != nil {
		return nil, err
	}
	return &LongTermMemory{
		handler: handler,
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

func (m *MemoryMenager) Summarize() (string,error){
	prompt,err := prompt.BuildSummarize()
	if err !=nil{
		return err
	}

}


func (m *MemoryMenager) UpdateLongTerm() (error) {
	if TriggerSummarization(){
	
	}
	return nil
}

