---

description: "Task list for Large Scale LLM Data Generation feature"
---

# Tasks: Large Scale LLM Data Generation

**Input**: Design documents from `/specs/003-large-scale-llm-gen/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Not requested in spec.md - no test tasks generated

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `src/`, `tests/` at repository root
- Paths follow plan.md structure: `src/cmd/synthdata/`, `src/internal/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan in src/cmd/synthdata/, src/internal/cli/, src/internal/config/, src/internal/errors/, src/internal/formatters/, src/internal/models/, src/internal/services/generator/, src/internal/services/llm/, src/internal/services/parser/, src/internal/services/validator/, src/internal/services/batch/
- [X] T002 Initialize Go 1.26.1 project with dependencies: go.mod with github.com/sashabaranov/go-openai, github.com/spf13/cobra, github.com/spf13/viper, github.com/stretchr/testify, github.com/sourcegraph/conc
- [X] T003 [P] Configure linting (golangci-lint) and formatting tools (gofmt) in .golangci.yml

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Create config package in src/internal/config/ with Viper for flags: --scale, --batch-size, --concurrency, --max-retries, --description, --output, --format, --config, --force
- [X] T005 [P] Create error types package in src/internal/errors/ with error codes: InvalidArgs, GenerationFailed, PartialSuccess
- [X] T006 Setup logging infrastructure in src/internal/logger.go with JSON Lines output to logs/ directory
- [X] T007 [P] Create JSON/CSV formatters in src/internal/formatters/json.go and src/internal/formatters/csv.go
- [X] T008 Verify existing LLM client in src/internal/services/llm/ is accessible for batch service

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

### Phase 2b: Testing Infrastructure (Required by Constitution)

**Purpose**: Unit test coverage for core logic (mandatory per constitution)

- [X] T043 [P] Create unit tests for config package in src/internal/config/
- [X] T044 [P] Create unit tests for error types in src/internal/errors/
- [X] T045 [P] Create unit tests for formatters (JSON/CSV) in src/internal/formatters/
- [X] T046 [P] Create unit tests for logger in src/internal/logger_test.go

---

## Phase 3: User Story 1 - Generate Large Dataset (Priority: P1) 🎯 MVP

**Goal**: Enable users to generate large datasets (e.g., 10,000 records) with automatic batch division and parallel processing using semaphore

**Independent Test**: Request 1000 records and verify exactly 1000 valid records are produced within reasonable timeframe

### Implementation for User Story 1

- [X] T009 [P] [US1] Create GenerationRequest model in src/internal/models/generation_request.go with fields: DescriptionFile, TargetCount, BatchSize, Concurrency, MaxRetries
- [X] T010 [P] [US1] Create BatchResult model in src/internal/models/batch_result.go with fields: BatchID, SuccessfulRecords, FailedRecords, LLMCallCount
- [X] T011 [P] [US1] Create FailedRecord model in src/internal/models/failed_record.go with fields: OriginalOutput, Error, RetryCount
- [X] T012 [P] [US1] Create GenerationSession model in src/internal/models/session.go with fields: Request, TotalBatches, CompletedBatches, TotalRecords, FailedRecords, StartTime, EndTime
- [X] T013 [US1] Implement BatchService core logic in src/internal/services/batch/service.go with semaphore-controlled concurrency (depends on T009, T010, T011, T012)
- [X] T014 [US1] Implement batch division logic in src/internal/services/batch/batcher.go to split TargetCount into batches based on BatchSize
- [X] T015 [US1] Implement parallel LLM call execution with semaphore in src/internal/services/batch/executor.go
- [X] T016 [US1] Add progress tracking and terminal output in BatchService
- [X] T017 [US1] Integrate batch service with CLI generate command in src/internal/cli/generate.go
- [X] T047 [US1] Unit tests for BatchService in src/internal/services/batch/service_test.go
- [X] T048 [US1] Unit tests for Batcher in src/internal/services/batch/batcher_test.go
- [X] T049 [US1] Unit tests for Executor in src/internal/services/batch/executor_test.go
- [X] T054 [US1] Add CLI error output to stderr in src/internal/cli/generate.go (FR-013)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Handle Partial Failures (Priority: P1)

**Goal**: Retain successfully parsed records when other records fail, queue failed records for retry

**Independent Test**: Inject malformed responses, verify valid records preserved while failures queued for retry

### Implementation for User Story 2

- [X] T018 [P] [US2] Implement output parser in src/internal/services/parser/output_parser.go to identify valid vs invalid records
- [X] T019 [US2] Implement partial result handling in BatchService to retain SuccessfulRecords even when FailedRecords exist
- [X] T020 [US2] Implement retry queue in src/internal/services/batch/retry_queue.go to manage failed records for retry
- [X] T021 [US2] Add retry counting and exhaustion logic in FailedRecord handling
- [X] T022 [US2] Implement exit code 3 for partial success scenario in CLI
- [X] T034 Handle empty LLM response edge case
- [X] T035 Handle all retries exhausted scenario (track permanently failed records)
- [X] T050 [US2] Unit tests for OutputParser in src/internal/services/parser/output_parser_test.go
- [X] T051 [US2] Unit tests for RetryQueue in src/internal/services/batch/retry_queue_test.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Retry Invalid Outputs (Priority: P1)

**Goal**: Send invalid content back to LLM along with original description for correction

**Independent Test**: Provide intentionally malformed JSON, verify LLM corrects it in retry

### Implementation for User Story 3

- [X] T023 [P] [US3] Implement retry payload builder in src/internal/services/batch/retry_payload.go to include original invalid content and description
- [X] T052 [US3] Unit tests for RetryPayload in src/internal/services/batch/retry_payload_test.go
- [X] T024 [US3] Implement retry execution logic in RetryQueue to call LLM with retry payloads
- [X] T025 [US3] Implement recovered record handling in BatchService to mark successfully retried records
- [X] T026 [US3] Add terminal and log output for recovered errors (mark as recovered)
- [X] T027 [US3] Implement exponential backoff for rate limit errors in retry logic
- [X] T055 [US3] Add terminal output formatting for recovered errors in src/internal/cli/ (FR-014)

**Checkpoint**: All P1 user stories should now be independently functional

---

## Phase 6: User Story 4 - Control Concurrency (Priority: P2)

**Goal**: Allow users to configure maximum parallel LLM calls to manage API usage and costs

**Independent Test**: Set concurrency to 2, verify no more than 2 LLM calls happen simultaneously

### Implementation for User Story 4

- [X] T028 [US4] Add --concurrency flag validation (must be > 0) in config
- [X] T029 [US4] Implement dynamic semaphore configuration in BatchService based on --concurrency flag
- [X] T030 [US4] Add concurrency limit enforcement verification in executor
- [X] T031 [US4] Document concurrency behavior in CLI help output
- [X] T053 [US4] Unit tests for concurrency enforcement in executor_test.go

**Checkpoint**: At this point, User Stories 1-4 should all work independently

---

## Phase 7: Edge Cases & Polish (Cross-Cutting)

**Purpose**: Handle edge cases and improvements affecting multiple user stories

- [X] T032 [P] Handle zero/negative scale input in validation
- [X] T033 [P] Handle semaphore set to 0 with sensible default
- [X] T036 Validate data description file existence and format
- [X] T037 Add --force flag behavior to overwrite existing output
- [ ] T038 Run quickstart.md validation scenarios

---

## Phase 8: Final Polish

**Purpose**: Final improvements and cleanup

- [ ] T039 Review and finalize all error messages for clarity
- [ ] T040 Verify JSON Lines logging format in logs/ directory
- [ ] T041 Performance testing: verify 1000 records generation within 5 minutes
- [ ] T042 Code cleanup and refactoring across batch service

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-6)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2)
- **Polish (Phase 7-8)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Depends on US1 models (GenerationRequest, BatchResult, FailedRecord)
- **User Story 3 (P1)**: Can start after Foundational (Phase 2) - Depends on US2 retry queue logic
- **User Story 4 (P2)**: Can start after Foundational (Phase 2) - Can integrate in parallel with US1-3

### Within Each User Story

- Models before services
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- Models within US1 marked [P] can run in parallel
- T032 and T033 (edge cases) can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all models for User Story 1 together:
Task: "Create GenerationRequest model in src/internal/models/generation_request.go"
Task: "Create BatchResult model in src/internal/models/batch_result.go"
Task: "Create FailedRecord model in src/internal/models/failed_record.go"
Task: "Create GenerationSession model in src/internal/models/session.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add User Story 4 → Test independently → Deploy/Demo
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (core batch processing)
   - Developer B: User Story 2 (partial failures handling)
   - Developer C: User Story 3 + 4 (retry logic + concurrency)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
