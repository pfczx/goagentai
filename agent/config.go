package agent

import (
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Config struct {
	Name                                        string  `yaml:"name"`
	Path                                        string  `yaml:"path"`
	Provider                                    string  `yaml:"provider"`
	IternalProvider                             string  `yaml:"iternalprovider"`
	Model                                       string  `yaml:"model"`
	Temperature                                 float64 `yaml:"temperature"`
	MemoryOn                                    bool    `yaml:"memory_on"`
	ShortTermMemoryLimit                        int     `yaml:"short_term_memory_limit"`
	ShortTermMemoryEvaluation                   bool    `yaml:"short_term_memory_evaluation_on"`
	LongTermMemoryChunksToAdd                   int     `yaml:"long_term_memory_chunks_to_be_added_as_context"`
	LongTermMemoryBufferSize                    int     `yaml:"long_term_memory_buffer_size"`
	LongTermMemoryChunkSize                     int     `yaml:"long_term_memory_chunk_size"`
	LongTermMemoryStorageSize                   int     `yaml:"long_term_memory_storage_size"`
	LongTermMemorySummarizationProvider         string  `yaml:"long_term_memory_summarization_provider"`
	LongTermMemorySummarizationInternalProvider string  `yaml:"long_term_memory_summarization_internal_provider"`
	LongTermMemorySummarizationModel            string  `yaml:"long_term_memory_summarization_model"`
	TokenBalanceLimit                           int     `yaml:"token_balance_limit"`
}

func DefaultConfig(name string, path string) *Config {
	return &Config{
		Name:                                name,
		Path:                                path,
		Provider:                            "HuggingFace",
		IternalProvider:                     "nscale",
		Model:                               "Qwen/Qwen3-4B-Instruct-2507",
		Temperature:                         50.0,
		MemoryOn:                            true,
		ShortTermMemoryLimit:                1,
		ShortTermMemoryEvaluation:           true,
		LongTermMemoryChunksToAdd:           5,
		LongTermMemoryBufferSize:            1,
		LongTermMemoryChunkSize:             100,
		LongTermMemoryStorageSize:           1000,
		LongTermMemorySummarizationProvider: "HuggingFace",
		LongTermMemorySummarizationInternalProvider: "nscale",
		LongTermMemorySummarizationModel:            "Qwen/Qwen3-4B-Instruct-2507",
		TokenBalanceLimit:                           1000000,
	}
}

var yamlComments = map[string]string{
	"name:": "  Profile name \n Do not touch \n  Create new profile and delete old one if necessary",

	"path:": " Path to profile directory \n Do not touch",

	"provider:": "  LLM Provider - service that provide models and resources \n Add key in .env file in .config/goagent \n Choose from: \n 1. HuggingFace \n 2. Groq \n OpenRouter",

	"iternalprovider:": "  Some providers like HuggingFace have internal providers \n Check avaible with list internal-providers",

	"model:": "  Model for generating answers \n Check avaible with list models",

	"temperature:": "  Model tendency to output less predictible answers \n Temp: 100 - very high random answers \n Temp: 0 - very predictible answears",

	"memory_on:": "  Memory System switch \n Set `false` when necessary, responses quality will drop but much more less input tokens will be used",

	"short_term_memory_limit:": "  Limit how much previous covnversations will be added to message \n Set high for wide context or low for saving tokens \n Recommended values are between 10-50",

	"short_term_memory_evaluation_on:": "  Evaluation switch \n set `false` if you do not want to rate response \n Recommended to set `true` for improvement of context menaging",

	"long_term_memory_chunks_to_be_added_as_context:": "  Number of long term memory chunks added for context \n Set low value for saving tokens \n Recommended values are between 10 - 30",

	"long_term_memory_buffer_size:": "  Number of old covnversations that will be summarized for one long term chunk \n Lower values will trigger summarization more often \n High values may result in loosing important information \n Recommended values beetwen 3-10",

	"long_term_memory_chunk_size:": "  Preffered number of words for one summarized chunk \n Low values saves output tokens but may lower quality \n Recommended +/- 100 for covnversations \n For more specific tasks (like coding) high values are recommended for memorizing details",

	"long_term_memory_storage_size:": "  Number of long term memory chunks stored in database for one profile \n If the limit is exceeded, the oldest chunks are removed ",

	"long_term_memory_summarization_provider:": "  Provider for summarizing",

	"long_term_memory_summarization_internal_provider:": "  Internal provider for summarizing",

	"long_term_memory_summarization_model:": "  Model for summarizing",
	"token_balance_limit:":                  " Token usage for profile, if reached waring will be printed",
}

func addYAMLComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")

	var out []string

	for _, line := range lines {
		trim := strings.TrimSpace(line)

		for key, comment := range yamlComments {
			if strings.HasPrefix(trim, key) {
				indent := line[:len(line)-len(strings.TrimLeft(line, " "))]
				commentLines := strings.Split(comment, "\n")

				for _, cLine := range commentLines {
					cLine = strings.TrimSpace(cLine)
					if cLine != "" {
						out = append(out, indent+"# "+cLine)
					} else {
						out = append(out, "")
					}
				}

				break
			}
		}

		out = append(out, line)
	}

	return []byte(strings.Join(out, "\n"))
}
func SaveConfig(path string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	data = addYAMLComments(data)

	return os.WriteFile(path, data, 0644)
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
