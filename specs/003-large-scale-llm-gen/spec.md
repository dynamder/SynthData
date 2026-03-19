# Feature Specification: Large Scale LLM Data Generation

**Feature Branch**: `003-large-scale-llm-gen`  
**Created**: 2026-03-19  
**Status**: Draft  
**Input**: User description: "now we need to add large scale data generation. If the scale is too large, we should divided into parts. and make parrallel llm calls to generate a small portion with semaphore, then iterate until the required scale is reached.

Since the large scale generation, there might be llm output invalid format. We should keep the record that parse successfully, and let llm retry for the others. for retry, we input the invalid content to llm, also with the description, to let it fix the format or content issue."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate Large Dataset (Priority: P1)

A user needs to generate a large synthetic dataset (e.g., 10,000 records) using LLM calls. The system should automatically divide this into manageable batches and process them in parallel up to a configurable concurrency limit.

**Why this priority**: This is the core value proposition - enabling users to generate datasets at scale without manual intervention or overwhelming the LLM API.

**Independent Test**: Can be fully tested by requesting 1000 records and verifying that exactly 1000 valid records are produced within a reasonable timeframe.

**Acceptance Scenarios**:

1. **Given** a user requests 1000 records, **When** the system processes the request, **Then** exactly 1000 valid records are returned
2. **Given** a user requests 50,000 records, **When** the system processes the request, **Then** the generation is divided into multiple batches and completes successfully with exactly 50,000 valid records

---

### User Story 2 - Handle Partial Failures (Priority: P1)

When LLM returns invalid or unparseable output during large-scale generation, the system should retain successfully parsed records and retry only the failed ones.

**Why this priority**: Prevents losing all progress when some LLM responses are malformed - a common occurrence at scale.

**Independent Test**: Can be tested by injecting malformed responses and verifying valid records are preserved while failures are retried.

**Acceptance Scenarios**:

1. **Given** 10 LLM calls produce 5 valid and 5 invalid responses, **When** the system processes them, **Then** 5 valid records are retained and 5 are queued for retry
2. **Given** retries successfully fix 3 of the 5 invalid responses, **When** retry completes, **Then** total of 8 valid records exist

---

### User Story 3 - Retry Invalid Outputs (Priority: P1)

When LLM produces invalid output, the system should send the invalid content back to the LLM along with the original data description, allowing the LLM to correct the format or content issues.

**Why this priority**: Provides automated recovery from malformed outputs without manual intervention.

**Independent Test**: Can be tested by providing intentionally malformed JSON and verifying the LLM corrects it in the retry.

**Acceptance Scenarios**:

1. **Given** LLM returns JSON with missing required fields, **When** the system retries with the invalid content and description, **Then** the LLM returns properly formatted data
2. **Given** LLM returns non-JSON text, **When** the system retries with the invalid content, **Then** the LLM converts it to valid JSON format

---

### User Story 4 - Control Concurrency (Priority: P2)

Users should be able to configure the maximum number of parallel LLM calls to avoid rate limiting and manage costs.

**Why this priority**: Essential for production use - users need control over API usage and costs.

**Independent Test**: Can be tested by setting concurrency to 2 and verifying no more than 2 LLM calls happen simultaneously.

**Acceptance Scenarios**:

1. **Given** concurrency is set to 3, **When** 10 batches are processed, **Then** no more than 3 LLM calls run concurrently
2. **Given** concurrency is set to 1, **When** batches are processed, **Then** calls happen sequentially

---

### Edge Cases

- What happens when all retries fail for a record?
- What happens when LLM returns completely empty response?
- What happens when scale is 0 or negative?
- What happens when semaphore is set to 0?
- How does the system handle partial network failures during batch processing?
- What happens when the dataset description itself is invalid or insufficient?
- How does the system handle rate limiting from the LLM API (exponential backoff parameters)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept a target record count for dataset generation
- **FR-002**: System MUST divide large generation tasks into smaller batches based on configurable batch size
- **FR-002a**: System MUST process multiple batches in parallel using a semaphore to limit concurrent LLM calls
- **FR-003**: System MUST iterate through batches until the required scale is reached
- **FR-004**: System MUST parse LLM output and identify valid versus invalid records
- **FR-005**: System MUST retain successfully parsed records even when other records in the batch fail
- **FR-006**: System MUST queue failed records for retry with the original invalid content and data description
- **FR-007**: System MUST provide a configurable maximum retry count per record
- **FR-008**: System MUST support a configurable semaphore (max concurrent calls) value
- **FR-009**: System MUST track and report the number of successfully generated records versus failed records
- **FR-010**: System MUST stop and return results when target scale is reached, even if mid-batch
- **FR-011**: System MUST detect rate limit errors and retry with exponential backoff
- **FR-012**: System MUST display progress information in the terminal during generation
- **FR-013**: System MUST output error messages to both command line and log file
- **FR-014**: System MUST mark errors as recovered in both terminal and log file when retries succeed
- **FR-015**: System MUST store logs in JSON Lines format in the logs/ directory

### Key Entities *(include if feature involves data)*

- **GenerationRequest**: Represents a request to generate N records, includes data description, batch size, concurrency limit, and retry settings
- **BatchResult**: Contains the results from a single batch of LLM calls - both successful and failed records
- **FailedRecord**: Contains the original LLM output that failed parsing and the associated error information for retry
- **GenerationSession**: Tracks the overall progress of a large-scale generation task across multiple batches

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can generate datasets up to 100,000 records in a single request
- **SC-002**: System maintains the requested concurrency limit (semaphore) throughout generation
- **SC-003**: At least 95% of LLM calls produce successfully parseable output on first attempt
- **SC-004**: Failed records are successfully recovered through retry mechanism at least 80% of the time
- **SC-005**: Users can generate 1000 records in under 5 minutes with default settings
- **SC-006**: All successfully parsed records are retained regardless of failures in other batches
- **SC-007**: Rate limit errors are automatically handled with exponential backoff without losing any records

## Assumptions

- The LLM API can handle the described retry logic (providing invalid content back with context)
- The data description provided by users is sufficient for the LLM to generate valid output
- Users have appropriate API rate limits for their concurrency settings
- The parsing format expected from LLM is consistent (e.g., JSON array of objects)

## Clarifications

### Session 2026-03-20

- Q: How should rate limit errors be handled? → A: Use exponential backoff for retry on rate limit errors
- Q: What format should the log file use? → A: JSON Lines (one JSON object per line)
- Q: Where should the log file be stored? → A: logs/ directory relative to working directory
