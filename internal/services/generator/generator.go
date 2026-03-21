package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	synthdatalog "github.com/dynamder/synthdata/internal"
	syntherror "github.com/dynamder/synthdata/internal/errors"
	"github.com/dynamder/synthdata/internal/formatters"
	"github.com/dynamder/synthdata/internal/models"
	"github.com/dynamder/synthdata/internal/services/llm"
)

type Generator struct {
	client llm.Client
}

func New(client llm.Client) *Generator {
	return &Generator{client: client}
}

func (g *Generator) Generate(descFile *models.DescriptionFile) ([]map[string]interface{}, error) {
	logger := synthdatalog.GetLogger()

	schemaJSON, err := json.Marshal(descFile.Schema)
	if err != nil {
		//return nil, fmt.Errorf("failed to marshal schema: %w", err)
		return nil, syntherror.Wrap(syntherror.CodeInvalidDescription, "failed to marshal schema.", err)
	}

	prompt := buildPrompt(descFile.Name, descFile.Description, string(schemaJSON), descFile.Count)
	response, err := g.client.Generate(prompt)
	if err != nil {
		logger.Error("LLM call failed", map[string]interface{}{
			"record_count": descFile.Count,
			"error":        err.Error(),
			"prompt":       prompt,
			"response":     response,
		})
		//return nil, fmt.Errorf("failed to generate from LLM: %w", err)
		return nil, err
	}

	var records []map[string]interface{}
	cleaned := cleanJSONResponse(response)
	if err := json.Unmarshal([]byte(cleaned), &records); err != nil {
		logger.Error("Failed to parse LLM response", map[string]interface{}{
			"record_count": descFile.Count,
			"error":        err.Error(),
			"response":     response,
		})
		//return nil, fmt.Errorf("failed to parse LLM response: %w", err)
		return nil,
			syntherror.Wrap(
				syntherror.ErrLLMCall.Code,
				"failed to parse LLM response",
				llm.NewLLmCallError(llm.GenerationFailed, err),
			)
	}

	logger.Info("Generation completed successfully", map[string]interface{}{
		"record_count": len(records),
	})

	return records, nil
}

func buildPrompt(name, description, schema string, count int) string {
	return fmt.Sprintf(`Generate exactly %d records for a dataset named "%s".
Description: %s
Schema: %s

Return ONLY a valid JSON array of objects with no additional text or explanation.
Each object must conform to the schema provided.
`, count, name, description, schema)
}

// TODO: Is this correct?
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
