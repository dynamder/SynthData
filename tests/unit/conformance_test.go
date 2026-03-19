package unit

import (
	"testing"

	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConformance_ValidRecords(t *testing.T) {
	schema := &models.SchemaField{
		Name:     "user",
		Type:     models.FieldTypeNested,
		Required: true,
		Children: []models.SchemaField{
			{Name: "id", Type: models.FieldTypeInteger, Required: true},
			{Name: "name", Type: models.FieldTypeString, Required: true},
			{Name: "email", Type: models.FieldTypeEmail, Required: true},
			{Name: "age", Type: models.FieldTypeInteger, Required: false},
		},
	}

	tests := []struct {
		name  string
		data  map[string]interface{}
		want  bool
		score float64
	}{
		{
			name: "all valid fields",
			data: map[string]interface{}{
				"id":    1,
				"name":  "John",
				"email": "john@example.com",
				"age":   25,
			},
			want:  true,
			score: 1.0,
		},
		{
			name: "missing optional field",
			data: map[string]interface{}{
				"id":    1,
				"name":  "John",
				"email": "john@example.com",
			},
			want:  false,
			score: 0.75,
		},
		{
			name: "one invalid type",
			data: map[string]interface{}{
				"id":    "not an int",
				"name":  "John",
				"email": "john@example.com",
				"age":   25,
			},
			want:  false,
			score: 0.75,
		},
		{
			name: "missing required field",
			data: map[string]interface{}{
				"id":   1,
				"name": "John",
			},
			want:  false,
			score: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.ValidateSchemaConformance(schema, tt.data)
			assert.Equal(t, tt.want, result.Valid)
			assert.Equal(t, tt.score, result.ActualScore)
		})
	}
}

func TestConformance_Threshold95(t *testing.T) {
	schema := &models.SchemaField{
		Name:     "user",
		Type:     models.FieldTypeNested,
		Required: true,
		Children: []models.SchemaField{
			{Name: "id", Type: models.FieldTypeInteger, Required: true},
			{Name: "name", Type: models.FieldTypeString, Required: true},
			{Name: "email", Type: models.FieldTypeEmail, Required: true},
			{Name: "age", Type: models.FieldTypeInteger, Required: false},
			{Name: "country", Type: models.FieldTypeString, Required: false},
			{Name: "phone", Type: models.FieldTypeString, Required: false},
			{Name: "city", Type: models.FieldTypeString, Required: false},
			{Name: "state", Type: models.FieldTypeString, Required: false},
			{Name: "zip", Type: models.FieldTypeString, Required: false},
			{Name: "status", Type: models.FieldTypeString, Required: false},
			{Name: "role", Type: models.FieldTypeString, Required: false},
			{Name: "department", Type: models.FieldTypeString, Required: false},
			{Name: "manager", Type: models.FieldTypeString, Required: false},
			{Name: "start_date", Type: models.FieldTypeDate, Required: false},
			{Name: "end_date", Type: models.FieldTypeDate, Required: false},
			{Name: "salary", Type: models.FieldTypeFloat, Required: false},
			{Name: "bonus", Type: models.FieldTypeFloat, Required: false},
			{Name: "active", Type: models.FieldTypeBoolean, Required: false},
			{Name: "verified", Type: models.FieldTypeBoolean, Required: false},
			{Name: "created_at", Type: models.FieldTypeDate, Required: false},
			{Name: "updated_at", Type: models.FieldTypeDate, Required: false},
		},
	}

	data := map[string]interface{}{
		"id":         1,
		"name":       "John",
		"email":      "john@example.com",
		"age":        25,
		"country":    "USA",
		"phone":      "123-456-7890",
		"city":       "NYC",
		"state":      "NY",
		"zip":        "10001",
		"status":     "active",
		"role":       "user",
		"department": "Engineering",
		"manager":    "Jane",
		"start_date": "2020-01-01",
		"salary":     100000.0,
		"bonus":      10000.0,
		"active":     true,
		"verified":   true,
		"created_at": "2020-01-01",
		"updated_at": "2024-01-01",
	}

	result := validator.ValidateSchemaConformance(schema, data)
	require.True(t, result.Valid)
	assert.GreaterOrEqual(t, result.ActualScore, 0.95)
	assert.Equal(t, 21, result.TotalFields)
}

func TestConformance_BelowThreshold(t *testing.T) {
	schema := &models.SchemaField{
		Name:     "user",
		Type:     models.FieldTypeNested,
		Required: true,
		Children: []models.SchemaField{
			{Name: "id", Type: models.FieldTypeInteger, Required: true},
			{Name: "name", Type: models.FieldTypeString, Required: true},
			{Name: "email", Type: models.FieldTypeEmail, Required: true},
			{Name: "age", Type: models.FieldTypeInteger, Required: false},
			{Name: "country", Type: models.FieldTypeString, Required: false},
			{Name: "phone", Type: models.FieldTypeString, Required: false},
			{Name: "address", Type: models.FieldTypeString, Required: false},
			{Name: "city", Type: models.FieldTypeString, Required: false},
			{Name: "state", Type: models.FieldTypeString, Required: false},
			{Name: "zip", Type: models.FieldTypeString, Required: false},
			{Name: "status", Type: models.FieldTypeString, Required: false},
			{Name: "role", Type: models.FieldTypeString, Required: false},
			{Name: "department", Type: models.FieldTypeString, Required: false},
			{Name: "manager", Type: models.FieldTypeString, Required: false},
			{Name: "start_date", Type: models.FieldTypeDate, Required: false},
			{Name: "end_date", Type: models.FieldTypeDate, Required: false},
			{Name: "salary", Type: models.FieldTypeFloat, Required: false},
			{Name: "bonus", Type: models.FieldTypeFloat, Required: false},
			{Name: "active", Type: models.FieldTypeBoolean, Required: false},
			{Name: "verified", Type: models.FieldTypeBoolean, Required: false},
		},
	}

	data := map[string]interface{}{
		"id":    1,
		"name":  "John",
		"email": "john@example.com",
	}

	result := validator.ValidateSchemaConformance(schema, data)
	require.False(t, result.Valid)
	assert.Less(t, result.ActualScore, 0.95)
	assert.Equal(t, 20, result.TotalFields)
}

func TestConformance_NilSchema(t *testing.T) {
	result := validator.ValidateSchemaConformance(nil, map[string]interface{}{})
	require.False(t, result.Valid)
	assert.Contains(t, result.InvalidFields[0], "nil")
}

func TestConformance_NilRecord(t *testing.T) {
	schema := &models.SchemaField{
		Name: "user",
		Type: models.FieldTypeString,
	}

	result := validator.ValidateRecordConformance(schema, nil)
	require.False(t, result.Valid)
}
