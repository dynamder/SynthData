# Data Model: Interactive CLI Wizard

## Entities

### ArgumentMetadata

Represents metadata for a CLI argument extracted from cobra flags.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Flag name (e.g., "description", "output") |
| Short | string | Short flag alias (e.g., "d", "o") |
| Description | string | Help text from flag definition |
| DefaultValue | any | Default value if optional |
| IsRequired | bool | Whether flag is required |
| Type | string | Value type: "string", "int", "bool" |
| ValidValues | []string | Allowed values if enumerated (optional) |

### PromptSession

Manages the interactive argument collection flow.

| Field | Type | Description |
|-------|------|-------------|
| Command | *cobra.Command | Reference to the command being prompted |
| ProvidedArgs | map[string]any | Args already provided on command line |
| CollectedArgs | map[string]any | Args collected interactively |
| CurrentIndex | int | Position in argument prompt order |

### ValidationResult

Result of validating user input.

| Field | Type | Description |
|-------|------|-------------|
| IsValid | bool | Whether input passes validation |
| ErrorMessage | string | Human-readable error if invalid |
| Value | any | Parsed/converted value if valid |

## State Transitions

```
PromptSession: New → Running → Completed
                    ↓
                  Cancelled (Ctrl+C)
```

## Validation Rules

- **Required flags**: Must have non-empty input or explicit default
- **Int flags**: Must parse as integer within any defined bounds
- **Bool flags**: Accept "y/yes/true/1" or "n/no/false/0"
- **String flags**: No additional validation unless enumerated
- **Enumerated flags**: Input must be in ValidValues list
