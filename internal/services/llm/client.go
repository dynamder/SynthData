package llm

import (
	"context"
	"fmt"

	"github.com/anomalyco/synthdata/internal/config"
	"github.com/anomalyco/synthdata/internal/errors"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
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
		return "", errors.Wrap("E004", "LLM client not initialized", nil)
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

	resp, err := c.client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", errors.Wrap("E004", fmt.Sprintf("failed to generate: %v", err), err)
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("E004", "no response from LLM")
	}

	return resp.Choices[0].Message.Content, nil
}
