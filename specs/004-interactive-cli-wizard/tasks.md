---

description: "Task list for Interactive CLI Wizard feature implementation"
---

# Tasks: Interactive CLI Argument Wizard

**Input**: Design documents from `/specs/004-interactive-cli-wizard/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are NOT requested in the feature specification - skipping test tasks.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go project**: `cmd/`, `internal/`, `test/` at repository root
- Using existing CLI structure in `internal/cli/`

---

## Phase 1: Setup (Project Initialization)

**Purpose**: No additional setup required - existing project structure (cmd/, internal/cli/) is sufficient for this feature

---

## Phase 2: Foundational (Core Prompt Infrastructure)

**Purpose**: Create the core prompt module that all user stories depend on

- [X] T001 [P] Create internal/cli/prompt/session.go with PromptSession struct and Run method
- [X] T002 [P] Create internal/cli/prompt/metadata.go with ArgumentMetadata struct
- [X] T003 [P] Create internal/cli/prompt/validator.go with ValidationResult struct and validation logic
- [X] T004 [P] Create internal/cli/prompt/prompter.go with buildPrompt function
- [X] T005 Add -i, --interactive flag to GenerateCmd in internal/cli/generate.go
- [X] T006 Wire interactive mode in runGenerate function to check -i flag and enter prompt session

**Checkpoint**: Core prompt infrastructure ready - user story implementation can now begin

---

## Phase 3: User Story 1 - Interactive Argument Input (Priority: P1) 🎯 MVP

**Goal**: Users can run CLI commands with the -i flag to enter arguments through an interactive wizard that prompts for each required argument one at a time.

**Independent Test**: Run `synthdata generate -i` and verify each argument is prompted for individually with clear descriptions

### Implementation for User Story 1

- [X] T007 [P] [US1] Implement PromptSession.Run to iterate through required arguments in internal/cli/prompt/session.go
- [X] T008 [P] [US1] Implement argument metadata extraction from cobra flags in internal/cli/prompt/metadata.go
- [X] T009 [US1] Display clear prompts with argument name, description, and example in internal/cli/prompt/prompter.go
- [X] T010 [US1] Collect user input using bufio.Scanner in internal/cli/prompt/session.go
- [X] T011 [US1] Store collected arguments and pass to command execution in internal/cli/prompt/session.go
- [X] T012 [US1] Handle explicit arguments passed alongside -i flag to skip their prompts

**Checkpoint**: At this point, User Story 1 should be fully functional - users can run `synthdata generate -i` and be prompted for each argument

---

## Phase 4: User Story 2 - Argument Validation and Help (Priority: P2)

**Goal**: Users receive real-time validation feedback and helpful hints while entering arguments interactively.

**Independent Test**: Enter invalid input during interactive mode and verify appropriate error messages appear

### Implementation for User Story 2

- [X] T013 [P] [US2] Implement input validation for string type arguments in internal/cli/prompt/validator.go
- [X] T014 [P] [US2] Implement input validation for int type arguments in internal/cli/prompt/validator.go
- [X] T015 [P] [US2] Implement input validation for bool type arguments in internal/cli/prompt/validator.go
- [X] T016 [US2] Implement enumerated/valid values validation in internal/cli/prompt/validator.go
- [X] T017 [US2] Display validation error messages and re-prompt on invalid input in internal/cli/prompt/session.go
- [X] T018 [US2] Add help keyword (? or help) handling to show argument details in internal/cli/prompt/prompter.go
- [X] T019 [US2] Implement re-prompt loop until valid input or cancellation in internal/cli/prompt/session.go

**Checkpoint**: User Story 2 complete - validation and help features working

---

## Phase 5: User Story 3 - Default Values and Optional Arguments (Priority: P3)

**Goal**: Users can accept suggested defaults for optional arguments or skip them entirely.

**Independent Test**: Run interactive mode with optional arguments that have defaults and verify default behavior

### Implementation for User Story 3

- [X] T020 [P] [US3] Detect optional arguments vs required in internal/cli/prompt/metadata.go
- [X] T021 [P] [US3] Display default value in prompt for optional arguments in internal/cli/prompt/prompter.go
- [X] T022 [US3] Implement empty input handling: use default for optional with default, skip for optional without default in internal/cli/prompt/session.go
- [X] T023 [US3] Implement required argument handling: re-prompt if empty in internal/cli/prompt/session.go

**Checkpoint**: All user stories complete - interactive CLI wizard fully functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T024 [P] Add graceful Ctrl+C handling with clean exit message in internal/cli/prompt/session.go
- [X] T025 [P] Add progress indication showing current argument position in internal/cli/prompt/prompter.go
- [ ] T026 Run quickstart.md validation to ensure documented behavior matches implementation
- [ ] T027 Update README or documentation with -i flag usage if needed

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - existing project structure sufficient
- **Foundational (Phase 2)**: Depends on Setup - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories should proceed in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Builds on US1 infrastructure
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Builds on US1/US2 infrastructure

### Within Each User Story

- Core infrastructure before enhancements
- US1 must be complete before US2 features can be tested fully
- US1 and US2 should be complete before US3

### Parallel Opportunities

- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Within each user story, parallelizable tasks are marked [P]
- Once Foundational phase completes, US1 is MVP - can be tested independently before moving to US2/US3

---

## Parallel Example: Foundational Phase

```bash
# Launch all foundational prompt modules in parallel:
Task: "Create internal/cli/prompt/session.go with PromptSession struct"
Task: "Create internal/cli/prompt/metadata.go with ArgumentMetadata struct"
Task: "Create internal/cli/prompt/validator.go with ValidationResult struct"
Task: "Create internal/cli/prompt/prompter.go with buildPrompt function"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (already done - existing structure)
2. Complete Phase 2: Foundational (core prompt infrastructure)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test interactive mode works with argument prompts
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Foundational → Core prompt infrastructure ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Sequential Strategy

1. Foundational Phase (T001-T006)
2. User Story 1 Phase (T007-T012) → MVP delivered
3. User Story 2 Phase (T013-T019)
4. User Story 3 Phase (T020-T023)
5. Polish Phase (T024-T027)

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- This feature enhances the existing generate command - no breaking changes to existing CLI behavior
