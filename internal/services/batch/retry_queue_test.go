package batch

import (
	"testing"
)

func TestRetryQueue_NewRetryQueue(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}

	queue := NewRetryQueue(client, 0)
	if queue.maxRetries != 3 {
		t.Errorf("default maxRetries = %d, want 3", queue.maxRetries)
	}

	queue2 := NewRetryQueue(client, 5)
	if queue2.maxRetries != 5 {
		t.Errorf("maxRetries = %d, want 5", queue2.maxRetries)
	}
}

func TestRetryQueue_Add(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	queue := NewRetryQueue(client, 3)

	queue.Add(FailedRecord{OriginalOutput: `{"id": 1}`, Error: "parse error"})
	queue.Add(FailedRecord{OriginalOutput: `{"id": 2}`, Error: "parse error"})

	if len(queue.queue) != 2 {
		t.Errorf("queue length = %d, want 2", len(queue.queue))
	}
}

func TestRetryQueue_FailedCount(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	queue := NewRetryQueue(client, 3)

	if queue.FailedCount() != 0 {
		t.Errorf("FailedCount() = %d, want 0", queue.FailedCount())
	}

	queue.Add(FailedRecord{OriginalOutput: `{"id": 1}`, Error: "test"})
	if queue.FailedCount() != 1 {
		t.Errorf("FailedCount() = %d, want 1", queue.FailedCount())
	}
}

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		attempt int
		minMs   int64
		maxMs   int64
	}{
		{0, 1000, 2000},
		{1, 2000, 4000},
		{2, 4000, 8000},
		{3, 8000, 16000},
		{10, 30000, 30000}, // capped at 30s
	}

	for _, tt := range tests {
		result := exponentialBackoff(tt.attempt)
		if result.Milliseconds() < tt.minMs || result.Milliseconds() > tt.maxMs {
			t.Errorf("exponentialBackoff(%d) = %v, want between %dms and %dms", tt.attempt, result, tt.minMs, tt.maxMs)
		}
	}
}
