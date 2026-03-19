package models

type LLMConfig struct {
	APIKey     string `json:"api_key" mapstructure:"api_key"`
	BaseURL    string `json:"base_url" mapstructure:"base_url"`
	Model      string `json:"model" mapstructure:"model"`
	MaxRetries int    `json:"max_retries" mapstructure:"max_retries"`
}
