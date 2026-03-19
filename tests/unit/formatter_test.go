package unit

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/anomalyco/synthdata/internal/services/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFormatter_UnsupportedFormat(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"xml", "xml"},
		{"yaml", "yaml"},
		{"txt", "txt"},
		{"unknown", "unknown"},
		{"", ""},
		{"JSON", "JSON"},
		{"CSV", "CSV"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter, err := generator.GetFormatter(tt.format)
			if tt.format == "json" || tt.format == "JSON" || tt.format == "csv" || tt.format == "CSV" {
				require.NoError(t, err, "json and csv should be supported")
				require.NotNil(t, formatter)
				return
			}
			assert.Error(t, err, "unsupported format should return error")
			assert.Nil(t, formatter)
			if err != nil {
				assert.Contains(t, err.Error(), "unsupported format", "Error message should mention unsupported format")
			}
		})
	}
}

func TestJSONFormatter_ParseValidation(t *testing.T) {
	f, err := generator.GetFormatter("json")
	require.NoError(t, err)
	require.NotNil(t, f)

	records := []map[string]interface{}{
		{"name": "Alice", "age": 25, "email": "alice@example.com"},
		{"name": "Bob", "age": 30, "email": "bob@example.com"},
	}

	data, err := f.Format(records)
	require.NoError(t, err, "Format should not return error")

	var parsed []map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err, "Output should be valid parseable JSON")

	assert.Len(t, parsed, 2, "Should have 2 records")
	assert.Equal(t, "Alice", parsed[0]["name"])
	assert.Equal(t, "Bob", parsed[1]["name"])
}

func TestCSVFormatter_ParseValidation(t *testing.T) {
	f, err := generator.GetFormatter("csv")
	require.NoError(t, err)
	require.NotNil(t, f)

	records := []map[string]interface{}{
		{"name": "Alice", "age": 25},
		{"name": "Bob", "age": 30},
	}

	data, err := f.Format(records)
	require.NoError(t, err, "Format should not return error")

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	require.GreaterOrEqual(t, len(lines), 2, "Should have at least header + 1 data row")

	headers := strings.Split(lines[0], ",")
	assert.Contains(t, headers, "name", "Headers should contain 'name'")
	assert.Contains(t, headers, "age", "Headers should contain 'age'")

	dataLine := strings.Split(lines[1], ",")
	assert.Len(t, dataLine, 2, "Data row should have 2 columns")
}

func TestCSVFormatter_ProperDelimiter(t *testing.T) {
	f, err := generator.GetFormatter("csv")
	require.NoError(t, err)

	records := []map[string]interface{}{
		{"name": "Test User", "value": "100"},
	}

	data, err := f.Format(records)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	require.GreaterOrEqual(t, len(lines), 2, "Should have header + data row")

	headers := strings.Split(lines[0], ",")
	assert.Contains(t, headers, "name")
	assert.Contains(t, headers, "value")
}
