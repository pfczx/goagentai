package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
)

type ShortTermMemory struct {
	Limit   int             `json:"short_term_history_size"`
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
		if _, err = os.Create(path); err != nil {
			return err
		}
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

func (m *MemoryMenager) AppendShortTermHistory(prompt string, response string, usefull ...bool) error {
	part := &ShortTermPart{
		Prompt:   prompt,
		Response: response,
	}
	if len(usefull) > 0 {
		part.Usefull = usefull[0]
	}
	m.ShortTermMemory.Content = append(m.ShortTermMemory.Content, *part)

	if len(m.ShortTermMemory.Content) > m.ShortTermMemory.Limit {
		//in future, send deleted to long term buffer and summarize
		m.ShortTermMemory.Content = slices.Delete(m.ShortTermMemory.Content, 0, 1)
	}

	memoryPath := filepath.Join(m.path, "shortTermMemory.json")
	err := SaveShortTermMemory(memoryPath, m.ShortTermMemory)
	if err != nil {
		return err
	}
	return nil
}

func (m *MemoryMenager) ShortTermMemoryString()  string {
	var out 

} 
