package errors

import (
	"errors"
	"testing"
)

func TestSynthdataError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SynthdataError
		expected string
	}{
		{
			name:     "without cause",
			err:      New("E001", "test error"),
			expected: "E001: test error",
		},
		{
			name:     "with cause",
			err:      Wrap("E002", "wrapped error", errors.New("original")),
			expected: "E002: wrapped error (caused by: original)",
		},
		{
			name:     "empty message",
			err:      New("E003", ""),
			expected: "E003: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestSynthdataError_Unwrap(t *testing.T) {
	original := errors.New("original error")
	err := Wrap("E001", "wrapped", original)

	got := err.Unwrap()
	if got != original {
		t.Errorf("Unwrap() = %v, want %v", got, original)
	}
}

func TestSynthdataError_Unwrap_Nil(t *testing.T) {
	err := New("E001", "no cause")
	got := err.Unwrap()
	if got != nil {
		t.Errorf("Unwrap() = %v, want nil", got)
	}
}

func TestNew(t *testing.T) {
	err := New("E001", "test message")

	if err.Code != "E001" {
		t.Errorf("Code = %s, want E001", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("Message = %s, want test message", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("Cause = %v, want nil", err.Cause)
	}
}

func TestWrap(t *testing.T) {
	original := errors.New("original")
	err := Wrap("E001", "wrapped", original)

	if err.Code != "E001" {
		t.Errorf("Code = %s, want E001", err.Code)
	}
	if err.Message != "wrapped" {
		t.Errorf("Message = %s, want wrapped", err.Message)
	}
	if err.Cause != original {
		t.Errorf("Cause = %v, want %v", err.Cause, original)
	}
}

func TestErrorCodes(t *testing.T) {
	tests := []*SynthdataError{
		ErrInvalidDescription,
		ErrMissingField,
		ErrInvalidFormat,
		ErrLLMConnection,
		ErrGenerationFailed,
		ErrFileWrite,
		ErrValidationFailed,
		ErrInvalidArgs,
		ErrPartialSuccess,
		ErrRetriesExhausted,
	}

	expectedCodes := []string{
		"E001", "E002", "E003", "E004", "E005", "E006", "E007", "E010", "E011", "E012",
	}

	if len(tests) != len(expectedCodes) {
		t.Fatalf("mismatch between test cases and expected codes")
	}

	for i, tt := range tests {
		if tt.Code != expectedCodes[i] {
			t.Errorf("error code %d = %s, want %s", i, tt.Code, expectedCodes[i])
		}
		if tt.Message == "" {
			t.Errorf("error %s has empty message", tt.Code)
		}
	}
}
