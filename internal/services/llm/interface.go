package llm

type Client interface {
	Generate(prompt string) (string, error)
}
