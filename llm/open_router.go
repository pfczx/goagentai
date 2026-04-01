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

type Openrouter struct {
	ApiKey           string
	InternalProvider string
	Model            string
}

func NewOpenRouter(model string, iternalProvider string) (*Openrouter, error) {
	key := os.Getenv("OPEN_ROUTER")
	provider := &Openrouter{
		ApiKey:           key,
		InternalProvider: iternalProvider,
		Model:            model,
	}
	if key == "" {
		fmt.Println("There is no key in .env file (OpenRouter), add your key to .env file in .config/goagent and restart or switch provider")
		return provider, nil
	}
	return provider, nil
}

func (o *Openrouter) Name() string {
	return "OpenRouter"
}

func (o *Openrouter) IternalProviderName() string {
	return o.InternalProvider
}

func (o *Openrouter) ModelName() string {
	return o.Model
}
func (o *Openrouter) SwitchModel(model string) error {
	o.Model = model
	return nil
}

func (o *Openrouter) SwitchIternalProvider(provider string) error {
	o.InternalProvider = provider
	return nil
}

func (o *Openrouter) Generate(message ChatMessage) (*ChatResponse, error) {
	endpoint := "https://openrouter.ai/api/v1/chat/completions"

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
					"url": part.ImageURL.Url,
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
	modelString := o.InternalProvider + "/" + o.Model
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

	req.Header.Set("Authorization", "Bearer "+o.ApiKey)
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

func (v *Openrouter) ListIternalProviders() ([]string, error) {
	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/providers", nil)
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
			Slug string `json:"slug"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	var providers []string
	for _, p := range raw.Data {
		providers = append(providers, p.Slug)
	}

	return providers, nil
}

func (o *Openrouter) ListProviderModels(provider string, withPhoto bool) ([]string, error) {
	if withPhoto{
		fmt.Println("OpenRouter does not support multimodal filter")
	}
	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/models/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+o.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Data []struct {
			Slug string `json:"canonical_slug"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	modelsMap := make(map[string]struct{})

	for _, m := range raw.Data {
		fullName := strings.Split(m.Slug, "/")
		if fullName[0] == o.InternalProvider {
			modelsMap[fullName[1]] = struct{}{}
		}

	}
	var models []string
	infoString := "Free models cannot be listed via api, visit https://openrouter.ai/models?q=free to check free models then paste model name from quickstart"
	models = append(models, infoString)
	for name := range modelsMap {
		models = append(models, name)
	}

	return models, nil
}
