package batch

import (
	"fmt"
	"testing"

	"github.com/dynamder/synthdata/internal/services/llm"
)

func TestRetryPayloadBuilder_Build(t *testing.T) {
	builder := NewRetryPayloadBuilder()

	failed := FailedRecord{
		OriginalOutput: `{"invalid": json}`,
		Error:          llm.NewLLmCallError(llm.InvalidFormat, nil),
		RetryCount:     1,
	}

	prompt := "Generate 10 records for test dataset"
	result := builder.Build(failed, prompt)

	if result == "" {
		t.Error("expected non-empty payload")
	}

	expectedSubstrings := []string{"failed to parse", "Invalid JSON", failed.OriginalOutput, prompt}
	for _, substr := range expectedSubstrings {
		if !contains(result, substr) {
			t.Errorf("expected payload to contain %q", substr)
		}
	}
}

func TestBuildRetryPayload(t *testing.T) {
	failed := FailedRecord{
		OriginalOutput: `{"name": "test"`,
		Error:          llm.NewLLmCallError(llm.InvalidFormat, fmt.Errorf("Missing closing brace")),
		RetryCount:     0,
	}

	promptTemplate := "Generate %d records"
	result := BuildRetryPayload(failed, promptTemplate)

	if result == "" {
		t.Error("expected non-empty payload")
	}

	if !contains(result, failed.OriginalOutput) {
		t.Errorf("expected payload to contain original output %q", failed.OriginalOutput)
	}

	if !contains(result, promptTemplate) {
		t.Errorf("expected payload to contain prompt template %q", promptTemplate)
	}
}

func TestBuildRetryPayload_MultipleRetries(t *testing.T) {
	failed := FailedRecord{
		OriginalOutput: `{"id": 1}`,
		RetryCount:     2,
		Error:          llm.NewLLmCallError(llm.InvalidFormat, nil),
	}

	prompt := "Generate records"
	result := BuildRetryPayload(failed, prompt)

	if result == "" {
		t.Error("expected non-empty payload")
	}

	if !contains(result, "failed to parse") {
		t.Error("expected payload to mention parse failure")
	}
}

func TestRetryPayloadBuilder_EmptyOriginalOutput(t *testing.T) {
	builder := NewRetryPayloadBuilder()

	failed := FailedRecord{
		OriginalOutput: "",
		RetryCount:     0,
		Error:          llm.NewLLmCallError(llm.InvalidFormat, nil),
	}

	result := builder.Build(failed, "test prompt")

	if result == "" {
		t.Error("expected non-empty payload even with empty original output")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
