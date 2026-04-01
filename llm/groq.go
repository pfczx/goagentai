package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

type Groq struct {
	ApiKey           string
	InternalProvider string
	Model            string
}

func NewGroq(model string, iternalProvider string) (*Groq, error) {
	key := os.Getenv("GROQ")
	provider := &Groq{
		ApiKey:           key,
		InternalProvider: iternalProvider,
		Model:            model,
	}
	if key == "" {
		fmt.Println("There is no key in .env file (Groq), add your key to .env file in .config/goagent and restart or switch provider")
		return provider, nil
	}
	return provider, nil
}

func (g *Groq) Name() string {
	return "Groq"
}

func (g *Groq) IternalProviderName() string {
	return g.InternalProvider
}

func (g *Groq) ModelName() string {
	return g.Model
}
func (g *Groq) SwitchModel(model string) error {
	g.Model = model
	return nil
}

func (g *Groq) SwitchIternalProvider(provider string) error {
	g.InternalProvider = provider
	return nil
}

func (g *Groq) Generate(message ChatMessage) (*ChatResponse, error) {
	endpoint := "https://api.groq.com/openai/v1/chat/completions"

	// Spłaszczamy wszystko do jednego stringa
	contentText := ""
	if message.SystemPrompt != "" {
		contentText += message.SystemPrompt + "\n"
	}

	for _, part := range message.Content {
		switch part.Type {
		case "text":
			contentText += part.Text + "\n"
		case "image_url":
			if part.ImageURL != nil {
				contentText += "[image: " + part.ImageURL.Url + "]\n"
			}
		}
	}

	// Tworzymy wiadomości
	messages := []map[string]interface{}{
		{
			"role":    "user",
			"content": contentText, // <-- musi być string
		},
	}

	// Payload
	payload := map[string]interface{}{
		"model":    g.Model,
		"stream":   false,
		"messages": messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}

	// Extract assistant text
	var responseText string
	if choices, ok := parsed["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := msg["content"].(string); ok {
					responseText = content
				}
			}
		}
	}

	// Extract usage
	usage := &Usage{}
	if u, ok := parsed["usage"].(map[string]interface{}); ok {
		if v, ok := u["prompt_tokens"].(float64); ok {
			usage.PromptTokens = int(v)
		}
		if v, ok := u["completion_tokens"].(float64); ok {
			usage.CompletionTokens = int(v)
		}
		if v, ok := u["total_tokens"].(float64); ok {
			usage.TotalTokens = int(v)
		}
	}

	if usage.TotalTokens == 0 {
		fmt.Println("Error occured, raw response: ", string(respBody))
	}

	return &ChatResponse{
		Text:  responseText,
		Usage: usage,
	}, nil
}
func (g *Groq) ListIternalProviders() ([]string, error) {
	//Groq is just groq
	var providers []string
	providers = append(providers, "groq")
	return providers, nil
}

func (g *Groq) ListProviderModels(provider string, withPhoto bool) ([]string, error) {
	if withPhoto {
		fmt.Println("Groq does not support multimodal filter")
	}
	req, err := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Data []struct {
			Id string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var models []string
	for _, m := range raw.Data {
		models = append(models, m.Id)
	}
	sort.Strings(models)
	return models, nil
}
