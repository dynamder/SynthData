# Quickstart: Synthetic Dataset Generation Tool

## Installation

```bash
go install github.com/soul-plan/synthdata@latest
```

## Configuration

Create a `config.toml` file:

```toml
[llm]
api_key = "your-api-key"
base_url = "https://api.openai.com/v1"  # Optional, for custom endpoints
model = "gpt-3.5-turbo"
temperature = 0.7
```

Or set environment variable: `export OPENAI_API_KEY=your-key`

## Usage

### Generate JSON dataset

```bash
synthdata generate --description schema.md --format json --scale 100 --output data.json
```

### Generate CSV dataset

```bash
synthdata generate --description schema.md --format csv --scale 1000 --output data.csv
```

### Use config file

```bash
synthdata generate --config config.toml --description schema.md --output data.json
```

### Overwrite existing file

```bash
synthdata generate --description schema.md --output existing.json --force
```

## Description File Format

You can provide the description in two ways:

### Option 1: Inline JSON description

Create a JSON file with inline description:

```json
{
  "name": "User Profiles",
  "description": "Generate user profiles with the following fields: id (integer, unique), name (string), email (string, valid format), age (integer, 18-100), country (string: US/UK/CA/AU)",
  "format": "json",
  "count": 100,
  "schema": {
    "name": "user",
    "type": "nested",
    "children": [
      {"name": "id", "type": "integer"},
      {"name": "name", "type": "string"},
      {"name": "email", "type": "string"},
      {"name": "age", "type": "integer"},
      {"name": "country", "type": "string"}
    ]
  }
}
```

### Option 2: Reference a Markdown file

Create a JSON file referencing a Markdown description file:

```json
{
  "name": "User Profiles",
  "description_file": "path/to/description.md",
  "format": "json",
  "count": 100,
  "schema": {
    "name": "user",
    "type": "nested",
    "children": [
      {"name": "id", "type": "integer"},
      {"name": "name", "type": "string"},
      {"name": "email", "type": "string"}
    ]
  }
}
```

Then create your Markdown description file (`path/to/description.md`):

```markdown
# Dataset: User Profiles

Generate user profiles with the following fields:
- id: integer, unique identifier
- name: string, full name
- email: string, valid email format
- age: integer, between 18 and 100
- country: string, must be one of US, UK, CA, AU
- created_at: date, ISO 8601 format
- tags: array of strings, 1-5 tags per user
```

**Note**: You must provide either "description" or "description_file" (or both) in the JSON.

## CLI Options

| Flag | Description | Default |
|------|-------------|---------|
| `--description, -d` | Path to description file | Required |
| `--format, -f` | Output format (json, csv) | json |
| `--output, -o` | Output file path | Required |
| `--scale, -s` | Number of records | 100 |
| `--config, -c` | Config file path | ~/.synthdata.toml |
| `--force` | Overwrite existing file | false |
| `--help` | Show help | - |

## Examples

See `examples/` directory for more examples.
