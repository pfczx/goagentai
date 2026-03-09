package memory

import (
	"os"
	"path/filepath"
)

type ShortTermMemory struct {
	Limit   int               `json:"short_term_history_size"`
	Content []ShortTermMemory `json:"conversation"`
}

type ShortTermPart struct {
	Author  string `json:"author"`
	Text    string `json:"text"`
	Usefull bool   `json:"usefull,omitempty"`
}

func InitShortMemoryFile(path string) error {
	path = filepath.Join(path, "shortTermMemory.json")
	if _, err := os.Stat(path); os.IsExist(err) {
		if _, err = os.Create(path); err != nil {
			return err
		}
	}
	return nil

}


