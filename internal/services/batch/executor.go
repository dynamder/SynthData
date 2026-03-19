package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/anomalyco/synthdata/internal/services/llm"
)

type Executor struct {
	client      llm.Client
	semaphore   chan struct{}
	concurrency int
	wg          sync.WaitGroup
	mu          sync.Mutex
}

func NewExecutor(client llm.Client, concurrency int) *Executor {
	if concurrency <= 0 {
		concurrency = 5
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

	if err != nil {
		result.FailedRecords = append(result.FailedRecords, FailedRecord{
			OriginalOutput: response,
			Error:          err.Error(),
			RetryCount:     0,
			BatchID:        result.BatchID,
		})
		result.EndTime = time.Now().UnixMilli()
		return result
	}

	records, err := parseRecords(response)
	if err != nil {
		result.FailedRecords = append(result.FailedRecords, FailedRecord{
			OriginalOutput: response,
			Error:          fmt.Sprintf("parse error: %v", err),
			RetryCount:     0,
			BatchID:        result.BatchID,
		})
	} else {
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
			e.semaphore <- struct{}{}

			result := e.ExecuteBatch(ctx, b, promptTemplate)

			atomic.AddInt64(&completed, 1)
			if progressChan != nil {
				progressChan <- ProgressUpdate{
					BatchID:   result.BatchID,
					Completed: int(atomic.LoadInt64(&completed)),
					Total:     len(batches),
					Success:   len(result.SuccessfulRecords),
					Failed:    len(result.FailedRecords),
					LLMCalls:  result.LLMCallCount,
					Duration:  result.EndTime - result.StartTime,
				}
			}

			e.mu.Lock()
			results = append(results, *result)
			e.mu.Unlock()

			<-e.semaphore
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
