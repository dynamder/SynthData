package formatters

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dynamder/synthdata/internal/models"
)

type JSONFormatter struct{}

func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

func (f *JSONFormatter) Format(records []map[string]interface{}) ([]byte, error) {
	if len(records) == 0 {
		return []byte("[]"), nil
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}

func (f *JSONFormatter) FormatSingle(record map[string]interface{}) ([]byte, error) {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return data, nil
}

func ParseJSONFieldType(fieldType string) models.FieldType {
	switch strings.ToLower(fieldType) {
	case "string":
		return models.FieldTypeString
	case "integer", "int":
		return models.FieldTypeInteger
	case "float", "number":
		return models.FieldTypeFloat
	case "boolean", "bool":
		return models.FieldTypeBoolean
	case "email":
		return models.FieldTypeEmail
	case "date", "datetime":
		return models.FieldTypeDate
	case "uuid":
		return models.FieldTypeUUID
	case "nested", "object":
		return models.FieldTypeNested
	default:
		return models.FieldTypeString
	}
}
