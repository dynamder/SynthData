package batch

import (
	"context"
	"fmt"
	"time"

	synthdatalog "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/services/llm"
)

type RetryProgressUpdate struct {
	Completed  int
	Total      int
	Recovered  int
	Failed     int
	DurationMs int64
}

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

func (r *RetryQueue) Process(ctx context.Context, originalPrompt string, onRecovered func(map[string]interface{}), progressChan chan<- RetryProgressUpdate) (int, error) {
	recovered := 0
	stillFailed := make([]FailedRecord, 0)
	totalFailed := r.FailedCount()

	synthdatalog.GetLogger().Info(
		fmt.Sprintf("Retry with max times %d for %d failed batches.", r.maxRetries, totalFailed),
		map[string]interface{}{
			"originalPrompt": originalPrompt,
		},
	)
	fmt.Printf("[Retry] Retry with max times %d for %d failed batches.\n", r.maxRetries, totalFailed)

	startTime := time.Now()

	for i, failed := range r.queue {
		select {
		case <-ctx.Done():
			r.queue = stillFailed
			return recovered, ctx.Err()
		default:
		}

		if failed.RetryCount >= r.maxRetries {
			stillFailed = append(stillFailed, failed)
			continue
		}

		failed.RetryCount++
		retryPayload := BuildRetryPayload(failed, originalPrompt)

		var sleepDuration time.Duration
		retrySuccess := false
		for attempt := 0; attempt < 3; attempt++ {
			select {
			case <-ctx.Done():
				r.queue = stillFailed
				return recovered, ctx.Err()
			default:
			}

			response, err := r.client.GenerateWithBatchSize(retryPayload, failed.RecordCount)
			if err != nil {
				var errorDetail string
				if llmErr, ok := err.(*llm.LLMCallError); ok {
					if llmErr.Detail != nil {
						errorDetail = fmt.Sprintf("%s: %v", llmErr.Message(), llmErr.Detail)
					} else {
						errorDetail = llmErr.Message()
					}
				} else {
					errorDetail = err.Error()
				}
				synthdatalog.GetLogger().Error("Retry attempt failed", map[string]interface{}{
					"batch_id":    failed.BatchID,
					"retry_count": failed.RetryCount,
					"attempt":     attempt + 1,
					"error":       errorDetail,
					"payload":     retryPayload,
				})
				sleepDuration = exponentialBackoff(attempt)
				time.Sleep(sleepDuration)
				continue
			}

			records, err := parseRecords(response)
			if err != nil || len(records) == 0 {
				synthdatalog.GetLogger().Error("Retry parse failed", map[string]interface{}{
					"batch_id":      failed.BatchID,
					"retry_count":   failed.RetryCount,
					"attempt":       attempt + 1,
					"parse_error":   err.Error(),
					"raw_response":  response,
					"cleaned_input": cleanJSONResponse(response),
				})
				sleepDuration = exponentialBackoff(attempt)
				time.Sleep(sleepDuration)
				continue
			}

			for _, record := range records {
				onRecovered(record)
			}
			recovered++
			retrySuccess = true
			break
		}

		if !retrySuccess {
			synthdatalog.GetLogger().Error("Retry exhausted, batch failed permanently", map[string]interface{}{
				"batch_id":     failed.BatchID,
				"retry_count":  failed.RetryCount,
				"max_retries":  r.maxRetries,
				"original_err": failed.Error.Error(),
				"payload":      retryPayload,
			})
		}

		if failed.RetryCount < r.maxRetries {
			stillFailed = append(stillFailed, failed)
		}

		if progressChan != nil {
			elapsed := time.Since(startTime)
			progressChan <- RetryProgressUpdate{
				Completed:  i + 1,
				Total:      totalFailed,
				Recovered:  recovered,
				Failed:     len(stillFailed),
				DurationMs: elapsed.Milliseconds(),
			}
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
