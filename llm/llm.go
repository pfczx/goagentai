package llm

import "fmt"

type ModelProvider interface {
	Generate(prompt string) (string, error)
	Name() string
	SwitchModel(model string) error
}

func NewProvider(name string, model string) (ModelProvider, error) {
	switch name {
	case "HuggingFace":
		return NewHuggingFace(model)
	default:
		return nil, fmt.Errorf("unknown provider %s", name)
	}
}
