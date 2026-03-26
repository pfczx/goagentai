package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pfczx/goagentai/llm"
	"github.com/pfczx/goagentai/memory"
)

type Profile struct {
	Name        string
	Path        string
	Config      *Config
	Provider    llm.ModelProvider
	Temperature float64
}

func NewProfile(name string, path string, config *Config, provider llm.ModelProvider, temperature float64) *Profile {
	return &Profile{
		Name:        name,
		Path:        path,
		Config:      config,
		Provider:    provider,
		Temperature: temperature,
	}
}

func LoadLatestUsedProfileName() (string, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, ".config", "goagent", "latestProfile")
	name, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(name), nil

}

func (p *Profile) SaveLatestUsedProfileName() error {
	path, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path = filepath.Join(path, ".config", "goagent", "latestProfile")
	err = os.WriteFile(path, []byte(p.Name), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) ProfileFromConfig() (*Profile, error) {
	provider, err := llm.NewProvider(c.Provider, c.Model, c.IternalProvider)
	if err != nil {
		return nil, err
	}
	return NewProfile(
		c.Name, c.Path, c, provider, c.Temperature,
	), nil
}

func (p *Profile) UpdateConfigFromProfile(m *memory.MemoryMenager) error {
	config := &Config{
		Path:                                p.Path,
		Name:                                p.Name,
		Provider:                            p.Provider.Name(),
		IternalProvider:                     p.Provider.IternalProviderName(),
		Model:                               p.Provider.ModelName(),
		Temperature:                         p.Temperature,
		MemoryOn:                            m.MemoryOn,
		ShortTermMemoryLimit:                m.ShortTermMemoryLimit,
		ShortTermMemoryEvaluation:           m.ShortTermMemoryEvaluation,
		LongTermMemoryChunksToAdd:           m.LongTermMemoryChunksToAdd,
		LongTermMemoryBufferSize:            m.LongTermMemoryBufferSize,
		LongTermMemoryChunkSize:             m.LongTermMemoryChunkSize,
		LongTermMemoryStorageSize:           m.LongTermMemoryStorageSize,
		LongTermMemorySummarizationProvider: m.LongTermMemorySummarizationProvider.Name(),
		LongTermMemorySummarizationInternalProvider: m.LongTermMemorySummarizationProvider.IternalProviderName(),
		LongTermMemorySummarizationModel:            m.LongTermMemorySummarizationProvider.ModelName(),
	}
	p.Config = config
	err := SaveConfig(filepath.Join(p.Path, "config"), config)
	if err != nil {
		return err
	}
	return nil
}

func InitProfile(args ...string) error {
	if ProfileExists(args[0]) {
		return fmt.Errorf("Profile with this name arleady exists")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(homeDir, ".config", "goagent", "profiles", args[0])
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	configPath := filepath.Join(path, "config")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(configPath, DefaultConfig(args[0], path)); err != nil {
			return err
		}
	}
	err = memory.InitShortMemoryFile(path)
	if err != nil {
		return err
	}
	return nil
}

func ProfileExists(name string) bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	path := filepath.Join(homeDir, ".config", "goagent", "profiles", name)

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func FirstInitialize() error {
	err := InitProfile("default")
	if err != nil {
		return err
	}

	path, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	envPath := filepath.Join(path, ".config", "goagent", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		content := `
HUGGING_FACE=

GROK=

OPENROUTER=
`
		err = os.WriteFile(envPath, []byte(content), 0644)
		if err != nil {
			return err
		}
	}

	latestProfilePath := filepath.Join(path, ".config", "goagent", "latestProfile")
	if _, err := os.Stat(latestProfilePath); os.IsNotExist(err) {
		err = os.WriteFile(latestProfilePath, []byte("default"), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
