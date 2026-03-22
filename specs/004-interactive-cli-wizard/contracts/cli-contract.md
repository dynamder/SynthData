# CLI Contract: Interactive Mode

## Command Interface

### Flag: `-i`, `--interactive`

**Type**: Boolean flag  
**Purpose**: Enable interactive argument input mode  
**Behavior**: When present, CLI enters wizard mode to collect arguments interactively

### Interaction Flow

```
$ synthdata generate -i

Welcome to Interactive Mode!
Press Ctrl+C at any time to exit.

▶ Description file (-d, --description):
  > Enter path to description file
  Example: schema.json

▶ Output file (-o, --output):
  > Enter output file path
  Example: output.json

... (continues for each required argument)
```

### Prompt Format

Each prompt displays:
1. Argument name and flag(s)
2. Description from help text
3. Example value (if available)
4. Default value (if optional)
5. Input prompt (`> `)

### Input Handling

| Input | Behavior |
|-------|----------|
| Empty + required arg | Show error, re-prompt |
| Empty + optional arg with default | Use default value |
| Empty + optional arg no default | Skip argument |
| `?` or `help` | Show detailed help for argument |
| `Ctrl+C` | Exit cleanly with message |
| Invalid input | Show validation error, re-prompt |

### Explicit Arguments with `-i`

When `-i` is combined with explicit arguments:
- Provided args are skipped in interactive prompts
- User only prompted for missing required args

Example:
```
$ synthdata generate -i -o output.json
# Only prompts for -d (description file)
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Validation error or user cancellation |
| 2 | Unexpected error |
