package llm

import (
	"fmt"
	"os"
)

type HuggingFace struct {
	ApiKey   string
	InternalProvider string
	Model    string
}

func NewHuggingFace(model string, iternalProvider string) (*HuggingFace, error) {
	key := os.Getenv("HUGGING_FACE")
	if key == "" {
		return nil, fmt.Errorf("There is no value for %s in .env file", key)
	}

	return &HuggingFace{
		ApiKey:   key,
		InternalProvider: iternalProvider,
		Model:    model,
	}, nil
}

func (h *HuggingFace) SwitchModel(model string) error {
	h.Model = model
	return nil
}

func (h *HuggingFace) SwitchIternalProvider(provider string) error {
	h.InternalProvider = provider
	return nil
}
func (h *HuggingFace) Name() string {
	return "Provider: HuggingFace"
}

func (h *HuggingFace) Generate(message ChatMessage) (*ChatResponse, error) {
	url := "https://router.huggingface.co/v1/chat/completions"
	return "", nil
}
