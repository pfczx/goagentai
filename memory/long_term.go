package memory

import (
	"time"
)

type LongTermMemory struct {
	Buffer  string
	LongTermm string
	Embedder           func(text string) []float32
}

type MemoryChunk struct {
	ID        string
	Summary   string
	Embedding []float32
	Keywords  []string
	TFIDF     []float32
	CreatedAt time.Time
}


