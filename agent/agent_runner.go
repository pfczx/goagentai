package agent

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pfczx/goagentai/memory"
	"github.com/pfczx/goagentai/workspace"
)

func InitMemoryMenagerFromConfig(config *Config) (*memory.MemoryMenager, error) {
	menager, err := memory.InitMenager(
		config.Path,
		config.Name,
		config.MemoryOn,
		config.ShortTermMemoryLimit,
		config.ShortTermMemoryEvaluation,
		config.LongTermMemoryBufferSize,
		config.LongTermMemoryChunkSize,
		config.LongTermMemoryStorageSize,
		config.LongTermMemoryChunksToAdd,
		config.LongTermMemorySummarizationProvider,
		config.LongTermMemorySummarizationInternalProvider,
		config.LongTermMemorySummarizationModel,
	)

	if err != nil {
		return nil, err
	}
	return menager, nil

}

func InitAgent(profileName string) (*Agent, error) {
	path, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ".config", "goagent", "profiles", profileName, "config")
	config, err := LoadConfig(path)
	if err != nil {
		return nil, err
	}
	profile, err := config.ProfileFromConfig()
	if err != nil {
		return nil, err
	}
	Memorymenager, err := InitMemoryMenagerFromConfig(config)
	if err != nil {
		return nil, err
	}
	WorkspaceMenager := workspace.NewWorkspaceMenager()

	return NewAgent(profile, Memorymenager, WorkspaceMenager), nil

}

func RunAsk(agent *Agent, args ...string) error {
	prompt := strings.Join(args, " ")
	var context []string
	//clearing workspace entries set to addOnce
	workspaceString := agent.WorspaceMenager.WorkspaceToString(true)
	if workspaceString != "" {
		context = append(context, workspaceString)
	}

	if agent.MemoryMenager.MemoryOn {
		//last n conversations
		shortTermString := agent.MemoryMenager.ShortTermToString()
		//prompt + short term used for searching relevant long term chunks
		longTermString, err := agent.MemoryMenager.GetLongTermContextString(fmt.Sprintf("%s %s", prompt, shortTermString))
		if err != nil {
			return err
		}

		if shortTermString != "" {
			context = append(context, shortTermString)
		}
		if longTermString != "" {
			context = append(context, longTermString)
		}

	}
	resp, err := agent.Ask(prompt, context)
	if err != nil {
		return err
	}
	out, err := glamour.Render(resp.Text, "auto")
	if err != nil {
		return err
	}
	fmt.Print(out)
	if resp.Usage != nil {
		fmt.Printf("Tokens prompt: %v completion: %v total: %v \n",
			resp.Usage.PromptTokens,
			resp.Usage.CompletionTokens,
			resp.Usage.TotalTokens)
		if agent.MemoryMenager.MemoryOn {
			usefull := false
			if agent.MemoryMenager.ShortTermMemoryEvaluation {
				sc := bufio.NewScanner(os.Stdin)
				fmt.Println("Was it usefull? [y/anything]")

				if sc.Scan() {
					if sc.Text() == "y" {
						usefull = true
					}
				}
			}
			//appending workspace string to prompt for more relevant summarization results
			if workspaceString != "" {
				agent.MemoryMenager.AppendShortTermHistory(prompt+" "+workspaceString, resp.Text, usefull)

			} else {
				agent.MemoryMenager.AppendShortTermHistory(prompt, resp.Text, usefull)

			}
			if err = agent.MemoryMenager.UpdateLongTerm(); err != nil {
				return err
			}
		}

	} else {
		fmt.Println("no token usage error : RunAsk")
	}

	return nil
}

func EditConfig(agent *Agent, args ...string) error {
	configPath := filepath.Join(agent.Profile.Path, "config")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}
	newAgent, err := InitAgent(agent.Profile.Name)
	if err != nil {
		return err
	}
	*agent = *newAgent
	return nil

}
