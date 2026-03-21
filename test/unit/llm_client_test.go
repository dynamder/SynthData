package unit

import (
	"testing"

	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/generator"
	"github.com/dynamder/synthdata/internal/services/llm"
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

func (m *mockLLMClient) GenerateWithBatchSize(prompt string, batchSize int) (string, error) {
	return m.Generate(prompt)
}

func TestLLMClientInterface(t *testing.T) {
	var _ llm.Client = &mockLLMClient{}
}

func TestMockLLMClient_Generate(t *testing.T) {
	mock := &mockLLMClient{
		response: `[{"name": "Test User", "age": 25}]`,
	}

	result, err := mock.Generate("test prompt")
	require.NoError(t, err)
	assert.Equal(t, `[{"name": "Test User", "age": 25}]`, result)
}

func TestMockLLMClient_GenerateError(t *testing.T) {
	mock := &mockLLMClient{
		err: assert.AnError,
	}

	result, err := mock.Generate("test prompt")
	require.Error(t, err)
	assert.Empty(t, result)
}

func TestGeneratorWithMockClient(t *testing.T) {
	mock := &mockLLMClient{
		response: `[{"name": "Test User", "age": 25}]`,
	}

	gen := generator.New(mock)

	descFile := &models.DescriptionFile{
		Name:        "Test Dataset",
		Description: "A test dataset",
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
	assert.Equal(t, "Test User", records[0]["name"])
}
