package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	origEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", origEnv)
	os.Unsetenv("OPENAI_API_KEY")

	LoadConfig()
	defer func() { cfg = nil }()

	c := GetConfig()

	if c.LLM.APIKey != "" {
		t.Errorf("expected empty API key, got %s", c.LLM.APIKey)
	}
	if c.LLM.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("expected default base URL, got %s", c.LLM.BaseURL)
	}
	if c.LLM.Model != "gpt-4o-mini" {
		t.Errorf("expected default model gpt-4o-mini, got %s", c.LLM.Model)
	}
	if c.LLM.MaxRetries != 3 {
		t.Errorf("expected default max retries 3, got %d", c.LLM.MaxRetries)
	}
}

func TestLoadConfig_EnvAPIKey(t *testing.T) {
	os.Setenv("OPENAI_API_KEY", "test-api-key")
	defer os.Unsetenv("OPENAI_API_KEY")

	LoadConfig()
	defer func() { cfg = nil }()

	c := GetConfig()

	if c.LLM.APIKey != "test-api-key" {
		t.Errorf("expected API key from env, got %s", c.LLM.APIKey)
	}
}

func TestBatchConfig_Defaults(t *testing.T) {
	origEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", origEnv)
	os.Unsetenv("OPENAI_API_KEY")

	LoadConfig()
	defer func() { cfg = nil }()

	c := GetConfig()
	batch := c.GetBatchConfig()

	if batch.BatchSize != 10 {
		t.Errorf("expected default batch size 10, got %d", batch.BatchSize)
	}
	if batch.Concurrency != 5 {
		t.Errorf("expected default concurrency 5, got %d", batch.Concurrency)
	}
	if batch.MaxRetries != 3 {
		t.Errorf("expected default max retries 3, got %d", batch.MaxRetries)
	}
}

func TestGetConfig_MultipleCalls(t *testing.T) {
	origEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", origEnv)
	os.Unsetenv("OPENAI_API_KEY")

	c1 := GetConfig()
	c2 := GetConfig()

	if c1 == nil || c2 == nil {
		t.Fatal("expected non-nil config")
	}

	if c1.LLM.Model != c2.LLM.Model {
		t.Errorf("expected same model, got %s and %s", c1.LLM.Model, c2.LLM.Model)
	}
}

func TestConfig_Structure(t *testing.T) {
	origEnv := os.Getenv("OPENAI_API_KEY")
	defer os.Setenv("OPENAI_API_KEY", origEnv)
	os.Setenv("OPENAI_API_KEY", "test-key")

	LoadConfig()
	defer func() { cfg = nil }()

	c := GetConfig()

	if c.LLM.APIKey != "test-key" {
		t.Error("LLM.APIKey should be initialized from env")
	}
	if c.Batch.BatchSize == 0 {
		t.Error("Batch.BatchSize should have default value")
	}
}
