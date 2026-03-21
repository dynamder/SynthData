package unit

import (
	"encoding/json"
	"testing"

	"github.com/dynamder/synthdata/internal/formatters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Format(t *testing.T) {
	f := formatters.NewJSONFormatter()

	tests := []struct {
		name     string
		records  []map[string]interface{}
		wantErr  bool
		validate func(*testing.T, []byte)
	}{
		{
			name:    "empty records",
			records: []map[string]interface{}{},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				assert.Equal(t, "[]", string(data))
			},
		},
		{
			name: "single record",
			records: []map[string]interface{}{
				{"name": "John", "age": 30},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var result []map[string]interface{}
				err := json.Unmarshal(data, &result)
				require.NoError(t, err)
				require.Len(t, result, 1)
				assert.Equal(t, "John", result[0]["name"])
				assert.Equal(t, float64(30), result[0]["age"])
			},
		},
		{
			name: "multiple records",
			records: []map[string]interface{}{
				{"name": "Alice", "age": 25},
				{"name": "Bob", "age": 35},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var result []map[string]interface{}
				err := json.Unmarshal(data, &result)
				require.NoError(t, err)
				require.Len(t, result, 2)
				assert.Equal(t, "Alice", result[0]["name"])
				assert.Equal(t, "Bob", result[1]["name"])
			},
		},
		{
			name: "nested objects",
			records: []map[string]interface{}{
				{
					"name":  "Test",
					"email": "test@example.com",
					"metadata": map[string]interface{}{
						"role":   "admin",
						"active": true,
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var result []map[string]interface{}
				err := json.Unmarshal(data, &result)
				require.NoError(t, err)
				metadata, ok := result[0]["metadata"].(map[string]interface{})
				require.True(t, ok, "metadata should be a nested object")
				assert.Equal(t, "admin", metadata["role"])
			},
		},
		{
			name: "all field types",
			records: []map[string]interface{}{
				{
					"string":  "hello",
					"integer": 42,
					"float":   3.14,
					"bool":    true,
					"email":   "user@test.com",
				},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				var result []map[string]interface{}
				err := json.Unmarshal(data, &result)
				require.NoError(t, err)
				assert.Equal(t, "hello", result[0]["string"])
				assert.Equal(t, float64(42), result[0]["integer"])
				assert.Equal(t, 3.14, result[0]["float"])
				assert.Equal(t, true, result[0]["bool"])
				assert.Equal(t, "user@test.com", result[0]["email"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := f.Format(tt.records)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			tt.validate(t, data)
		})
	}
}

func TestJSONFormatter_FormatSingle(t *testing.T) {
	f := formatters.NewJSONFormatter()

	record := map[string]interface{}{
		"name": "Test User",
		"age":  25,
	}

	data, err := f.FormatSingle(record)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	assert.Equal(t, "Test User", result["name"])
	assert.Equal(t, float64(25), result["age"])
}

func TestJSONFormatter_FormatSingle_Empty(t *testing.T) {
	f := formatters.NewJSONFormatter()

	record := map[string]interface{}{}

	data, err := f.FormatSingle(record)
	require.NoError(t, err)

	assert.Equal(t, "{}", string(data))
}

func TestParseJSONFieldType(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"string", "string"},
		{"integer", "integer"},
		{"int", "integer"},
		{"float", "float"},
		{"number", "float"},
		{"boolean", "boolean"},
		{"bool", "boolean"},
		{"email", "email"},
		{"date", "date"},
		{"datetime", "date"},
		{"uuid", "uuid"},
		{"nested", "nested"},
		{"object", "nested"},
		{"unknown", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatters.ParseJSONFieldType(tt.input)
			assert.Equal(t, tt.expected, string(result))
		})
	}
}

func TestJSONFormatter_IsValidOutput(t *testing.T) {
	f := formatters.NewJSONFormatter()

	records := []map[string]interface{}{
		{"name": "test"},
	}

	data, err := f.Format(records)
	require.NoError(t, err)

	var parsed []map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Output should be valid parseable JSON")

	assert.Len(t, parsed, 1)
}
