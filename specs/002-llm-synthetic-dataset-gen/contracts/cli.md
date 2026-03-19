# CLI Contract

## Command Structure

```
synthdata [OPTIONS] <command> [arguments]
```

## Commands

### generate

Generate a synthetic dataset from a description file.

**Usage:**
```
synthdata generate [OPTIONS]
```

**Options:**
| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `-d, --description` | string | Yes | Path to description Markdown file |
| `-o, --output` | string | Yes | Output file path |
| `-f, --format` | string | No | Output format: json, csv (default: json) |
| `-s, --scale` | int | No | Number of records (default: 100) |
| `-c, --config` | string | No | Config file path |
| `--force` | bool | No | Overwrite existing file |

**Exit Codes:**
- 0: Success
- 1: Validation error
- 2: LLM API error
- 3: File I/O error

### config

Show or edit configuration.

**Usage:**
```
synthdata config [command]
```

**Subcommands:**
- `show`: Display current configuration
- `init`: Create default config file

### version

Show version information.

**Usage:**
```
synthdata version
```

## Config File (TOML)

Location: `~/.synthdata.toml` or path specified by `--config`

```toml
[llm]
api_key = ""
base_url = "https://api.openai.com/v1"
model = "gpt-3.5-turbo"
temperature = 0.7
max_tokens = 2048

[output]
default_format = "json"
default_scale = 100
```

## Environment Variables

- `OPENAI_API_KEY`: LLM API key (overrides config)
- `OPENAI_BASE_URL`: Custom API endpoint
- `SYNTHDATA_CONFIG`: Config file path
