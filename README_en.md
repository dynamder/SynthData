# SynthData

A synthetic dataset generation tool powered by Large Language Models (LLM).

[中文版本](./README_cn.md)

## Overview

SynthData generates synthetic datasets from description files using LLM. It supports various data formats (JSON, CSV), large-scale generation with batching, and offers an interactive CLI wizard for guided data generation.

## Features

- **LLM-Powered Generation**: Generate synthetic data using OpenAI-compatible APIs
- **Multiple Output Formats**: Support for JSON and CSV output
- **Large-Scale Generation**: Batch processing with concurrency control for generating large datasets
- **Interactive Mode**: CLI wizard that guides you through the generation process
- **Schema Validation**: Validate description files before generation

## Installation

```bash
git clone https://github.com/dynamder/synthdata.git
cd synthdata
go build -o synthdata ./cmd/synthdata
```

## Configuration

Create a configuration file (e.g., `configs/default.toml`):

```toml
[llm]
api_key = "your-api-key"
base_url = "https://api.openai.com/v1"
model = "gpt-4o-mini"
max_retries = 3
```

Supported LLM providers: OpenAI, SiliconFlow, Azure OpenAI, and any OpenAI-compatible API.

## Quick Start

```bash
synthdata generate -d description.md -o output.json -s 100
```

## Usage

### Command Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--description` | `-d` | Path to description file | (required) |
| `--output` | `-o` | Output file path | (required) |
| `--format` | `-f` | Output format (json, csv) | json |
| `--scale` | `-s` | Number of records to generate | 10 |
| `--config` | `-c` | Config file path | configs/default.toml |
| `--batch-size` | | Records per batch | 10 |
| `--concurrency` | | Max parallel LLM calls | 5 |
| `--max-retries` | | Max retry attempts | 3 |
| `--force` | | Overwrite existing output | false |
| `--verbose` | `-v` | Enable verbose logging | false |
| `--interactive` | `-i` | Enable interactive wizard | false |

### Interactive Mode

Launch the interactive wizard to guide you through the process:

```bash
synthdata generate --interactive
```

### Large-Scale Generation

For generating large datasets with batch processing:

```bash
synthdata generate -d description.md -o output.json -s 10000 --batch-size 100 --concurrency 10
```

## Description File Format

The description file defines the data structure and generation rules:

```json
{
  "name": "Dataset Name",
  "description_file": "description.md",
  "format": "json",
  "count": 100,
  "schema": {
    "name": "table_name",
    "type": "nested",
    "children": [
      { "name": "id", "type": "integer" },
      { "name": "username", "type": "string" },
      { "name": "email", "type": "string" }
    ]
  }
}
```

For more examples, see `examples/bilibili_chat_description/bilibili_chat.json`.

## Examples

Generate JSON output:
```bash
synthdata generate -d examples/bilibili_chat_description/bilibili_chat.json -o output.json -s 50
```

Generate CSV output:
```bash
synthdata generate -d description.md -o output.csv -f csv -s 100
```

Use custom config:
```bash
synthdata generate -d description.md -o output.json -c my_config.toml
```
