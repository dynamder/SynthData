package formatters

import (
	"testing"
)

func TestCSVFormatter_Format(t *testing.T) {
	formatter := NewCSVFormatter()

	tests := []struct {
		name    string
		records []map[string]interface{}
		wantLen int
	}{
		{
			name:    "empty records",
			records: []map[string]interface{}{},
			wantLen: 0,
		},
		{
			name: "single record",
			records: []map[string]interface{}{
				{"name": "test", "age": "25"},
			},
			wantLen: 1,
		},
		{
			name: "multiple records",
			records: []map[string]interface{}{
				{"name": "test1", "age": "25"},
				{"name": "test2", "age": "30"},
			},
			wantLen: 2,
		},
		{
			name: "records with different keys",
			records: []map[string]interface{}{
				{"name": "test1"},
				{"age": "25"},
			},
			wantLen: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := formatter.Format(tt.records)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantLen == 0 {
				if len(data) != 0 {
					t.Errorf("expected empty output, got %s", string(data))
				}
				return
			}

			if len(data) == 0 {
				t.Fatal("expected non-empty output")
			}
		})
	}
}

func TestCSVFormatter_Format_Headers(t *testing.T) {
	formatter := NewCSVFormatter()

	records := []map[string]interface{}{
		{"name": "test1", "age": "25"},
		{"name": "test2", "city": "NYC"},
	}

	data, err := formatter.Format(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestCSVFormatter_Format_EmptyValue(t *testing.T) {
	formatter := NewCSVFormatter()

	records := []map[string]interface{}{
		{"name": "test", "age": ""},
	}

	data, err := formatter.Format(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestCSVFormatter_Format_NestedRecord(t *testing.T) {
	formatter := NewCSVFormatter()

	records := []map[string]interface{}{
		{"name": "test", "data": map[string]string{"key": "value"}},
	}

	data, err := formatter.Format(records)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output for nested record")
	}
}

func TestParseCSVFieldType(t *testing.T) {
	result := ParseCSVFieldType("string")
	if string(result) != "string" {
		t.Errorf("ParseCSVFieldType(\"string\") = %q, want \"string\"", result)
	}

	result = ParseCSVFieldType("integer")
	if string(result) != "integer" {
		t.Errorf("ParseCSVFieldType(\"integer\") = %q, want \"integer\"", result)
	}
}
