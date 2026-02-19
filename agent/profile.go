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

func (c *Config) ProfileFromConfig() (*Profile, error) {
	provider, err := llm.NewProvider(c.Name,c.Model,c.IternalProvider)
	if err != nil {
		return nil, err
	}
	return NewProfile(
		c.Name, c.Path, c, provider, c.Temperature,
	), nil
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

/*
func InitDefaultProfile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pathToDefault := filepath.Join(homeDir, ".config", "goagent", "profiles", "default")

	if _, err := os.Stat(pathToDefault); os.IsNotExist(err) {
		if err := os.MkdirAll(pathToDefault, 0755); err != nil {
			return err
		}
	}

	configPath := filepath.Join(pathToDefault, "config.json")

	if _, err = os.Stat(configPath); err != nil {
		return err
	}
	if err = SaveConfig(configPath, DefaultConfig()); err != nil {
		return err
	}

	return nil
}
*/
