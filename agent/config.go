package agent

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name        string  `json:"name"`
	Provider    string  `json:"provider"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
}

func DefaultConfig(name string) *Config {
	return &Config{
		Name:        name,
		Provider:    "HuggingFace",
		Model:       "Model",
		Temperature: 50.0,
	}
}

func SaveConfig(path string, config *Config) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
