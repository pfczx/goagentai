package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/writeas/go-strip-markdown"
)

type ShortTermMemory struct {
	Content []ShortTermPart `json:"conversation"`
}

type ShortTermPart struct {
	Prompt   string `json:"prompt"`
	Response string `json:"llm_response"`
	Usefull  bool   `json:"usefull,omitempty"`
}

func InitShortMemoryFile(path string) error {
	path = filepath.Join(path, "shortTermMemory.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		memory := ShortTermMemory{
			Content: []ShortTermPart{},
		}

		data, err := json.MarshalIndent(memory, "", " ")
		if err != nil {
			return err
		}

		return os.WriteFile(path, data, 0644)
	}

	return nil
}

func LoadShortTermMemory(path string) (*ShortTermMemory, error) {
	memoryPath := filepath.Join(path, "shortTermMemory.json")
	data, err := os.ReadFile(memoryPath)
	if err != nil {
		return nil, err
	}
	var memory ShortTermMemory
	err = json.Unmarshal(data, &memory)
	if err != nil {
		return nil, err
	}
	return &memory, err
}

func SaveShortTermMemory(path string, memory *ShortTermMemory) error {
	memoryPath := filepath.Join(path, "shortTermMemory.json")
	data, err := json.MarshalIndent(memory, "", " ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(memoryPath, data, 0644); err != nil {
		return err
	}

	return nil

}

func (m *MemoryMenager) ShortTermToString() string {
	var out strings.Builder
	info := fmt.Sprintf("this is last conversations for additional context: ")
	out.WriteString(info)

	for _, part := range m.ShortTermMemory.Content {
		stringPart := fmt.Sprintf("User: %s Agent: %s ", part.Prompt, part.Response)
		if m.ShortTermMemoryEvaluation {
			stringPart = stringPart + fmt.Sprintf("Usefull: %t", part.Usefull)
		}
		out.WriteString(stringPart)
	}
	return out.String()

}

func (m *MemoryMenager) AppendShortTermHistory(prompt string, response string, usefull bool) error {
	responseCleared := func(text string) string {
		text = stripmd.Strip(text)
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "|", " ")
		text = strings.ReplaceAll(text,"-"," ")
		return strings.TrimSpace(text)

	}(response)
	part := ShortTermPart{
		Prompt:   prompt,
		Response: responseCleared,
	}
	if m.ShortTermMemoryEvaluation {
		part.Usefull = usefull
	}
	m.ShortTermMemory.Content = append(m.ShortTermMemory.Content, part)

	if len(m.ShortTermMemory.Content) > m.ShortTermMemoryLimit {
		//in future, send deleted to long term buffer and summarize
		//remove oldest messages
		m.ShortTermMemory.Content = slices.Delete(m.ShortTermMemory.Content, 0, len(m.ShortTermMemory.Content)-m.ShortTermMemoryLimit)
	}

	err := SaveShortTermMemory(m.path, m.ShortTermMemory)
	if err != nil {
		return err
	}
	return nil
}
