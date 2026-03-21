package integration

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getHeaderIndex(headers []string) map[string]int {
	index := make(map[string]int)
	for i, h := range headers {
		index[h] = i
	}
	return index
}

func TestGenerateCSVDataset(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "John Doe", "email": "john@example.com", "age": 30},
			{"name": "Jane Smith", "email": "jane@example.com", "age": 25},
			{"name": "Bob Wilson", "email": "bob@example.com", "age": 35}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users Dataset",
		Description: "A dataset containing user information",
		Format:      "csv",
		Count:       3,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 3)

	assert.Equal(t, "John Doe", records[0]["name"])
	assert.Equal(t, "john@example.com", records[0]["email"])
	assert.Equal(t, float64(30), records[0]["age"])

	assert.Equal(t, "Jane Smith", records[1]["name"])
	assert.Equal(t, "Bob Wilson", records[2]["name"])
}

func TestGenerateCSVDataset_WriteToFile(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "Test User", "email": "test@example.com", "age": 30}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Test Dataset",
		Description: "Test",
		Format:      "csv",
		Count:       1,
		Schema: models.SchemaField{
			Name: "test",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)

	formatter, err := generator.GetFormatter("csv")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.csv")

	err = os.WriteFile(outputPath, data, 0644)
	require.NoError(t, err)

	readData, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(readData)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err)
	require.Len(t, csvRecords, 2)
	headers := csvRecords[0]
	row := csvRecords[1]
	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[h] = i
	}
	assert.Equal(t, "Test User", row[headerIndex["name"]])
	assert.Equal(t, "test@example.com", row[headerIndex["email"]])
	assert.Equal(t, "30", row[headerIndex["age"]])
}

func TestGenerateCSVDataset_ProperHeaders(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"name": "Alice", "email": "alice@example.com"},
			{"name": "Bob", "email": "bob@example.com"}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Users",
		Description: "Test",
		Format:      "csv",
		Count:       2,
		Schema: models.SchemaField{
			Name: "user",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)

	formatter, err := generator.GetFormatter("csv")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(data)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err)

	require.Len(t, csvRecords, 3)
	headers := csvRecords[0]
	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[h] = i
	}
	assert.Contains(t, headers, "name", "Headers should contain 'name'")
	assert.Contains(t, headers, "email", "Headers should contain 'email'")
	assert.Equal(t, "Alice", csvRecords[1][headerIndex["name"]])
	assert.Equal(t, "alice@example.com", csvRecords[1][headerIndex["email"]])
	assert.Equal(t, "Bob", csvRecords[2][headerIndex["name"]])
	assert.Equal(t, "bob@example.com", csvRecords[2][headerIndex["email"]])
}

func TestGenerateCSVDataset_ValidCSVOutput(t *testing.T) {
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
		Format:      "csv",
		Count:       2,
		Schema: models.SchemaField{
			Name: "test",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)

	formatter, err := generator.GetFormatter("csv")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(data)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err, "Output should be valid CSV")
	assert.Len(t, csvRecords, 3, "Should have header + 2 data rows")
}

func TestGenerateCSV_WithNumericFields(t *testing.T) {
	mock := &mockLLMClient{
		response: `[
			{"product": "Widget A", "price": 19.99, "quantity": 100},
			{"product": "Widget B", "price": 29.99, "quantity": 50}
		]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Products",
		Description: "Products",
		Format:      "csv",
		Count:       2,
		Schema: models.SchemaField{
			Name: "product",
			Type: models.FieldTypeString,
		},
	}

	records, err := gen.Generate(descFile)
	require.NoError(t, err)
	require.Len(t, records, 2)

	formatter, err := generator.GetFormatter("csv")
	require.NoError(t, err)

	data, err := formatter.Format(records)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(data)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err)

	require.Len(t, csvRecords, 3)
	headerIndex := getHeaderIndex(csvRecords[0])
	assert.Equal(t, "Widget A", csvRecords[1][headerIndex["product"]])
	assert.Equal(t, "19.99", csvRecords[1][headerIndex["price"]])
	assert.Equal(t, "100", csvRecords[1][headerIndex["quantity"]])
	assert.Equal(t, "Widget B", csvRecords[2][headerIndex["product"]])
	assert.Equal(t, "29.99", csvRecords[2][headerIndex["price"]])
	assert.Equal(t, "50", csvRecords[2][headerIndex["quantity"]])
}
