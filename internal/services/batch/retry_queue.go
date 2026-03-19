package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/anomalyco/synthdata/internal/services/llm"
)

type RetryQueue struct {
	client     llm.Client
	maxRetries int
	queue      []FailedRecord
	mu         int
}

func NewRetryQueue(client llm.Client, maxRetries int) *RetryQueue {
	if maxRetries <= 0 {
		maxRetries = 3
	}
	return &RetryQueue{
		client:     client,
		maxRetries: maxRetries,
		queue:      make([]FailedRecord, 0),
	}
}

func (r *RetryQueue) Add(failed FailedRecord) {
	r.queue = append(r.queue, failed)
}

func (r *RetryQueue) Process(ctx context.Context, promptTemplate string, onRecovered func(map[string]interface{})) (int, error) {
	recovered := 0
	stillFailed := make([]FailedRecord, 0)

	for _, failed := range r.queue {
		if failed.RetryCount >= r.maxRetries {
			stillFailed = append(stillFailed, failed)
			continue
		}

		failed.RetryCount++
		retryPayload := buildRetryPayload(failed, promptTemplate)

		var sleepDuration time.Duration
		for attempt := 0; attempt < 3; attempt++ {
			select {
			case <-ctx.Done():
				return recovered, ctx.Err()
			default:
			}

			response, err := r.client.Generate(retryPayload)
			if err != nil {
				sleepDuration = exponentialBackoff(attempt)
				time.Sleep(sleepDuration)
				continue
			}

			records, err := parseRecords(response)
			if err != nil || len(records) == 0 {
				sleepDuration = exponentialBackoff(attempt)
				time.Sleep(sleepDuration)
				continue
			}

			for _, record := range records {
				onRecovered(record)
			}
			recovered++
			break
		}

		if failed.RetryCount < r.maxRetries {
			stillFailed = append(stillFailed, failed)
		}
	}

	r.queue = stillFailed
	return recovered, nil
}

func (r *RetryQueue) FailedCount() int {
	return len(r.queue)
}

func buildRetryPayload(failed FailedRecord, promptTemplate string) string {
	return fmt.Sprintf(`The following JSON failed to parse. Please fix any JSON syntax errors and return valid JSON:

Invalid JSON:
%s

%s`, failed.OriginalOutput, promptTemplate)
}

func exponentialBackoff(attempt int) time.Duration {
	base := time.Second
	maxDelay := 30 * time.Second
	delay := base * (1 << attempt)
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
