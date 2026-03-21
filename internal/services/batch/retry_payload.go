package batch

import (
	"fmt"

	"github.com/dynamder/synthdata/internal/services/llm"
)

type RetryPayloadBuilder struct{}

func NewRetryPayloadBuilder() *RetryPayloadBuilder {
	return &RetryPayloadBuilder{}
}

func (b *RetryPayloadBuilder) Build(failed FailedRecord, originalPrompt string) string {
	switch failed.Error.ErrorCode {
	case llm.InvalidFormat:
		{
			return fmt.Sprintf(
				`The following JSON failed to parse. Please fix any JSON syntax errors and return valid JSON:

Invalid JSON:
%s

Original request:
%s`, failed.OriginalOutput, originalPrompt,
			)
		}
	case llm.NoResponse:
		{
			return originalPrompt
		}
	case llm.GenerationFailed:
		{
			//This is mostly caused by schema validation failed(in the future)
			return failed.OriginalOutput
		}
	default:
		{
			return ""
		}
	}
}

func BuildRetryPayload(failed FailedRecord, promptTemplate string) string {
	return NewRetryPayloadBuilder().Build(failed, promptTemplate)
}
