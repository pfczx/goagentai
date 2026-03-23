package prompt

import (
	"fmt"

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

func BuildSummarize(preffered_size int, conversations []string) (llm.ChatMessage, error) {
	message := llm.ChatMessage{
		SystemPrompt: fmt.Sprintf(`
You are an AI assistant. Given the following conversation history, summarize it concisely in no more than %d words. Only include necessary information be specyfic if needed and do not add any new information. If there is information about whether reply was useful, mention it. Generate response in simple string. Do not use any special characters and fomratting. 
		`, preffered_size),
		Content: []llm.ContentPart{},
	}
	for _, conversation := range conversations {
		message.Content = append(message.Content,
			llm.ContentPart{
				Type: "text",
				Text: conversation,
			})

	}
	return message, nil
}
