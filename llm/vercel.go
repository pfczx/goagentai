package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Vercel struct {
	ApiKey           string
	InternalProvider string
	Model            string
}

func NewVercel(model string, iternalProvider string) (*Vercel, error) {
	key := os.Getenv("VERCEL")
	provider := &Vercel{
		ApiKey:           key,
		InternalProvider: iternalProvider,
		Model:            model,
	}
	if key == "" {
		fmt.Println("There is no key in .env file (Vercel), add your key to .env file in .config/goagent and restart or switch provider")
		return provider, nil
	}
	return provider, nil
}

func (v *Vercel) Name() string {
	return "Vercel"
}

func (v *Vercel) IternalProviderName() string {
	return v.InternalProvider
}

func (v *Vercel) ModelName() string {
	return v.Model
}
func (v *Vercel) SwitchModel(model string) error {
	v.Model = model
	return nil
}

func (v *Vercel) SwitchIternalProvider(provider string) error {
	v.InternalProvider = provider
	return nil
}

func (v *Vercel) Generate(message ChatMessage) (*ChatResponse, error) {
	endpoint := "https://ai-gateway.vercel.sh/v1/chat/completions"

	//building reques from ChatMessage struct
	var contentParts []map[string]interface{}

	for _, part := range message.Content {
		switch part.Type {
		case "text":
			contentParts = append(contentParts, map[string]interface{}{
				"type": part.Type,
				"text": part.Text,
			})
		case "image_url":
			contentParts = append(contentParts, map[string]interface{}{
				"type": "image_url",
				"image_url": map[string]interface{}{
					"url": part.ImageURL,
				},
			})
		}
	}

	var messages []map[string]interface{}

	if message.SystemPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": message.SystemPrompt,
		})
	}

	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": contentParts,
	})

	// Final payload
	modelString := v.InternalProvider + "/" + v.Model
	payload := map[string]interface{}{
		"model":    modelString,
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

	req.Header.Set("Authorization", "Bearer "+v.ApiKey)
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

	var responseText string
	var usage *Usage

	// Extract assistant text
	if choices, ok := parsed["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := msg["content"].(string); ok {
					responseText = content
				}
			}
		}
	}
	usage = &Usage{}
	// Extract usage
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

func (v *Vercel) ListIternalProviders() ([]string, error) {
	req, err := http.NewRequest("GET", "https://ai-gateway.vercel.sh/v1/models", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+v.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Data []struct {
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	set := make(map[string]struct{})

	for _, m := range raw.Data {
		if m.OwnedBy != "" {
			set[m.OwnedBy] = struct{}{}
		}
	}

	var providers []string
	for p := range set {
		providers = append(providers, p)
	}

	return providers, nil
}

func (v *Vercel) ListProviderModels(provider string, withPhoto bool) ([]string, error) {
	req, err := http.NewRequest("GET", "https://ai-gateway.vercel.sh/v1/models", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+v.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Data []struct {
			ID      string `json:"id"`
			OwnedBy string `json:"owned_by"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var models []string

	for _, m := range raw.Data {
		if m.OwnedBy == provider {
			fullName := strings.Split(m.ID, "/")
			models = append(models, fullName[1])
		}
	}

	return models, nil
}
