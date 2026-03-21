package unit

import (
	"testing"

	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	mock := &mockLLMClient{
		response: `[{"name": "Alice", "email": "alice@example.com"}]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users",
		Description: "A list of users",
		Format:      "json",
		Count:       1,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "Alice", records[0]["name"])
	assert.Equal(t, "alice@example.com", records[0]["email"])
}

func TestGenerator_GenerateMultipleRecords(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "Alice", "age": 25},
			{"name": "Bob", "age": 30},
			{"name": "Charlie", "age": 35}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users",
		Description: "A list of users",
		Format:      "json",
		Count:       3,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 3)
	assert.Equal(t, "Alice", records[0]["name"])
	assert.Equal(t, "Bob", records[1]["name"])
	assert.Equal(t, "Charlie", records[2]["name"])
}

func TestGenerator_GenerateError(t *testing.T) {
	mock := &mockLLMClient{
		err: assert.AnError,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users",
		Description: "A list of users",
		Format:      "json",
		Count:       1,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.Error(t, err)
	assert.Nil(t, records)
}

func TestGenerator_GenerateInvalidJSON(t *testing.T) {
	mock := &mockLLMClient{
		response: `not valid json`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users",
		Description: "A list of users",
		Format:      "json",
		Count:       1,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse LLM response")
	assert.Nil(t, records)
}

func TestGenerator_CleanJSONResponse(t *testing.T) {
	mock := &mockLLMClient{
		response: "```json\n[{\"name\": \"Test\"}]\n```",
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Test",
		Description: "Test",
		Format:      "json",
		Count:       1,
		Schema: models.SchemaField{
			Name: "test",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 1)
	assert.Equal(t, "Test", records[0]["name"])
}

func TestGenerator_GetFormatter(t *testing.T) {
	tests := []struct {
		format  string
		wantErr bool
	}{
		{"json", false},
		{"JSON", false},
		{"csv", false},
		{"CSV", false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			formatter, err := generator.GetFormatter(tt.format)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, formatter)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, formatter)

			records := []map[string]interface{}{
				{"name": "Test"},
			}
			data, err := formatter.Format(records)
			require.NoError(t, err)
			assert.NotEmpty(t, data)
		})
	}
}
