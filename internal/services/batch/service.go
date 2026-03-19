package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anomalyco/synthdata/internal/formatters"
	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/llm"
	"github.com/anomalyco/synthdata/internal/services/parser"
	"github.com/anomalyco/synthdata/internal/services/validator"
)

type Formatter interface {
	Format(records []map[string]interface{}) ([]byte, error)
}

type Service struct {
	client   llm.Client
	batcher  *Batcher
	executor *Executor
}

func NewService(client llm.Client, concurrency, maxRetries int) *Service {
	return &Service{
		client:   client,
		batcher:  NewBatcher(),
		executor: NewExecutor(client, concurrency),
	}
}

func (s *Service) Generate(ctx context.Context, req GenerationRequest) (*GenerationSession, error) {
	if req.TargetCount <= 0 {
		return nil, fmt.Errorf("target count must be greater than 0")
	}
	if req.BatchSize <= 0 {
		req.BatchSize = 100
	}
	if req.Concurrency <= 0 {
		req.Concurrency = 5
	}

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

	promptTemplate := s.buildPromptTemplate(descFile)

	progressChan := make(chan ProgressUpdate, 10)
	go s.printProgress(progressChan, session)

	results := s.executor.ExecuteBatches(ctx, batches, promptTemplate, progressChan)

	session.BatchResults = results
	for _, r := range results {
		session.TotalRecords += len(r.SuccessfulRecords)
		session.FailedRecords += len(r.FailedRecords)
		session.CompletedBatches++
	}

	recovered := 0
	if session.FailedRecords > 0 && req.MaxRetries > 0 {
		recovered = s.processRetries(ctx, results, promptTemplate, session)
		session.RecoveredRecords = recovered
		session.FailedRecords -= recovered
		session.TotalRecords += recovered
	}

	session.EndTime = time.Now().UnixMilli()

	if err := s.writeOutput(session, req.Format, req.Output); err != nil {
		return session, fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Completed: %d records generated, %d failed, %d recovered\n",
		session.TotalRecords, session.FailedRecords, session.RecoveredRecords)

	return session, nil
}

func (s *Service) processRetries(ctx context.Context, results []BatchResult, promptTemplate string, session *GenerationSession) int {
	var allFailed []FailedRecord
	for _, r := range results {
		allFailed = append(allFailed, r.FailedRecords...)
	}

	retryQueue := NewRetryQueue(s.client, session.Request.MaxRetries)
	for _, failed := range allFailed {
		retryQueue.Add(failed)
	}

	recovered := 0
	retryQueue.Process(ctx, promptTemplate, func(record map[string]interface{}) {
		recovered++
		for i := range session.BatchResults {
			session.BatchResults[i].SuccessfulRecords = append(
				session.BatchResults[i].SuccessfulRecords, record,
			)
		}
	})

	return recovered
}

func (s *Service) buildPromptTemplate(descFile *models.DescriptionFile) string {
	schemaJSON, _ := json.Marshal(descFile.Schema)
	return fmt.Sprintf(`Generate %%d records for a dataset named "%s".
Description: %s
Schema: %s

Return ONLY a valid JSON array of objects with no additional text or explanation.
Each object must conform to the schema provided.`, descFile.Name, descFile.Description, schemaJSON)
}

func (s *Service) printProgress(progressChan chan ProgressUpdate, session *GenerationSession) {
	for update := range progressChan {
		percent := float64(update.Completed) / float64(update.Total) * 100
		fmt.Printf("\r[%d/%d] %.1f%% - Success: %d, Failed: %d",
			update.Completed, update.Total, percent, update.Success, update.Failed)
	}
	fmt.Println()
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
