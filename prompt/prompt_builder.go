package prompt

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/kbinani/screenshot"
	"github.com/pfczx/goagentai/llm"
)

func screenshotAndConvert() (string, error) {
	n := screenshot.NumActiveDisplays()

	if n == 0 {
		return "", fmt.Errorf("no displays found")
	}

	var displayIndex int

	if n == 1 {
		displayIndex = 0
	} else {
		fmt.Println("Select display:")

		for i := 0; i < n; i++ {
			b := screenshot.GetDisplayBounds(i)
			fmt.Printf("[%d] %dx%d\n", i, b.Dx(), b.Dy())
		}

		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Enter display number: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			idx, err := strconv.Atoi(input)
			if err != nil || idx < 0 || idx >= n {
				fmt.Println("Invalid selection, try again.")
				continue
			}

			displayIndex = idx
			break
		}
	}

	bounds := screenshot.GetDisplayBounds(displayIndex)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func BuildAsk(prompt string, context []string, triggerScreenshot bool) (llm.ChatMessage, error) {
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

	if triggerScreenshot {
		imgConverted, err := screenshotAndConvert()
		if err != nil {
			return llm.ChatMessage{}, err
		}
		imageString := fmt.Sprintf("data:image/png;base64,%s", imgConverted)

		message.Content = append(message.Content,
			llm.ContentPart{
				Type: "image_url",
				ImageURL: &llm.Img{
					Url: imageString,
				},
			})
	}
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
