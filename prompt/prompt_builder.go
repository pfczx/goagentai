package prompt

import (
	"github.com/pfczx/goagentai/llm"
)

func BuildAsk(prompt string, context []string) (llm.ChatMessage, error) {
	message := llm.ChatMessage{
		SystemPrompt: "You are goagentai cli assistant.\n Answer clearly and concisely. \n Use markdown formatting.",
		Content:      []llm.ContentPart{},
	}
	if len(context) > 0 {
		for _, contextText := range context {
			message.Content = append(message.Content,
				llm.ContentPart{
					Type: "text",
					Text: contextText,
				})
		}
	}

	message.Content = append(message.Content,
		llm.ContentPart{
			Type: "text",
			Text: prompt,
		})

	return message, nil
}
