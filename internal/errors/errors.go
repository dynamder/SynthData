package errors

import "fmt"

type SynthdataError struct {
	Code    string
	Message string
	Cause   error
}

func (e *SynthdataError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SynthdataError) Unwrap() error {
	return e.Cause
}

func New(code, message string) *SynthdataError {
	return &SynthdataError{Code: code, Message: message}
}

func Wrap(code, message string, cause error) *SynthdataError {
	return &SynthdataError{Code: code, Message: message, Cause: cause}
}

var (
	ErrInvalidDescription = New("E001", "invalid description file")
	ErrMissingField       = New("E002", "missing required field")
	ErrInvalidFormat      = New("E003", "invalid output format")
	ErrLLMConnection      = New("E004", "failed to connect to LLM")
	ErrGenerationFailed   = New("E005", "dataset generation failed")
	ErrFileWrite          = New("E006", "failed to write output file")
	ErrValidationFailed   = New("E007", "validation failed")
	ErrInvalidArgs        = New("E010", "invalid command arguments")
	ErrPartialSuccess     = New("E011", "partial success - some records failed")
	ErrRetriesExhausted   = New("E012", "max retries exhausted for record")
)
