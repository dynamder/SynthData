---

description: "Task list for Synthetic Dataset Generation Tool"
---

# Tasks: Synthetic Dataset Generation Tool

**Input**: Design documents from `/specs/002-llm-synthetic-dataset-gen/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: REQUIRED per spec.md (User Scenarios & Testing is mandatory)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Initialize Go module in repository root with go mod init
- [X] T002 [P] Install Cobra CLI: go install github.com/spf13/cobra-cli@latest
- [X] T003 [P] Initialize Cobra project structure using cobra-cli
- [X] T004 Add dependencies: go get github.com/spf13/cobra github.com/spf13/viper github.com/sashabaranov/go-openai github.com/stretchr/testify
- [X] T005 Create project directory structure per plan.md: cmd/, internal/, config/, tests/
- [X] T006 Configure Go formatting (gofmt) and linting

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**CRITICAL**: No user story work can begin until this phase is complete

- [X] T007 Create LLM configuration struct in internal/models/config.go
- [X] T008 [P] Create dataset configuration struct in internal/models/dataset.go
- [X] T009 [P] Implement Viper config loading in internal/config/config.go
- [X] T010 Implement LLM client wrapper in internal/services/llm/client.go
- [X] T010a [P] Define LLM client interface for abstraction in internal/services/llm/interface.go
- [X] T011 Create schema parsing types in internal/models/schema.go
- [X] T012 Setup error types and handling in internal/errors/errors.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Generate Synthetic Dataset (Priority: P1) 🎯 MVP

**Goal**: Generate synthetic dataset from description file with specified format and scale

**Independent Test**: Provide valid description file + output spec → verify generated file matches format and record count

### Tests for User Story 1 (REQUIRED per spec.md)

- [X] T013 [P] [US1] Unit test for JSON formatter in tests/unit/formatter_json_test.go
- [X] T014 [P] [US1] Unit test for CSV formatter in tests/unit/formatter_csv_test.go
- [X] T015 [US1] Integration test: generate JSON dataset in tests/integration/generate_json_test.go
- [X] T016 [US1] Integration test: generate CSV dataset in tests/integration/generate_csv_test.go
- [X] T016a [P] [US1] Unit test for LLM client in tests/unit/llm_client_test.go
- [X] T016b [P] [US1] Unit test for generator service in tests/unit/generator_test.go

### Implementation for User Story 1

- [X] T017 [P] [US1] Implement JSON formatter in internal/formatters/json.go
- [X] T018 [P] [US1] Implement CSV formatter in internal/formatters/csv.go
- [X] T018a [P] [US1] Implement nested JSON structure handling in internal/formatters/json.go
- [X] T019 [US1] Implement dataset generator service in internal/services/generator/generator.go
- [X] T020 [US1] Create generate command in cmd/synthdata/main.go with Cobra
- [X] T021 [US1] Implement file output handling with --force flag logic
- [X] T022 [US1] Add CLI flags: --description, --output, --format, --scale, --config, --force
- [X] T023 [US1] Wire up LLM client to generator, generator to formatter, formatter to CLI

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Validate Description File (Priority: P2)

**Goal**: Validate description file before generation with clear error messages

**Independent Test**: Provide invalid description file → verify clear error messages returned

### Tests for User Story 2 (REQUIRED per spec.md)

- [X] T024 [P] [US2] Unit test: invalid syntax detection in tests/unit/validator_test.go
- [X] T025 [P] [US2] Unit test: missing required fields detection in tests/unit/validator_test.go
- [X] T026 [US2] Integration test: validation error messages in tests/integration/validation_test.go

### Implementation for User Story 2

- [X] T027 [P] [US2] Implement description file parser in internal/services/parser/parser.go
- [X] T028 [P] [US2] Implement schema validator in internal/services/validator/validator.go
- [X] T029 [US2] Add validation step before generation in generate command
- [X] T030 [US2] Implement structured error messages with field/line info
- [X] T031 [US2] Integrate validation into CLI flow
- [X] T031a [US2] Implement schema conformance validator in internal/services/validator/conformance.go
- [X] T031b [US2] Test schema conformance (95% threshold) in tests/unit/conformance_test.go

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Specify Multiple Output Formats (Priority: P3)

**Goal**: Support various output formats (JSON, CSV) and handle unsupported formats gracefully

**Independent Test**: Request same dataset in different formats → verify each is valid and parseable

### Tests for User Story 3 (REQUIRED per spec.md)

- [X] T032 [P] [US3] Test: unsupported format error in tests/unit/formatter_test.go
- [X] T033 [P] [US3] Test: JSON parse validation in tests/unit/formatter_test.go
- [X] T034 [P] [US3] Test: CSV parse validation in tests/unit/formatter_test.go

### Implementation for User Story 3

- [X] T035 [P] [US3] Add format detection and validation in formatters
- [X] T036 [P] [US3] Implement error for unsupported formats
- [X] T037 [US3] Ensure generated JSON is parseable by standard JSON parsers
- [X] T038 [US3] Ensure generated CSV has correct headers and delimiters

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T039 Add TOML config file support with defaults
- [X] T040 [P] Update quickstart.md with accurate usage examples
- [X] T041 Add environment variable support (OPENAI_API_KEY, etc.)
- [X] T042 Run full test suite and fix any failures
- [X] T043 Validate against success criteria: SC-001 through SC-006
- [X] T043a Performance test: generate 1000 records in under 30s in tests/performance/benchmark_test.go
- [X] T043b Scale test: generate 10k records without crash in tests/performance/scale_test.go

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Models before services
- Services before CLI integration
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- All formatters (JSON, CSV) can be implemented in parallel (T017, T018)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Unit test for JSON formatter in tests/unit/formatter_json_test.go"
Task: "Unit test for CSV formatter in tests/unit/formatter_csv_test.go"

# Launch all formatters for User Story 1 together:
Task: "Implement JSON formatter in internal/formatters/json.go"
Task: "Implement CSV formatter in internal/formatters/csv.go"
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
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Summary

| Metric | Value |
|--------|-------|
| Total Tasks | 48 |
| User Story 1 Tasks | 11 |
| User Story 2 Tasks | 8 |
| User Story 3 Tasks | 6 |
| Setup Tasks | 6 |
| Foundational Tasks | 6 |
| Polish Tasks | 5 |

- **MVP Scope**: User Story 1 (Phase 3) after Phase 1 + Phase 2
- **Independent Test Criteria**: Each user story has clear acceptance scenarios in spec.md
- **Parallel Opportunities**: 15 tasks marked [P] can run in parallel

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
