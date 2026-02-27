package agent

import (
	"os"
	"path/filepath"

	"github.com/pfczx/goagentai/llm"
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

func (p *Profile) UpdateConfigFromProfile() error {
	config := &Config{
		Path:            p.Path,
		Name:            p.Name,
		Provider:        p.Provider.Name(),
		IternalProvider: p.Provider.IternalProviderName(),
		Model:           p.Provider.ModelName(),
		Temperature:     p.Temperature,
	}

	err := SaveConfig(p.Path, config)
	if err != nil {
		return err
	}
	return nil
}

func InitProfile(args ...string) error {
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
	configPath := filepath.Join(path, "config.json")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		if err := SaveConfig(configPath, DefaultConfig(args[0], path)); err != nil {
			return err
		}
	}
	return nil
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
