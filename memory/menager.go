package memory

import ()

type MemoryMenager struct {
	path                      string
	MemoryOn                  bool
	ShortTermMemoryLimit      int
	ShortTermMemoryEvaluation bool
	ShortTermMemory           *ShortTermMemory
}

func InitMenager(path string, memoryOn bool, shortTermMemoryLimit int, shortTermMemoryEvaluation bool,
	shortTermMemory *ShortTermMemory) (*MemoryMenager, error) {

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
