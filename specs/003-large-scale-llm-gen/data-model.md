# Data Model: Large Scale LLM Generation

## Entities

### GenerationRequest
Represents a request to generate N records with batch processing options.

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| DescriptionFile | *models.DescriptionFile | Yes | Original description file with schema |
| TargetCount | int | Yes | Total number of records to generate |
| BatchSize | int | No | Records per batch (default: 100) |
| Concurrency | int | No | Max concurrent LLM calls (default: 5) |
| MaxRetries | int | No | Max retries per failed record (default: 3) |

### BatchResult
Contains results from a single batch of LLM calls.

| Field | Type | Description |
|-------|------|-------------|
| BatchID | int | Sequential batch number |
| SuccessfulRecords | []map[string]interface{} | Records that parsed successfully |
| FailedRecords | []FailedRecord | Records that failed to parse |
| LLMCallCount | int | Number of LLM calls made |

### FailedRecord
Contains original LLM output that failed parsing and error information.

| Field | Type | Description |
|-------|------|-------------|
| OriginalOutput | string | Raw LLM response that failed to parse |
| Error | string | Parsing error message |
| RetryCount | int | Number of retries attempted (0 = first attempt) |

### GenerationSession
Tracks overall progress of a large-scale generation task.

| Field | Type | Description |
|-------|------|-------------|
| Request | GenerationRequest | Original request parameters |
| TotalBatches | int | Total batches to process |
| CompletedBatches | int | Batches completed |
| TotalRecords | int | Successfully generated records |
| FailedRecords | []FailedRecord | All failed records (exhausted retries) |
| StartTime | time.Time | When generation started |
| EndTime | time.Time | When generation ended |

---

## State Transitions

### GenerationSession States
```
Pending -> Running -> Completed
                  -> PartiallyCompleted (some records failed permanently)
                  -> Failed (all records failed)
```

### FailedRecord States
```
Pending -> Retrying -> Recovered (retry succeeded)
                  -> Exhausted (max retries reached)
```

---

## Validation Rules

1. **TargetCount**: Must be > 0
2. **BatchSize**: Must be > 0, default 100
3. **Concurrency**: Must be > 0, default 5
4. **MaxRetries**: Must be >= 0, default 3
5. **BatchSize <= TargetCount**: If TargetCount < BatchSize, single batch used

---

## Relationships

```
GenerationRequest (1) -> (many) BatchResult
BatchResult (1) -> (many) FailedRecord
GenerationSession (1) -> (1) GenerationRequest
GenerationSession (1) -> (many) BatchResult
```
