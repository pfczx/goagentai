package memory

import ()

type MemoryMenager struct {
	MemoryOn                  bool
	ShortTermMemoryLimit      int
	ShortTermMemoryEvaluation bool
	ShortMemory               ShortTermMemory
}

func InitMenager(path string) (*MemoryMenager, error) {
	return &MemoryMenager{
		MemoryOn: true,
		ShortTermMemoryLimit: 20,
		ShortTermMemoryEvaluation: true,
		ShortMemory: ShortTermMemory{},
	},nil

}
