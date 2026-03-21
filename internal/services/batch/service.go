package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	synthdatalog "github.com/dynamder/synthdata/internal"
	"github.com/dynamder/synthdata/internal/formatters"
	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/llm"
	"github.com/dynamder/synthdata/internal/services/parser"
	"github.com/dynamder/synthdata/internal/services/validator"
)

type Formatter interface {
	Format(records []map[string]interface{}) ([]byte, error)
}

type Service struct {
	client   llm.Client
	batcher  *Batcher
	executor *Executor
	verbose  bool
}

func NewService(client llm.Client, concurrency, maxRetries int, verbose bool) *Service {
	return &Service{
		client:   client,
		batcher:  NewBatcher(),
		executor: NewExecutor(client, concurrency),
		verbose:  verbose,
	}
}

func (s *Service) Generate(ctx context.Context, req GenerationRequest) (*GenerationSession, error) {
	logger := synthdatalog.GetLogger()
	if req.TargetCount <= 0 {
		return nil, fmt.Errorf("target count must be greater than 0")
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 10
	}
	if req.Concurrency <= 0 {
		req.Concurrency = 5
	}

	logger.Info("Starting batch generation", map[string]interface{}{
		"target":      req.TargetCount,
		"batch_size":  req.BatchSize,
		"concurrency": req.Concurrency,
	})

	fmt.Printf("Starting batch generation: target=%d, batch_size=%d, concurrency=%d\n",
		req.TargetCount, req.BatchSize, req.Concurrency)

	descFile, err := parser.ParseDescriptionFile(req.DescriptionFile)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if err := validator.ValidateDescription(descFile); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	batches := s.batcher.DivideIntoBatches(req.TargetCount, req.BatchSize)
	session := NewSession(req, len(batches))

	promptTemplate := s.buildDescPrompt(descFile)

	progressChan := make(chan ProgressUpdate, 10)
	go s.printProgress(progressChan, session)

	results := s.executor.ExecuteBatches(ctx, batches, promptTemplate, progressChan)

	session.BatchResults = results
	for _, r := range results {
		session.TotalRecords += len(r.SuccessfulRecords)
		for _, failed := range r.FailedRecords {
			session.FailedRecords += failed.RecordCount
		}
		session.CompletedBatches++
	}

	recovered := 0
	if session.FailedRecords > 0 && req.MaxRetries > 0 {
		recovered = s.processRetries(ctx, results, promptTemplate, session)
		session.RecoveredRecords = recovered
		if recovered > session.FailedRecords {
			recovered = session.FailedRecords
		}
		session.FailedRecords -= recovered
	}

	session.EndTime = time.Now().UnixMilli()

	logger.Info("Batch generation completed", map[string]interface{}{
		"total_records":     session.TotalRecords,
		"failed_records":    session.FailedRecords,
		"recovered_records": session.RecoveredRecords,
		"duration_ms":       session.EndTime - session.StartTime,
	})

	if err := s.writeOutput(session, req.Format, req.Output); err != nil {
		return session, fmt.Errorf("failed to write output: %w", err)
	}

	elapsed := time.Duration(session.EndTime-session.StartTime) * time.Millisecond
	failedBatches := 0
	for _, r := range session.BatchResults {
		if len(r.FailedRecords) > 0 {
			failedBatches++
		}
	}

	fmt.Printf("\n========================================\n")
	fmt.Printf("Generation Complete!\n")
	fmt.Printf("========================================\n")
	fmt.Printf("Total Records:   %d / %d\n", session.TotalRecords, req.TargetCount)
	fmt.Printf("Success Batches: %d / %d\n", len(session.BatchResults)-failedBatches, len(session.BatchResults))
	fmt.Printf("Failed Records: %d\n", session.FailedRecords)
	fmt.Printf("Recovered:      %d\n", session.RecoveredRecords)
	fmt.Printf("Output File:    %s\n", req.Output)
	fmt.Printf("Elapsed Time:   %s\n", elapsed.Round(time.Second))
	fmt.Printf("========================================\n")

	return session, nil
}

func (s *Service) processRetries(ctx context.Context, results []BatchResult, originalPrompt string, session *GenerationSession) int {
	var allFailed []FailedRecord
	for _, r := range results {
		allFailed = append(allFailed, r.FailedRecords...)
	}

	retryQueue := NewRetryQueue(s.client, session.Request.MaxRetries)
	for _, failed := range allFailed {
		retryQueue.Add(failed)
	}

	recovered := 0
	var recoveredRecords []map[string]interface{}
	retryQueue.Process(ctx, originalPrompt, func(record map[string]interface{}) {
		recovered++
		recoveredRecords = append(recoveredRecords, record)
	})

	for i := range session.BatchResults {
		if len(session.BatchResults[i].FailedRecords) > 0 {
			session.BatchResults[i].SuccessfulRecords = append(
				session.BatchResults[i].SuccessfulRecords, recoveredRecords...,
			)
		}
	}

	return recovered
}

func (s *Service) buildDescPrompt(descFile *models.DescriptionFile) string {
	schemaJSON, _ := json.Marshal(descFile.Schema)
	return fmt.Sprintf(`Generate %%d records for a dataset named "%s".
Description: %s
Schema: %s

Return ONLY a valid JSON array of objects with no additional text or explanation.
Each object must conform to the schema provided.`, descFile.Name, descFile.Description, schemaJSON)
}

func (s *Service) printProgress(progressChan chan ProgressUpdate, session *GenerationSession) {
	startTime := time.Now()
	var cumulativeSuccess int
	var cumulativeFailed int

	for update := range progressChan {
		elapsed := time.Since(startTime)
		percent := float64(update.Completed) / float64(update.Total) * 100

		cumulativeSuccess += update.Success
		cumulativeFailed += update.Failed
		totalValidRecords := cumulativeSuccess

		avgTimePerBatch := float64(elapsed.Milliseconds()) / float64(update.Completed)
		etaSeconds := avgTimePerBatch * float64(update.Total-update.Completed) / 1000

		if update.Completed > 1 && (update.Completed%5 == 0 || update.Completed == update.Total) {
			fmt.Printf("\n[Phase Summary] Batch %d/%d (%.1f%%) | Elapsed: %s | Avg: %.1fs/batch | ETA: %.1fs | Records: %d/%d | Success: %d | Failed: %d",
				update.Completed, update.Total, percent,
				elapsed.Round(time.Second),
				avgTimePerBatch/1000,
				etaSeconds,
				totalValidRecords, session.Request.TargetCount,
				cumulativeSuccess, cumulativeFailed)
		} else if update.Completed == 1 {
			fmt.Printf("\n[Phase Summary] Batch %d/%d (%.1f%%) | Elapsed: %s | Records: %d/%d | Success: %d | Failed: %d",
				update.Completed, update.Total, percent,
				elapsed.Round(time.Second),
				totalValidRecords, session.Request.TargetCount,
				cumulativeSuccess, cumulativeFailed)
		}
	}
}

func (s *Service) writeOutput(session *GenerationSession, format, output string) error {
	allRecords := make([]map[string]interface{}, 0)
	for _, r := range session.BatchResults {
		allRecords = append(allRecords, r.SuccessfulRecords...)
	}

	formatter, err := GetFormatter(format)
	if err != nil {
		return fmt.Errorf("formatter error: %w", err)
	}

	outputData, err := formatter.Format(allRecords)
	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	if err := os.WriteFile(output, outputData, 0644); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}

func GetFormatter(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case "json":
		return formatters.NewJSONFormatter(), nil
	case "csv":
		return formatters.NewCSVFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
