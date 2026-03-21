package unit

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/dynamder/synthdata/internal/formatters"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getHeaderIndex(headers []string) map[string]int {
	index := make(map[string]int)
	for i, h := range headers {
		index[h] = i
	}
	return index
}

func TestCSVFormatter_Format(t *testing.T) {
	f := formatters.NewCSVFormatter()

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
				assert.Empty(t, string(data))
			},
		},
		{
			name: "single record",
			records: []map[string]interface{}{
				{"name": "John", "age": "30"},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				r := csv.NewReader(strings.NewReader(string(data)))
				records, err := r.ReadAll()
				require.NoError(t, err)
				require.Len(t, records, 2)
				headerIndex := getHeaderIndex(records[0])
				assert.Equal(t, "John", records[1][headerIndex["name"]])
				assert.Equal(t, "30", records[1][headerIndex["age"]])
			},
		},
		{
			name: "multiple records",
			records: []map[string]interface{}{
				{"name": "Alice", "age": "25"},
				{"name": "Bob", "age": "35"},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				r := csv.NewReader(strings.NewReader(string(data)))
				records, err := r.ReadAll()
				require.NoError(t, err)
				require.Len(t, records, 3)
				headerIndex := getHeaderIndex(records[0])
				assert.Equal(t, "Alice", records[1][headerIndex["name"]])
				assert.Equal(t, "25", records[1][headerIndex["age"]])
				assert.Equal(t, "Bob", records[2][headerIndex["name"]])
				assert.Equal(t, "35", records[2][headerIndex["age"]])
			},
		},
		{
			name: "different fields per record",
			records: []map[string]interface{}{
				{"name": "Alice"},
				{"name": "Bob", "age": "35"},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				r := csv.NewReader(strings.NewReader(string(data)))
				records, err := r.ReadAll()
				require.NoError(t, err)
				require.Len(t, records, 3)
				headerIndex := getHeaderIndex(records[0])
				assert.Equal(t, "Alice", records[1][headerIndex["name"]])
				assert.Equal(t, "", records[1][headerIndex["age"]])
				assert.Equal(t, "Bob", records[2][headerIndex["name"]])
				assert.Equal(t, "35", records[2][headerIndex["age"]])
			},
		},
		{
			name: "handles nil values",
			records: []map[string]interface{}{
				{"name": "John", "email": nil},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				r := csv.NewReader(strings.NewReader(string(data)))
				records, err := r.ReadAll()
				require.NoError(t, err)
				require.Len(t, records, 2)
				headerIndex := getHeaderIndex(records[0])
				assert.Equal(t, "John", records[1][headerIndex["name"]])
				assert.Equal(t, "", records[1][headerIndex["email"]])
			},
		},
		{
			name: "handles numeric values",
			records: []map[string]interface{}{
				{"count": 100, "rate": 3.5},
			},
			wantErr: false,
			validate: func(t *testing.T, data []byte) {
				r := csv.NewReader(strings.NewReader(string(data)))
				records, err := r.ReadAll()
				require.NoError(t, err)
				require.Len(t, records, 2)
				headerIndex := getHeaderIndex(records[0])
				assert.Equal(t, "100", records[1][headerIndex["count"]])
				assert.Equal(t, "3.5", records[1][headerIndex["rate"]])
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

func TestCSVFormatter_IncludesHeaders(t *testing.T) {
	f := formatters.NewCSVFormatter()

	records := []map[string]interface{}{
		{"name": "John", "age": "30"},
	}

	data, err := f.Format(records)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(data)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err)

	require.Len(t, csvRecords, 2)
	headers := csvRecords[0]
	assert.Contains(t, headers, "name", "Headers should contain 'name'")
	assert.Contains(t, headers, "age", "Headers should contain 'age'")
}

func TestCSVFormatter_ProperlyFormatted(t *testing.T) {
	f := formatters.NewCSVFormatter()

	records := []map[string]interface{}{
		{"name": "John", "age": "30"},
		{"name": "Alice", "age": "25"},
	}

	data, err := f.Format(records)
	require.NoError(t, err)

	r := csv.NewReader(strings.NewReader(string(data)))
	csvRecords, err := r.ReadAll()
	require.NoError(t, err, "Output should be valid parseable CSV")

	assert.Len(t, csvRecords, 3, "Should have header + 2 data rows")
}

func TestParseCSVFieldType(t *testing.T) {
	result := formatters.ParseCSVFieldType("string")
	assert.Equal(t, "string", string(result))

	result = formatters.ParseCSVFieldType("integer")
	assert.Equal(t, "integer", string(result))

	result = formatters.ParseCSVFieldType("unknown")
	assert.Equal(t, "string", string(result))
}
