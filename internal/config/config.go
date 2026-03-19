package config

import (
	"fmt"
	"os"
	"path/filepath"

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
	cfg.AddConfigPath("$HOME/.synthdata")
	if exePath, err := os.Executable(); err == nil {
		cfg.AddConfigPath(filepath.Dir(exePath))
		cfg.AddConfigPath(filepath.Dir(exePath) + "/config")
	}
	cfg.SetConfigType("toml")

	cfg.SetDefault("llm.api_key", os.Getenv("OPENAI_API_KEY"))
	cfg.SetDefault("llm.base_url", "https://api.openai.com/v1")
	cfg.SetDefault("llm.model", "gpt-4o-mini")
	cfg.SetDefault("llm.max_retries", 3)

	err := cfg.ReadInConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config file not found, using defaults: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Loaded config from: %s\n", cfg.ConfigFileUsed())
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
