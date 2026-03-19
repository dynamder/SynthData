package performance

import (
	"os"
	"testing"

	"github.com/anomalyco/synthdata/internal/config"
	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/generator"
	"github.com/anomalyco/synthdata/internal/services/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate10kRecordsWithoutCrash(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping scale test")
	}

	config.LoadConfig()
	client := llm.NewOpenAIClient()
	gen := generator.New(client)

	desc := &models.DescriptionFile{
		Name:  "Scale Test Dataset",
		Count: 10000,
		Schema: models.SchemaField{
			Name: "root",
			Children: []models.SchemaField{
				{Name: "id", Type: models.FieldTypeInteger},
				{Name: "name", Type: models.FieldTypeString},
				{Name: "email", Type: models.FieldTypeString},
				{Name: "age", Type: models.FieldTypeInteger},
			},
		},
	}

	records, err := gen.Generate(desc)

	require.NoError(t, err, "SC-005: Should generate 10k records without crashing")
	assert.Len(t, records, 10000, "SC-005: Should return exactly 10000 records")
}
