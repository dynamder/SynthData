package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var cfg *viper.Viper

type Config struct {
	LLM   LLMConfig   `mapstructure:"llm"`
	Batch BatchConfig `mapstructure:"batch"`
}

type LLMConfig struct {
	APIKey     string `mapstructure:"api_key"`
	BaseURL    string `mapstructure:"base_url"`
	Model      string `mapstructure:"model"`
	MaxRetries int    `mapstructure:"max_retries"`
}

type BatchConfig struct {
	BatchSize   int `mapstructure:"batch_size"`
	Concurrency int `mapstructure:"concurrency"`
	MaxRetries  int `mapstructure:"max_retries"`
}

func (c *Config) GetBatchConfig() BatchConfig {
	return BatchConfig{
		BatchSize:   cfg.GetInt("batch.batch_size"),
		Concurrency: cfg.GetInt("batch.concurrency"),
		MaxRetries:  cfg.GetInt("batch.max_retries"),
	}
}

var configFileUsed string

func LoadConfig() {
	cfg = viper.New()
	cfg.SetConfigName("default")
	cfg.AddConfigPath("./configs")
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
	cfg.SetDefault("batch.batch_size", 10)
	cfg.SetDefault("batch.concurrency", 5)
	cfg.SetDefault("batch.max_retries", 3)

	err := cfg.ReadInConfig()
	if err != nil {
		configFileUsed = ""
		fmt.Fprintf(os.Stderr, "Config file not found: %v\n", err)
	} else {
		configFileUsed = cfg.ConfigFileUsed()
		fmt.Fprintf(os.Stderr, "Loaded config from: %s\n", configFileUsed)
	}
}

func GetConfig() *Config {
	if cfg == nil {
		LoadConfig()
	}

	if configFileUsed == "" {
		fmt.Println("Error: No config file found.")
		fmt.Println("")
		fmt.Println("Please create a config file at ./configs/default.toml with the following content:")
		fmt.Println("")
		fmt.Println("  [llm]")
		fmt.Println("  api_key = \"sk-your-api-key-here\"")
		fmt.Println("  base_url = \"https://api.openai.com/v1\"")
		fmt.Println("  model = \"gpt-4o-mini\"")
		fmt.Println("")
		fmt.Println("  [batch]")
		fmt.Println("  batch_size = 10")
		fmt.Println("  concurrency = 5")
		fmt.Println("  max_retries = 3")
		fmt.Println("")
		fmt.Println("Or set the OPENAI_API_KEY environment variable.")
		os.Exit(1)
	}

	var c Config
	err := cfg.Unmarshal(&c)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	return &c
}
