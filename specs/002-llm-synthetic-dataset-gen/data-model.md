# Data Model: Synthetic Dataset Generation Tool

## Entities

### 1. DatasetConfig

Configuration for dataset generation.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| DescriptionFile | string | Path to Markdown description file | Required, must exist |
| OutputFormat | string | Output format (json, csv) | Required, enum: json, csv |
| OutputPath | string | Path for output file | Required |
| Scale | int | Number of records to generate | Required, 0-10000 |
| Force | bool | Overwrite existing file | Default: false |

### 2. LLMConfig

LLM API configuration.

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| APIKey | string | API authentication key | Env: OPENAI_API_KEY |
| BaseURL | string | OpenAI-compatible API endpoint | https://api.openai.com/v1 |
| Model | string | Model name | gpt-3.5-turbo |
| Temperature | float | Sampling temperature | 0.7 |
| MaxTokens | int | Max tokens per request | 2048 |

### 3. GenerationRequest

Request sent to LLM for data generation.

| Field | Type | Description |
|-------|------|-------------|
| Prompt | string | Full prompt with schema description |
| Schema | SchemaDefinition | Data schema from description |
| Count | int | Number of records requested |

### 4. SchemaDefinition

Data schema defined in description file.

| Field | Type | Description |
|-------|------|-------------|
| Fields | []FieldDefinition | List of fields |
| Name | string | Dataset name |

### 5. FieldDefinition

Definition for a single field.

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Field name |
| Type | string | Data type (string, int, float, bool, date, array) |
| Constraints | []Constraint | Field constraints |

### 6. Constraint

Validation constraint for a field.

| Field | Type | Description |
|-------|------|-------------|
| Type | string | Constraint type (min, max, enum, minLength, maxLength) |
| Value | any | Constraint value |

## State Transitions

```
IDLE → VALIDATING → GENERATING → FORMATTING → COMPLETE
                ↓
              ERROR
```

## Validation Rules

1. Description file must exist and be valid Markdown
2. Scale must be 0-10000
3. Output format must be supported (json, csv)
4. Output file must not exist unless --force is set
5. LLM config must have valid API key or error message
