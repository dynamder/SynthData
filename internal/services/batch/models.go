package batch

import (
	"time"

	"github.com/dynamder/synthdata/internal/services/llm"
)

type GenerationRequest struct {
	DescriptionFile string
	TargetCount     int
	BatchSize       int
	Concurrency     int
	MaxRetries      int
	Format          string
	Output          string
	Verbose         bool
}

type BatchResult struct {
	BatchID           string
	SuccessfulRecords []map[string]interface{}
	FailedRecords     []FailedRecord //TODO: name this field FailureDetail since this slice will always have len 0 or 1.
	LLMCallCount      int
	StartTime         int64
	EndTime           int64
}

func (b *BatchResult) RecordCount() int {
	return len(b.SuccessfulRecords) + len(b.FailedRecords)
}

type FailedRecord struct {
	OriginalOutput string
	Error          *llm.LLMCallError
	RetryCount     int
	BatchID        string
	RecordCount    int //why?
}

type GenerationSession struct {
	Request          GenerationRequest
	TotalBatches     int
	CompletedBatches int
	TotalRecords     int
	FailedRecords    int
	RecoveredRecords int
	StartTime        int64
	EndTime          int64
	BatchResults     []BatchResult
}

func NewSession(req GenerationRequest, totalBatches int) *GenerationSession {
	return &GenerationSession{
		Request:      req,
		TotalBatches: totalBatches,
		StartTime:    nowMillis(),
	}
}

func nowMillis() int64 {
	return time.Now().UnixMilli()
}
