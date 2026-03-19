package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

type Config struct {
	LLM LLMConfig `mapstructure:"llm"`
}

type LLMConfig struct {
	APIKey     string `mapstructure:"api_key"`
	BaseURL    string `mapstructure:"base_url"`
	Model      string `mapstructure:"model"`
	MaxRetries int    `mapstructure:"max_retries"`
}

func LoadConfig() {
	cfg = viper.New()
	cfg.SetConfigName("default")
	cfg.AddConfigPath("./config")
	cfg.AddConfigPath("$HOME/.config/synthdata")
	cfg.SetConfigType("toml")

	cfg.SetDefault("llm.api_key", os.Getenv("OPENAI_API_KEY"))
	cfg.SetDefault("llm.base_url", "https://api.openai.com/v1")
	cfg.SetDefault("llm.model", "gpt-4o-mini")
	cfg.SetDefault("llm.max_retries", 3)

	err := cfg.ReadInConfig()
	if err != nil {
		_ = fmt.Errorf("config file not found, using defaults: %w", err)
	}
}

func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}
	var c Config
	err := cfg.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	return &c
}
