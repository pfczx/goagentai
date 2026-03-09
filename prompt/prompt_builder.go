package prompt

import (
	"github.com/pfczx/goagentai/llm"
)

type Context interface{
	AddContext()
}

func BuildAsk(prompt string) (llm.ChatMessage, error) {
	message := llm.ChatMessage{
		SystemPrompt: "You are goagentai cli assistant.\n Answer clearly and concisely. \n Use markdown formatting.",
		Content:      []llm.ContentPart{},
	}
	message.Content = append(message.Content,
		llm.ContentPart{
			Type: "text",
			Text: prompt,
		})

	return message, nil
}
