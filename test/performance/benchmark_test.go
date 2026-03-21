package performance

import (
	"os"
	"testing"
	"time"

	"github.com/dynamder/synthdata/internal/config"
	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/generator"
	"github.com/dynamder/synthdata/internal/services/llm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkGenerate1000Records(b *testing.B) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		b.Skip("OPENAI_API_KEY not set, skipping performance test")
	}

	config.LoadConfig()
	client := llm.NewOpenAIClient()
	gen := generator.New(client)

	desc := &models.DescriptionFile{
		Name:  "Benchmark Dataset",
		Count: 1000,
		Schema: models.SchemaField{
			Name: "root",
			Children: []models.SchemaField{
				{Name: "id", Type: models.FieldTypeInteger},
				{Name: "name", Type: models.FieldTypeString},
				{Name: "email", Type: models.FieldTypeString},
			},
		},
	}

	start := time.Now()
	records, err := gen.Generate(desc)
	elapsed := time.Since(start)

	require.NoError(b, err)
	assert.Len(b, records, 1000)

	b.ReportMetric(float64(elapsed.Milliseconds()), "ms")
	assert.Less(b, elapsed, 30*time.Second, "SC-001: 1000 records should generate in under 30 seconds")
}

func TestGenerate1000RecordsUnder30Seconds(t *testing.T) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("OPENAI_API_KEY not set, skipping performance test")
	}

	config.LoadConfig()
	client := llm.NewOpenAIClient()
	gen := generator.New(client)

	desc := &models.DescriptionFile{
		Name:  "Performance Test Dataset",
		Count: 1000,
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

	start := time.Now()
	records, err := gen.Generate(desc)
	elapsed := time.Since(start)

	require.NoError(t, err)
	assert.Len(t, records, 1000)
	assert.Less(t, elapsed, 30*time.Second, "SC-001: 1000 records should generate in under 30 seconds")
}
