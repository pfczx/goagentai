package agent

import (
	"encoding/json"
	"os"
)

type Config struct {
	Name                                        string  `json:"name"`
	Path                                        string  `json:"path"`
	Provider                                    string  `json:"provider"`
	IternalProvider                             string  `json:"iternalprovider"`
	Model                                       string  `json:"model"`
	Temperature                                 float64 `json:"temperature"`
	MemoryOn                                    bool    `json:"memory_on"`
	ShortTermMemoryLimit                        int     `json:"short_term_memory_limit"`
	ShortTermMemoryEvaluation                   bool    `json:"short_term_memory_evaluation_on"`
	LongTermMemoryChunksToAdd                   int     `json:"long_term_memory_chunks_to_be_added_as_context"`
	LongTermMemoryBufferSize                    int     `json:"long_term_memory_buffer_size"`
	LongTermMemoryChunkSize                     int     `json:"long_term_memory_chunk_size"`
	LongTermMemoryStorageSize                   int     `json:"long_term_memory_storage_size"`
	LongTermMemorySummarizationProvider         string  `json:"long_term_memory_summarization_provider"`
	LongTermMemorySummarizationInternalProvider string  `json:"long_term_memory_summarization_internal_provider"`
	LongTermMemorySummarizationModel            string  `json:"long_term_memory_summarization_model"`
}

func DefaultConfig(name string, path string) *Config {
	return &Config{
		Name:                                name,
		Path:                                path,
		Provider:                            "HuggingFace",
		IternalProvider:                     "groq",
		Model:                               "openai/gpt-oss-20b",
		Temperature:                         50.0,
		MemoryOn:                            true,
		ShortTermMemoryLimit:                20,
		ShortTermMemoryEvaluation:           true,
		LongTermMemoryChunksToAdd:           5,
		LongTermMemoryBufferSize:            3,
		LongTermMemoryChunkSize:             100,
		LongTermMemoryStorageSize:           1000,
		LongTermMemorySummarizationProvider: "HuggingFace",
		LongTermMemorySummarizationInternalProvider: "groq",
		LongTermMemorySummarizationModel:            "openai/gpt-oss-20b",
	}
}

func SaveConfig(path string, config *Config) error {
	data, err := json.MarshalIndent(config, "", " ")
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
