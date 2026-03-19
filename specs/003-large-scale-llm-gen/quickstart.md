# Quick Start: Large Scale LLM Generation

## Prerequisites

- Go 1.26+
- OpenAI API key (or compatible API)

## Installation

```bash
go install ./cmd/synthdata
```

## Basic Usage

Generate 1000 records with default settings (batch size: 100, concurrency: 5):

```bash
synthdata generate -d schema.json -o output.json --scale 1000
```

## Advanced Usage

### Control Concurrency

Limit to 3 concurrent LLM calls:

```bash
synthdata generate -d schema.json -o output.json --scale 5000 --concurrency 3
```

### Adjust Batch Size

Process 50 records per batch:

```bash
synthdata generate -d schema.json -o output.json --scale 1000 --batch-size 50
```

### Increase Retry Attempts

Allow up to 5 retries for failed records:

```bash
synthdata generate -d schema.json -o output.json --scale 1000 --max-retries 5
```

### Combine Options

Full configuration:

```bash
synthdata generate \
  -d schema.json \
  -o output.json \
  --scale 10000 \
  --batch-size 200 \
  --concurrency 10 \
  --max-retries 5
```

## Configuration File

Create `config.yaml`:

```yaml
llm:
  api_key: "your-api-key"
  model: "gpt-4o-mini"
  base_url: ""  # For OpenAI-compatible APIs
```

Use custom config:

```bash
synthdata generate -d schema.json -o output.json --scale 1000 -c config.yaml
```

## Description File Format

```json
{
  "name": "users",
  "description": "Synthetic user data for testing",
  "schema": {
    "fields": [
      {"name": "id", "type": "integer", "unique": true},
      {"name": "name", "type": "string"},
      {"name": "email", "type": "string", "unique": true}
    ]
  }
}
```

## Output Formats

- **JSON**: `--format json` (default)
- **CSV**: `--format csv`

## Logs

Logs are stored in `logs/` directory in JSON Lines format:

```bash
tail -f logs/generation-2026-03-20.jsonl
```
