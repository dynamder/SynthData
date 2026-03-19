# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- Cobra (CLI framework) - github.com/spf13/cobra
- go-openai (LLM API client) - github.com/sashabaranov/go-openai
- Viper (config) - github.com/spf13/viper (for TOML config support)
**Storage**: Files (input description files, output datasets)  
**Testing**: Go built-in testing + testify for assertions  
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)  
**Project Type**: CLI tool  
**Performance Goals**: Generate 1000 records in <30 seconds  
**Constraints**: <100MB memory, offline-capable LLM support (Ollama)  
**Scale/Scope**: Single-user CLI, 10k records max per run

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Concise Code | PASS | Single-purpose CLI tool, no unnecessary complexity |
| II. Code Quality | PASS | Standard library + well-maintained dependencies |
| III. Maintainability | PASS | Clear module boundaries: cli/, services/, models/ |
| IV. Readability | PASS | Go's simplicity + cobra's structured commands |
| V. Simplicity First | PASS | MVP-first: core generation first, then formats |
| VI. Modular Architecture | PASS | Separate LLM client, generator, formatter modules |
| VII. MVP-First Development | PASS | Core JSON/CSV generation first, validate description file |

## Project Structure

### Documentation (this feature)

```text
specs/002-llm-synthetic-dataset-gen/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI command schema)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/synthdata/
└── main.go              # Entry point with Cobra root command

internal/
├── cli/
│   └── generate.go      # Generate command implementation
├── config/
│   └── config.go        # Viper config loading
├── services/
│   ├── llm/
│   │   └── client.go    # go-openai wrapper
│   └── generator/
│       └── generator.go # Dataset generation logic
└── models/
    └── dataset.go       # Data structures

config/
└── default.toml         # Default configuration

tests/
├── unit/
│   └── generator_test.go
└── integration/
    └── e2e_test.go
```

**Structure Decision**: Single Go module with clear package separation. CLI (Cobra), LLM client (go-openai), and generator logic are separate packages. Uses standard Go project layout with `cmd/`, `internal/`, and `tests/`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
