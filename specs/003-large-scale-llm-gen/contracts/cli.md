# CLI Contract: Large Scale Generation

## Command: `generate` (Extended)

Extended flags for large-scale generation:

| Flag | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `--scale`, `-s` | int | No | 10 | Number of records to generate |
| `--batch-size` | int | No | 100 | Records per batch |
| `--concurrency` | int | No | 5 | Max concurrent LLM calls |
| `--max-retries` | int | No | 3 | Max retries per failed record |
| `--description`, `-d` | string | Yes | - | Path to description file |
| `--output`, `-o` | string | Yes | - | Output file path |
| `--format`, `-f` | string | No | json | Output format: json, csv |
| `--config`, `-c` | string | No | - | Config file path |
| `--force` | bool | No | false | Overwrite existing output |

## Output

### Success Output
```
Generating {target} records with {concurrency} concurrent calls...
Batch {i}/{total}: {success_count} success, {fail_count} failed
...
Progress: {total_generated}/{target} records ({percentage}%)
Successfully generated {total} records to {output_file}
```

### Partial Success Output
```
Generating {target} records with {concurrency} concurrent calls...
Progress: {total_generated}/{target} records ({percentage}%)
Generation completed with {failed_count} permanent failures.
Successfully generated {total} records to {output_file}
Run with higher concurrency or more retries to recover failed records.
```

### Error Output
```
Error: {error_message}
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Invalid arguments |
| 2 | Generation failed |
| 3 | Partial success (some records failed permanently) |
