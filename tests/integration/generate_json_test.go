package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLLMClient struct {
	response string
	err      error
}

func (m *mockLLMClient) Generate(prompt string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

func TestGenerateJSONDataset(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"id": 1, "name": "John Doe", "email": "john@example.com", "age": 30},
			{"id": 2, "name": "Jane Smith", "email": "jane@example.com", "age": 25},
			{"id": 3, "name": "Bob Wilson", "email": "bob@example.com", "age": 35}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users Dataset",
		Description: "A dataset containing user information",
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

	assert.Equal(t, float64(1), records[0]["id"])
	assert.Equal(t, "John Doe", records[0]["name"])
	assert.Equal(t, "john@example.com", records[0]["email"])
	assert.Equal(t, float64(30), records[0]["age"])

	assert.Equal(t, float64(2), records[1]["id"])
	assert.Equal(t, "Jane Smith", records[1]["name"])

	assert.Equal(t, float64(3), records[2]["id"])
	assert.Equal(t, "Bob Wilson", records[2]["name"])
}

func TestGenerateJSONDataset_WriteToFile(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "Test User", "email": "test@example.com"}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Test Dataset",
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

	formatter, err := generator.GetFormatter("json")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.json")

	err = os.WriteFile(outputPath, data, 0644)
	require.NoError(t, err)

	readData, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	var parsed []map[string]interface{}
	err = json.Unmarshal(readData, &parsed)
	require.NoError(t, err)
	require.Len(t, parsed, 1)
	assert.Equal(t, "Test User", parsed[0]["name"])
}

func TestGenerateJSONDataset_NestedStructures(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{
				"name": "Product A",
				"price": 99.99,
				"metadata": {
					"category": "electronics",
					"inStock": true
				}
			}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Products",
		Description: "Products with metadata",
		Format:      "json",
		Count:       1,
		Schema: models.SchemaField{
			Name: "product",
			Type: models.FieldTypeNested,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 1)

	metadata, ok := records[0]["metadata"].(map[string]interface{})
	require.True(t, ok, "metadata should be a nested object")
	assert.Equal(t, "electronics", metadata["category"])
	assert.Equal(t, true, metadata["inStock"])
}

func TestGenerateJSONDataset_ValidJSONOutput(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "User1", "value": 100},
			{"name": "User2", "value": 200}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Test",
		Description: "Test",
		Format:      "json",
		Count:       2,
		Schema: models.SchemaField{
			Name: "test",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)

	formatter, err := generator.GetFormatter("json")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	var parsed []map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Output should be valid JSON")
	assert.Len(t, parsed, 2)
}
