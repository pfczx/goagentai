package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Profile struct {
	Name   string
	Path   string
	Config Config
}

func NewProfile(name string, path string, config Config) *Profile {
	return &Profile{
		Name:   name,
		Path:   path,
		Config: config,
	}
}

func InitDefaultProfile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	pathToDefault := filepath.Join(homeDir, ".config", "goagent", "profiles", "default")
	if _, err := os.Stat(pathToDefault); err != nil {
		return fmt.Errorf("Default directory already exists in " + pathToDefault)
	} else {
		if err = os.MkdirAll(pathToDefault, 0755); err != nil {
			return err
		}
		configPath := filepath.Join(pathToDefault, "config.json")
		defaultConfig := DefaultConfig()
		if _, err = os.Stat(configPath); err != nil {
			return err
		}
		data, err := json.Marshal(defaultConfig)
		if err != nil {
			return err
		}
		if err = os.WriteFile(configPath, []byte(data), 0644); err != nil {
			return err
		}

	}

	return nil
}
