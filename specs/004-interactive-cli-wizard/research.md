# Research: Interactive CLI Argument Wizard

## Clarifications Resolved

### Array/Multi-value Arguments

**Decision**: Interactive mode will support single-value arguments (string, int, bool) initially. Multi-value/array arguments will be prompted with a comma-separated input format.

**Rationale**: 
- Current CLI has no array arguments (only string, int, bool flags)
- Comma-separated is a common, intuitive pattern for CLI arrays
- Allows MVP delivery with clear extension path

**Alternatives considered**:
- Repeat prompt for each value (more interactive but slower)
- Accept JSON input (more flexible but less user-friendly)

## Best Practices for Interactive CLI

### Library Selection

Using Go's standard library with minimal external dependencies:
- `bufio.Scanner` for line input
- `fmt` for formatted prompts
- No additional library needed - cobra already handles flag parsing

### Design Patterns

1. **Prompt Session**: Dedicated struct managing interactive flow
2. **Argument Metadata**: Extract flag metadata (description, default, required) from cobra flags
3. **Validation Layer**: Separate validation logic from prompt logic
4. **Graceful Exit**: Handle Ctrl+C with clean exit

## Implementation Approach

1. Add `-i` flag to root or specific commands
2. Create `internal/cli/prompt/` package with:
   - `session.go`: PromptSession struct and flow control
   - `validator.go`: Input validation logic
   - `prompter.go`: Individual prompt builders
3. Modify command execution to check `-i` flag and enter interactive mode if set
