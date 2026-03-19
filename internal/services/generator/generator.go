package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anomalyco/synthdata/internal/formatters"
	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/llm"
)

type Generator struct {
	client llm.Client
}

func New(client llm.Client) *Generator {
	return &Generator{client: client}
}

func (g *Generator) Generate(descFile *models.DescriptionFile) ([]map[string]interface{}, error) {
	schemaJSON, err := json.Marshal(descFile.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal schema: %w", err)
	}

	prompt := buildPrompt(descFile.Name, descFile.Description, string(schemaJSON), descFile.Count)
	response, err := g.client.Generate(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate from LLM: %w", err)
	}

	var records []map[string]interface{}
	cleaned := cleanJSONResponse(response)
	if err := json.Unmarshal([]byte(cleaned), &records); err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	return records, nil
}

func buildPrompt(name, description, schema string, count int) string {
	return fmt.Sprintf(`Generate %d records for a dataset named "%s".
Description: %s
Schema: %s

Return ONLY a valid JSON array of objects with no additional text or explanation.
Each object must conform to the schema provided.
`, count, name, description, schema)
}

func cleanJSONResponse(response string) string {
	response = strings.TrimSpace(response)
	response = strings.Trim(response, "`")
	if idx := strings.Index(response, "["); idx > 0 {
		response = response[idx:]
	}
	if idx := strings.LastIndex(response, "]"); idx >= 0 && idx < len(response)-1 {
		response = response[:idx+1]
	}
	return response
}

type Formatter interface {
	Format(records []map[string]interface{}) ([]byte, error)
}

func GetFormatter(format string) (Formatter, error) {
	switch strings.ToLower(format) {
	case "json":
		return formatters.NewJSONFormatter(), nil
	case "csv":
		return formatters.NewCSVFormatter(), nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
