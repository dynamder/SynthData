# Research: Large Scale LLM Data Generation

## Overview
Research findings for implementing large-scale LLM data generation with batch processing, concurrency control, and retry logic.

---

## Decision: Concurrency Control Using Semaphore

**Decision**: Use `golang.org/x/sync/semaphore` or `github.com/sourcegraph/conc` for concurrency control

**Rationale**: 
- The codebase already has `github.com/sourcegraph/conc` as an indirect dependency
- Semaphore pattern is ideal for limiting concurrent LLM API calls
- Provides clean API for acquire/release operations

**Alternatives Considered**:
- Go channels with buffered channel as semaphore - more complex, harder to reason about
- Worker pool pattern - overkill for this use case
- sync.WaitGroup alone - no concurrency limiting

---

## Decision: Retry Logic with Exponential Backoff

**Decision**: Implement retry with exponential backoff and jitter using built-in `time` package

**Rationale**:
- Required by spec for rate limit handling
- Simpler than adding external dependency
- Jitter prevents thundering herd problem

**Alternatives Considered**:
- `github.com/cenkalti/backoff` - external dependency, not needed for simple use case
- Linear backoff - less effective for rate limits

---

## Decision: Batch Processing Architecture

**Decision**: Create new `batch` service that wraps existing generator

**Rationale**:
- Existing `generator.Generate()` works for single batch
- New batch service orchestrates multiple calls with concurrency control
- Separation of concerns - generator handles single prompt, batch handles orchestration

**Alternatives Considered**:
- Modify existing generator - violates single responsibility
- Create entirely new generator - code duplication

---

## Decision: JSON Lines Logging

**Decision**: Store logs in `logs/` directory using JSON Lines format (one JSON object per line)

**Rationale**:
- Required by spec clarification
- Easy to parse and filter
- Append-friendly for long-running processes

**Alternatives Considered**:
- Single JSON file - requires re-reading entire file on append
- Plain text - harder to parse programmatically

---

## Decision: Progress Display

**Decision**: Display progress in terminal using simple print statements with carriage return

**Rationale**:
- CLI tool - terminal output is expected
- Simpler than progress bar libraries
- Works reliably across terminals

**Alternatives Considered**:
- `cheggaa/pb` - external dependency, adds weight
- `tty`-based progress bars - more complex

---

## Summary

All technical decisions are aligned with the Constitution principles:
- **Simplicity First**: Using stdlib where possible, minimal external dependencies
- **Concise Code**: Single-purpose modules for batch processing
- **Maintainability**: Clear separation between generator and batch services
- **MVP-First**: Core batch + semaphore first, retry as enhancement

No unresolved clarifications remain. Proceed to Phase 1: Design.
