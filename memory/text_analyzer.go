package memory

import (
	"math"
	"sort"
	"strings"
	"unicode"

)

type TextAnalyzer struct {
	stopwords       map[string]bool
	keywordTreshold float32
}

func NewTextAnalyzer() *TextAnalyzer {
	return &TextAnalyzer{
		stopwords: Stopwords,
	}
}

func (t *TextAnalyzer) removeStopwords(words []string) []string {
	var filtered []string

	for _, w := range words {
		if !t.stopwords[w] {
			filtered = append(filtered, w)
		}
	}

	return filtered
}
func (t *TextAnalyzer) tokenize(text string) []string {
	text = strings.ToLower(text)

	var cleaned strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsSpace(r) {
			cleaned.WriteRune(r)
		}
	}

	return t.removeStopwords(strings.Fields(cleaned.String()))
}

func (t *TextAnalyzer) ComputeTF(text string) map[string]float32 {
	words := t.tokenize(text)

	tf := make(map[string]float32)
	total := float32(len(words))

	if total == 0 {
		return tf
	}

	for _, w := range words {
		tf[w]++
	}

	for k, v := range tf {
		tf[k] = v / total
	}

	return tf
}
func (t *TextAnalyzer) CosineSimilarity(a, b map[string]float32) float32 {
	var dot, magA, magB float32
	for k, v := range a {
		if bVal, ok := b[k]; ok {
			dot += v * bVal
		}
	}
	for _, v := range a {
		magA += v * v
	}
	for _, v := range b {
		magB += v * v
	}
	if magA == 0 || magB == 0 {
		return 0
	}
	return dot / (float32(math.Sqrt(float64(magA))) * float32(math.Sqrt(float64(magB))))
}



func (t *TextAnalyzer) SelectRelevantChunks(query string, chunks []MemoryChunk, idf map[string]float32, topN int) []MemoryChunk {
	queryTF := t.ComputeTF(query)
	queryTFIDF := make(map[string]float32)
	for w, tfVal := range queryTF {
		if idfVal, ok := idf[w]; ok {
			queryTFIDF[w] = tfVal * idfVal
		} else {
			queryTFIDF[w] = tfVal
		}
	}

	type scoredChunk struct {
		chunk MemoryChunk
		score float32
	}

	var scored []scoredChunk
	for _, c := range chunks {
		//only words from query
		score := float32(0)
		for w, qVal := range queryTFIDF {
			if cVal, ok := c.TF[w]; ok {
				score += qVal * cVal
			}
		}
		scored = append(scored, scoredChunk{chunk: c, score: score})
	}


	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	var result []MemoryChunk
	for i := 0; i < topN && i < len(scored); i++ {
		result = append(result, scored[i].chunk)
	}
	return result
}
