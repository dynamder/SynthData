package batch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	synthdatalog "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/services/llm"
)

type Executor struct {
	client      llm.Client
	semaphore   chan struct{} //TODO: use the official semaphore package for better handling.
	concurrency int
	wg          sync.WaitGroup
	mu          sync.Mutex
}

func NewExecutor(client llm.Client, concurrency int) *Executor {
	if concurrency <= 0 {
		concurrency = 5 //TODO: many place hardcoded the default value, which is tedious.
	}
	return &Executor{
		client:      client,
		semaphore:   make(chan struct{}, concurrency),
		concurrency: concurrency,
	}
}

func (e *Executor) ExecuteBatch(ctx context.Context, batch Batch, promptTemplate string) *BatchResult {
	result := &BatchResult{
		BatchID:   fmt.Sprintf("batch-%d", batch.ID),
		StartTime: time.Now().UnixMilli(),
	}

	e.semaphore <- struct{}{}
	defer func() {
		<-e.semaphore
	}()

	prompt := fmt.Sprintf(promptTemplate, batch.Size)
	response, err := e.client.Generate(prompt)
	result.LLMCallCount = 1

	logger := synthdatalog.GetLogger()

	if err != nil {
		logger.Error("LLM call failed", map[string]interface{}{
			"batch_id":   result.BatchID,
			"batch_size": batch.Size,
			"error":      err.Error(),
			"prompt":     prompt,
			"response":   response,
		})
		if llmErr, ok := errors.AsType[*llm.LLMCallError](err); ok {
			result.FailedRecords = append(result.FailedRecords, FailedRecord{
				OriginalOutput: response,
				Error:          llmErr, //if err return from client, it must be this type
				RetryCount:     0,
				BatchID:        result.BatchID,
				RecordCount:    batch.Size,
			})
		} else {
			logger.Error(
				fmt.Errorf("Unknown LLM generation failure: %w", err).Error(),
				map[string]interface{}{
					"batch_id":   result.BatchID,
					"batch_size": batch.Size,
					"error":      err.Error(),
					"response":   response,
				},
			) //TODO: check if it is right.
			result.FailedRecords = append(result.FailedRecords, FailedRecord{
				OriginalOutput: response,
				Error:          llm.NewLLmCallError(llm.Unknown, err), //if err return from client, it must be this type
				RetryCount:     0,
				BatchID:        result.BatchID,
				RecordCount:    batch.Size,
			})
		}

		result.EndTime = time.Now().UnixMilli()
		return result
	}

	records, err := parseRecords(response)
	if err != nil {
		logger.Error("Failed to parse LLM response", map[string]interface{}{
			"batch_id":   result.BatchID,
			"batch_size": batch.Size,
			"error":      err.Error(),
			"response":   response,
		})
		result.FailedRecords = append(result.FailedRecords, FailedRecord{
			OriginalOutput: response,
			Error:          llm.NewLLmCallError(llm.InvalidFormat, err),
			RetryCount:     0,
			BatchID:        result.BatchID,
			RecordCount:    batch.Size,
		})
	} else {
		logger.Info("Batch completed successfully", map[string]interface{}{
			"batch_id":      result.BatchID,
			"batch_size":    batch.Size,
			"records_count": len(records),
		})
		result.SuccessfulRecords = records
	}

	result.EndTime = time.Now().UnixMilli()
	return result
}

func (e *Executor) ExecuteBatches(ctx context.Context, batches []Batch, promptTemplate string, progressChan chan<- ProgressUpdate) []BatchResult {
	results := make([]BatchResult, 0, len(batches))
	var completed int64

	for _, batch := range batches {
		select {
		case <-ctx.Done():
			return results
		default:
		}

		e.wg.Add(1)
		go func(b Batch) {
			defer e.wg.Done()

			result := e.ExecuteBatch(ctx, b, promptTemplate)

			atomic.AddInt64(&completed, 1)
			if progressChan != nil {
				failedCount := 0
				for _, f := range result.FailedRecords {
					failedCount += f.RecordCount
				}
				progressChan <- ProgressUpdate{
					BatchID:   result.BatchID,
					Completed: int(atomic.LoadInt64(&completed)),
					Total:     len(batches),
					BatchSize: b.Size,
					Success:   len(result.SuccessfulRecords),
					Failed:    failedCount,
					LLMCalls:  result.LLMCallCount,
					Duration:  result.EndTime - result.StartTime,
				}
			}

			e.mu.Lock()
			results = append(results, *result)
			e.mu.Unlock()
		}(batch)
	}

	e.wg.Wait()
	close(progressChan)
	return results
}

func (e *Executor) Concurrency() int {
	return e.concurrency
}

type ProgressUpdate struct {
	BatchID   string
	Completed int
	Total     int
	BatchSize int
	Success   int
	Failed    int
	LLMCalls  int
	Duration  int64
}

func parseRecords(response string) ([]map[string]interface{}, error) {
	cleaned := cleanJSONResponse(response)
	var records []map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &records); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return records, nil
}

func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)
	response = strings.Trim(response, "`")
	if idx := strings.Index(response, "["); idx > 0 {
		response = response[idx:]
	}
	if idx := strings.LastIndex(response, "]"); idx >= 0 && idx < len(response)-1 {
		response = response[:idx+1]
	}
	return response
}
