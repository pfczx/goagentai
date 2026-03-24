package memory

import "github.com/pfczx/goagentai/llm"

type MemoryMenager struct {
	path                                string
	profile                             string
	MemoryOn                            bool
	ShortTermMemoryLimit                int
	ShortTermMemoryEvaluation           bool
	ShortTermMemory                     *ShortTermMemory
	LongTermMemoryBufferSize            int
	LongTermMemoryChunkSize             int
	LongTermMemoryStorageSize           int
	LongTermMemoryChunksToAdd           int
	LongTermMemorySummarizationProvider llm.ModelProvider
	LongTermMemory                      *LongTermMemory
}

func InitMenager(path string, profile string, memoryOn bool,
	shortTermMemoryLimit int, shortTermMemoryEvaluation bool,
	longTermMemoryBufferSize int, longTermMemoryChunkSize int, longTermMemoryStorageSize int, longTermMemoryChunkToAdd int,
	longTermMemorySummarizationProvider string, longTermMemorySummarizationInternalProvider string, longTermMemorySummarizationModel string,
) (*MemoryMenager, error) {

	shortMemory, err := LoadShortTermMemory(path)
	if err != nil {
		return nil, err
	}
	summarizationProvider, err := llm.NewProvider(longTermMemorySummarizationProvider,
		longTermMemorySummarizationModel,
		longTermMemorySummarizationInternalProvider)
	if err != nil {
		return nil, err
	}
	longTermMemory, err := NewLongTermMemory(profile)
	if err != nil {
		return nil, err
	}

	return &MemoryMenager{
		path:                                path,
		profile:                             profile,
		MemoryOn:                            memoryOn,
		ShortTermMemoryLimit:                shortTermMemoryLimit,
		ShortTermMemoryEvaluation:           shortTermMemoryEvaluation,
		ShortTermMemory:                     shortMemory,
		LongTermMemoryBufferSize:            longTermMemoryBufferSize,
		LongTermMemoryChunkSize:             longTermMemoryChunkSize,
		LongTermMemoryStorageSize:           longTermMemoryStorageSize,
		LongTermMemoryChunksToAdd:           longTermMemoryChunkToAdd,
		LongTermMemorySummarizationProvider: summarizationProvider,
		LongTermMemory:                      longTermMemory,
	}, nil
}
