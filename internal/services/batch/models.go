package batch

type GenerationRequest struct {
	DescriptionFile string
	TargetCount     int
	BatchSize       int
	Concurrency     int
	MaxRetries      int
	Format          string
	Output          string
}

type BatchResult struct {
	BatchID           string
	SuccessfulRecords []map[string]interface{}
	FailedRecords     []FailedRecord
	LLMCallCount      int
	StartTime         int64
	EndTime           int64
}

func (b *BatchResult) RecordCount() int {
	return len(b.SuccessfulRecords) + len(b.FailedRecords)
}

type FailedRecord struct {
	OriginalOutput string
	Error          string
	RetryCount     int
	BatchID        string
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
	return int64(0)
}
