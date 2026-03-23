package memory

import (
	"github.com/duynguyendang/meb/vector"
	"github.com/securisec/go-keywords"
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
		//low treshold tf values will be used as first filter when appending context to ask
		keywordTreshold: 0.1,
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

func (t *TextAnalyzer) ExtractKeywords(text string) []string {
	opts := keywords.ExtractOptions{
		IgnoreLength: 0,
	}
	keywords, _ := keywords.Extract(text, opts)

	return keywords
}

func tfMapsToVectorsShared(queryTF, chunkTF map[string]float32) ([]float32, []float32) {
	vecQuery := []float32{}
	vecChunk := []float32{}

	vocab := make(map[string]struct{})
	for w := range queryTF {
		vocab[w] = struct{}{}
	}
	for w := range chunkTF {
		vocab[w] = struct{}{}
	}
	for w := range vocab {
		vecQuery = append(vecQuery, queryTF[w])
		vecChunk = append(vecChunk, chunkTF[w])
	}
	return vecQuery, vecChunk
}

func (t *TextAnalyzer) SelectRelevantChunks(input string, chunks []MemoryChunk, topN int) []MemoryChunk {
	queryTF := t.ComputeTF(input)
	queryKeywords := t.ExtractKeywords(input)

	type scoredChunk struct {
		chunk MemoryChunk
		score float32
	}
	//skip chunk if similar keyword not detected
	var scored []scoredChunk
	queryKWMap := make(map[string]struct{}, len(queryKeywords))
	for _, kw := range queryKeywords {
		queryKWMap[kw] = struct{}{}
	}

	for _, c := range chunks {
		overlap := false
		for _, kw := range c.Keywords {
			if _, ok := queryKWMap[kw]; ok {
				overlap = true
				break
			}
		}
		if !overlap {
			continue
		}
		queryVector, chunkVector := tfMapsToVectorsShared(queryTF, c.TF)
		score := vector.CosineSimilarity(queryVector, chunkVector)
		scored = append(scored, scoredChunk{chunk: c, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	result := []MemoryChunk{}
	for i := 0; i < len(scored) && i < topN; i++ {
		result = append(result, scored[i].chunk)
	}

	return result
}
