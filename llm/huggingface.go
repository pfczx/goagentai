package llm

import (
	"fmt"
	"os"
)

type HuggingFace struct {
	ApiKey string
	Model  string
}

func NewHuggingFace(model string) (*HuggingFace, error) {
	key := os.Getenv("HUGGING_FACE")
	if key == "" {
		return nil, fmt.Errorf("There is no value for %s in .env file", key)
	}

	return &HuggingFace{
		ApiKey: key,
		Model:  model,
	}, nil
}

func (h *HuggingFace) SwitchModel(model string) error {
	h.Model = model
	return nil
}
func (h *HuggingFace) Name() string {
	return "Provider: HuggingFace"
}

func (h *HuggingFace) Generate(prompt string) (string, error) {
	return "", nil
}
