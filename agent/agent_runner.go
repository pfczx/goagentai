package agent

import (
	"os"
	"path/filepath"
)

func InitAgent(profileName string) (*Agent, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ".config", "goagent", profileName, "config.json")
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	profile, err := config.ProfileFromConfig()
	if err != nil {
		return nil, err
	}

	return NewAgent(profile), nil
}
