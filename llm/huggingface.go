package llm

type HuggingFace struct {
	ApiKey string
	Model  string
}

func (h *HuggingFace) Name() string {
	return "Provider: HuggingFace"
}

func (h *HuggingFace) Generate(prompt string) (string, error) {
	return "", nil
}
