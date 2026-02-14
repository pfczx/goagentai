package agent

import (
)

type Config struct {
	Provider string`json:"provider"`
	Model string `json:"model"`
	Temperature float64 `json:"temperature"`	
	
}

func DefaultConfig() *Config{
	return &Config{
		Provider: "HuggingFace",
		Model: "Model",
		Temperature: 50.0,
	}
}


