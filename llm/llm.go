package llm

import (
	"fmt"
)

// input for model include text and photo
type ChatMessage struct {
	SystemPrompt string        `json:"system_prompt"`
	Content      []ContentPart `json:"content"`
}

type ContentPart struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ImageURL *Img   `json:"image_url,omitempty"`
}

type Img struct {
	Url string `json:"url"`
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
	ModelName() string
	IternalProviderName() string
	SwitchModel(model string) error
	SwitchIternalProvider(provider string) error
	ListIternalProviders() ([]string, error)
	ListProviderModels(provider string, withPhoto bool) ([]string, error)
}

func NewProvider(name string, model string, iternalProvider string) (ModelProvider, error) {
	switch name {
	case "HuggingFace":
		return NewHuggingFace(model, iternalProvider)
	case "OpenRouter":
		return NewOpenRouter(model, iternalProvider)
	case "Groq":
		return NewGroq(model, iternalProvider)
	default:
		return nil, fmt.Errorf("unknown provider %s", name)
	}
}

func ListProviders() []string {
	//currently avaible providers
	providers := []string{
		"HuggingFace",
		"OpenRouter",
		"Groq",
	}
	return providers
}
