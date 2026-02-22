package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"io"
)

type HuggingFace struct {
	ApiKey           string
	InternalProvider string
	Model            string
}

func NewHuggingFace(model string, iternalProvider string) (*HuggingFace, error) {
	key := os.Getenv("HUGGING_FACE")
	if key == "" {
		return nil, fmt.Errorf("There is no value for %s in .env file", key)
	}

	return &HuggingFace{
		ApiKey:           key,
		InternalProvider: iternalProvider,
		Model:            model,
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
	endpoint := "https://router.huggingface.co/v1/chat/completions"

	//building reques from ChatMessage struct 
	var contentParts []map[string]interface{}

	for _, part := range message.Content {
		switch part.Type {
		case "text":
			contentParts = append(contentParts, map[string]interface{}{
				"type": "text",
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

	if message.SystenPrompt != "" {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": message.SystenPrompt,
		})
	}

	messages = append(messages, map[string]interface{}{
		"role":    "user",
		"content": contentParts,
	})

	// Final payload
	modelString := h.Model + ":" + h.InternalProvider
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

	req.Header.Set("Authorization", "Bearer "+h.ApiKey)
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

	// Extract usage
	if u, ok := parsed["usage"].(map[string]interface{}); ok {
		usage = &Usage{}

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

	return &ChatResponse{
		Text:  responseText,
		Usage: usage,
	}, nil
}
