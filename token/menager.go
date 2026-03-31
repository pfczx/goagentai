package token

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)


type TokenMenager struct {
	mu sync.Mutex
	path              string
	tokenBalanceLimit int
	tokenBalanceUsed  int
}


func InitTokenMenager(path string, balanceLimit int) (*TokenMenager, error) {
	tm := &TokenMenager{
		path:              path,
		tokenBalanceLimit: balanceLimit,
	}
	err := tm.LoadUsage()
	if err != nil {
		return nil, err
	}
	return tm, nil

}

func (t *TokenMenager) SaveUsage() error {
	data, err := json.Marshal(t.tokenBalanceUsed)
	if err != nil {
		return err
	}
	path := filepath.Join(t.path, "tokenBalanceUsed")
	return os.WriteFile(path, data, 0644)
}

func (t *TokenMenager) LoadUsage() error {
	path := filepath.Join(t.path, "tokenBalanceUsed")
	file, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.tokenBalanceUsed = 0
			return nil
		}
		return err
	}
	var value int
	err = json.Unmarshal(file, &value)
	if err != nil {
		return err
	}
	t.tokenBalanceUsed = value
	return nil
}
