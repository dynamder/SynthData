package batch

import (
	"context"
	"fmt"
	"time"

	synthdatalog "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/services/llm"
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

func (r *RetryQueue) Process(ctx context.Context, originalPrompt string, onRecovered func(map[string]interface{})) (int, error) {
	recovered := 0
	stillFailed := make([]FailedRecord, 0)

	synthdatalog.GetLogger().Info(
		fmt.Sprintf("Retry with max times %d for %d failed batches.", r.maxRetries, r.FailedCount()),
		map[string]interface{}{
			"originalPrompt": originalPrompt,
		},
	)
	fmt.Printf("[Retry] Retry with max times %d for %d failed batches.\n", r.maxRetries, r.FailedCount())

	for _, failed := range r.queue {
		if failed.RetryCount >= r.maxRetries {
			stillFailed = append(stillFailed, failed)
			continue
		}

		failed.RetryCount++
		retryPayload := BuildRetryPayload(failed, originalPrompt)

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

	synthdatalog.GetLogger().Info(
		fmt.Sprintf("Retry completed. Recovered %d records, %d still fail", recovered, r.FailedCount()),
		nil,
	)
	fmt.Printf("[Retry] Retry completed. Recovered %d records, %d still fail\n", recovered, r.FailedCount())

	return recovered, nil
}

func (r *RetryQueue) FailedCount() int {
	return len(r.queue)
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
