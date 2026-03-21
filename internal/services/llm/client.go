package llm

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dynamder/synthdata/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

// Special Error for LLM, need this for different retry methods.
const (
	ClientNotInitialized = iota
	GenerationFailed
	NoResponse
	InvalidFormat
	Unknown
)

var errorCodeMap = map[int]string{
	ClientNotInitialized: "Client Not Initialized",
	GenerationFailed:     "LLM Generation Failed",
	InvalidFormat:        "Invalid Format for LLM Generation",
	NoResponse:           "No Response from LLM",
	Unknown:              "Unknown LLM error",
}

type LLMCallError struct {
	ErrorCode int
	Detail    error
}

func (error *LLMCallError) Error() string {
	return fmt.Sprintf("LLM Call failed: %s", errorCodeMap[error.ErrorCode])
}
func NewLLmCallError(code int, detail error) *LLMCallError {
	return &LLMCallError{ErrorCode: code, Detail: detail}
}

func NewOpenAIClient() *OpenAIClient {
	cfg := config.GetConfig()
	return NewOpenAIClientWithConfig(cfg.LLM.APIKey, cfg.LLM.BaseURL, cfg.LLM.Model)
}

func NewOpenAIClientWithConfig(apiKey, baseURL, model string) *OpenAIClient {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	cfg.HTTPClient = &http.Client{
		Timeout: 120 * time.Second,
	}
	client := openai.NewClientWithConfig(cfg)
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &OpenAIClient{
		client: client,
		model:  model,
	}
}

func (c *OpenAIClient) Generate(prompt string) (string, error) {
	if c.client == nil {
		return "", NewLLmCallError(ClientNotInitialized, nil)
	}

	req := openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", NewLLmCallError(GenerationFailed, err)
	}

	if len(resp.Choices) == 0 {
		return "", NewLLmCallError(NoResponse, nil)
	}

	return resp.Choices[0].Message.Content, nil
}
