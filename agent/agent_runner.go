package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
)

func InitAgent(profileName string) (*Agent, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ".config", "goagent", "profiles", profileName, "config.json")
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

func RunAsk(agent *Agent, args ...string) error {
	resp, err := agent.Ask(strings.Join(args," "))
	if err != nil {
		return err
	}
	out, err := glamour.Render(resp.Text, "auto")
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil

}
