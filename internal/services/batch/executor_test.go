package batch

import (
	"context"
	"testing"
	"time"
)

type mockClient struct {
	generateFunc func(prompt string) (string, error)
}

func (m *mockClient) Generate(prompt string) (string, error) {
	return m.generateFunc(prompt)
}

func (m *mockClient) GenerateWithBatchSize(prompt string, batchSize int) (string, error) {
	return m.generateFunc(prompt)
}

func TestExecutor_NewExecutor(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}

	exec := NewExecutor(client, 0)
	if exec.concurrency != 5 {
		t.Errorf("default concurrency = %d, want 5", exec.concurrency)
	}

	exec2 := NewExecutor(client, 3)
	if exec2.concurrency != 3 {
		t.Errorf("concurrency = %d, want 3", exec2.concurrency)
	}
}

func TestExecutor_ExecuteBatch(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) {
		return `[{"id": 1, "name": "test"}]`, nil
	}}
	exec := NewExecutor(client, 1)

	batch := Batch{ID: 1, Start: 0, End: 10, Size: 10}
	result := exec.ExecuteBatch(context.Background(), batch, "Generate %d records")

	if len(result.SuccessfulRecords) != 1 {
		t.Errorf("successful records = %d, want 1", len(result.SuccessfulRecords))
	}
	if result.LLMCallCount != 1 {
		t.Errorf("LLM call count = %d, want 1", result.LLMCallCount)
	}
}

func TestExecutor_ExecuteBatch_ParseError(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) {
		return "not valid json", nil
	}}
	exec := NewExecutor(client, 1)

	batch := Batch{ID: 1, Start: 0, End: 10, Size: 10}
	result := exec.ExecuteBatch(context.Background(), batch, "Generate %d records")

	if len(result.FailedRecords) == 0 {
		t.Error("expected failed record for parse error")
	}
}

func TestExecutor_ExecuteBatch_LLMError(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) {
		return "", ErrLLMConnection
	}}
	exec := NewExecutor(client, 1)

	batch := Batch{ID: 1, Start: 0, End: 10, Size: 10}
	result := exec.ExecuteBatch(context.Background(), batch, "Generate %d records")

	if len(result.FailedRecords) == 0 {
		t.Error("expected failed record for LLM error")
	}
}

func TestExecutor_Concurrency(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	exec := NewExecutor(client, 3)

	if exec.Concurrency() != 3 {
		t.Errorf("Concurrency() = %d, want 3", exec.Concurrency())
	}
}

func TestCleanJSONResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"[{\"a\":1}]", "[{\"a\":1}]"},
		{"```json\n[{\"a\":1}]\n```", "[{\"a\":1}]"},
		{"text before [{\"a\":1}]", "[{\"a\":1}]"},
		{"[{\"a\":1}]text after", "[{\"a\":1}]"},
		{"  [{\"a\":1}]  ", "[{\"a\":1}]"},
	}

	for _, tt := range tests {
		result := cleanJSONResponse(tt.input)
		if result != tt.expected {
			t.Errorf("cleanJSONResponse(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

var ErrLLMConnection = &testError{message: "connection error"}

type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}

func init() {
	_ = time.Second
}
