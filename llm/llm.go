package llm

import (
	"fmt"
)

// input for model include text and photo
type ChatMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ContentPart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
}

// output inlude md sring response, in future forced json with action like docker ps
type ChatResponse struct {
	Text  string
	Usage *Usage
}

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

type ModelProvider interface {
	Generate(message ChatMessage) (*ChatResponse, error)
	Name() string
	SwitchModel(model string) error
	SwitchIternalProvider(provider string) error
}

func NewProvider(name string, model string, iternalProvider string) (ModelProvider, error) {
	switch name {
	case "HuggingFace":
		return NewHuggingFace(model, iternalProvider)
	default:
		return nil, fmt.Errorf("unknown provider %s", name)
	}
}
