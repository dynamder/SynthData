package llm

type Client interface {
	Generate(prompt string) (string, error)
	GenerateWithBatchSize(prompt string, batchSize int) (string, error)
}
