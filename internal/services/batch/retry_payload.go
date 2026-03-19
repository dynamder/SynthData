package batch

import "fmt"

type RetryPayloadBuilder struct{}

func NewRetryPayloadBuilder() *RetryPayloadBuilder {
	return &RetryPayloadBuilder{}
}

func (b *RetryPayloadBuilder) Build(failed FailedRecord, originalPrompt string) string {
	return fmt.Sprintf(`The following JSON failed to parse. Please fix any JSON syntax errors and return valid JSON:

Invalid JSON:
%s

Original request:
%s`, failed.OriginalOutput, originalPrompt)
}

func BuildRetryPayload(failed FailedRecord, promptTemplate string) string {
	return fmt.Sprintf(`The following JSON failed to parse. Please fix any JSON syntax errors and return valid JSON:

Invalid JSON:
%s

%s`, failed.OriginalOutput, promptTemplate)
}
