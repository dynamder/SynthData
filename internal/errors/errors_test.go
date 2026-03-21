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
			err:      New(1, "test error"),
			expected: "E001: test error",
		},
		{
			name:     "with cause",
			err:      Wrap(2, "wrapped error", errors.New("original")),
			expected: "E002: wrapped error (caused by: original)",
		},
		{
			name:     "empty message",
			err:      New(3, ""),
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
	err := Wrap(1, "wrapped", original)

	got := err.Unwrap()
	if got != original {
		t.Errorf("Unwrap() = %v, want %v", got, original)
	}
}

func TestSynthdataError_Unwrap_Nil(t *testing.T) {
	err := New(1, "no cause")
	got := err.Unwrap()
	if got != nil {
		t.Errorf("Unwrap() = %v, want nil", got)
	}
}

func TestNew(t *testing.T) {
	err := New(1, "test message")

	if err.Code != 1 {
		t.Errorf("Code = %d, want E001", err.Code)
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
	err := Wrap(1, "wrapped", original)

	if err.Code != 1 {
		t.Errorf("Code = %d, want E001", err.Code)
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
		ErrLLMCall,
		ErrGenerationFailed,
		ErrFileWrite,
		ErrValidationFailed,
		ErrInvalidArgs,
		ErrPartialSuccess,
		ErrRetriesExhausted,
	}

	expectedCodes := []int{
		1, 2, 3, 4, 5, 6, 7, 10, 11, 12,
	}

	if len(tests) != len(expectedCodes) {
		t.Fatalf("mismatch between test cases and expected codes")
	}

	for i, tt := range tests {
		if tt.Code != expectedCodes[i] {
			t.Errorf("error code %d = %d, want %d", i, tt.Code, expectedCodes[i])
		}
		if tt.Message == "" {
			t.Errorf("error %d has empty message", tt.Code)
		}
	}
}
