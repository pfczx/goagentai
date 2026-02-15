package llm

type ModelProvider interface {
	Generate(prompt string) (string, error)
	Name() string
}
