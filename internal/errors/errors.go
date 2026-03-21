package errors

import "fmt"

type SynthdataError struct {
	Code    int
	Message string
	Cause   error
}

func (e *SynthdataError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%d: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (e *SynthdataError) Unwrap() error {
	return e.Cause
}

func New(code int, message string) *SynthdataError {
	return &SynthdataError{Code: code, Message: message}
}

func Wrap(code int, message string, cause error) *SynthdataError {
	return &SynthdataError{Code: code, Message: message, Cause: cause}
}

var (
	ErrInvalidDescription = New(1, "invalid description file")
	ErrMissingField       = New(2, "missing required field")
	ErrInvalidFormat      = New(3, "invalid output format")
	ErrLLMCall            = New(4, "failed to connect to LLM")
	ErrGenerationFailed   = New(5, "dataset generation failed")
	ErrFileWrite          = New(6, "failed to write output file")
	ErrValidationFailed   = New(7, "validation failed")
	ErrInvalidArgs        = New(10, "invalid command arguments")
	ErrPartialSuccess     = New(11, "partial success - some records failed")
	ErrRetriesExhausted   = New(12, "max retries exhausted for record")
)

const (
	CodeInvalidDescription = iota
	CodeMissingField
	CodeInvalidFormat
	CodeLLMCall
	CodeGenerationFailed
	CodeFileWrite
	CodeValidationFailed
	CodeInvalidArgs
	CodePartialSuccess
	CodeRetriesExhausted
)
