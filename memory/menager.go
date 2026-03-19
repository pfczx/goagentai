package memory

import "github.com/pfczx/goagentai/llm"

type MemoryMenager struct {
	path                                string
	MemoryOn                            bool
	ShortTermMemoryLimit                int
	ShortTermMemoryEvaluation           bool
	ShortTermMemory                     *ShortTermMemory
	LongTermMemoryBufferSize            int
	LongTermMemoryChunkSize             int
	LongTermMemoryStorageSize           int
	LongTermMemorySummarizationProvider llm.ModelProvider
	LongTermMemory                      *LongTermMemory
}

func InitMenager(path string, memoryOn bool, shortTermMemoryLimit int, shortTermMemoryEvaluation bool) (*MemoryMenager, error) {

	shortMemory, err := LoadShortTermMemory(path)
	if err != nil {
		return nil, err
	}
	return &MemoryMenager{
		path:                      path,
		MemoryOn:                  memoryOn,
		ShortTermMemoryLimit:      shortTermMemoryLimit,
		ShortTermMemoryEvaluation: shortTermMemoryEvaluation,
		ShortTermMemory:           shortMemory,
	}, nil

}
