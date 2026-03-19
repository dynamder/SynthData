package validator

import (
	"fmt"
	"reflect"

	"github.com/anomalyco/synthdata/internal/models"
)

const DefaultConformanceThreshold = 0.95

type ConformanceResult struct {
	Valid         bool
	Threshold     float64
	ActualScore   float64
	TotalFields   int
	ValidFields   int
	InvalidFields []string
}

func ValidateSchemaConformance(schema *models.SchemaField, data map[string]interface{}) *ConformanceResult {
	result := &ConformanceResult{
		Threshold:     DefaultConformanceThreshold,
		TotalFields:   0,
		ValidFields:   0,
		InvalidFields: []string{},
	}

	if schema == nil {
		result.Valid = false
		result.InvalidFields = append(result.InvalidFields, "schema is nil")
		return result
	}

	if schema.Type == models.FieldTypeNested && len(schema.Children) > 0 {
		for _, child := range schema.Children {
			validateField(&child, data, result)
		}
	} else {
		validateField(schema, data, result)
	}

	if result.TotalFields > 0 {
		result.ActualScore = float64(result.ValidFields) / float64(result.TotalFields)
	} else {
		result.ActualScore = 0
	}

	result.Valid = result.ActualScore >= result.Threshold

	return result
}

func validateField(field *models.SchemaField, data map[string]interface{}, result *ConformanceResult) {
	result.TotalFields++

	value, exists := data[field.Name]
	if !exists {
		if field.Required {
			result.InvalidFields = append(result.InvalidFields, fmt.Sprintf("%s: missing required field", field.Name))
		}
		return
	}

	if !isValidType(value, field.Type) {
		result.InvalidFields = append(result.InvalidFields, fmt.Sprintf("%s: invalid type, expected %s", field.Name, field.Type))
		return
	}

	result.ValidFields++

	if field.Type == models.FieldTypeNested && len(field.Children) > 0 {
		nestedData, ok := value.(map[string]interface{})
		if ok {
			for _, child := range field.Children {
				validateField(&child, nestedData, result)
			}
		}
	}
}

func isValidType(value interface{}, expectedType models.FieldType) bool {
	if value == nil {
		return false
	}

	switch expectedType {
	case models.FieldTypeString, models.FieldTypeEmail:
		_, ok := value.(string)
		return ok
	case models.FieldTypeInteger:
		switch value.(type) {
		case int, int64, float64:
			return true
		default:
			return false
		}
	case models.FieldTypeFloat:
		switch value.(type) {
		case float64:
			return true
		default:
			return false
		}
	case models.FieldTypeBoolean:
		_, ok := value.(bool)
		return ok
	case models.FieldTypeDate:
		_, ok := value.(string)
		return ok
	case models.FieldTypeUUID:
		_, ok := value.(string)
		return ok
	case models.FieldTypeNested:
		_, ok := value.(map[string]interface{})
		return ok
	default:
		return false
	}
}

func ValidateRecordConformance(schema *models.SchemaField, record interface{}) *ConformanceResult {
	recordVal := reflect.ValueOf(record)
	if recordVal.Kind() == reflect.Map {
		data := make(map[string]interface{})
		for _, key := range recordVal.MapKeys() {
			data[key.String()] = recordVal.MapIndex(key).Interface()
		}
		return ValidateSchemaConformance(schema, data)
	}

	result := &ConformanceResult{
		Valid:         false,
		InvalidFields: []string{"record is not a map"},
	}
	return result
}
