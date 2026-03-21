package formatters

import (
	"testing"
)

func TestJSONFormatter_Format(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name           string
		records        []map[string]interface{}
		expectEmptyArr bool
	}{
		{
			name:           "empty records",
			records:        []map[string]interface{}{},
			expectEmptyArr: true,
		},
		{
			name: "single record",
			records: []map[string]interface{}{
				{"name": "test", "age": 25},
			},
			expectEmptyArr: false,
		},
		{
			name: "multiple records",
			records: []map[string]interface{}{
				{"name": "test1", "age": 25},
				{"name": "test2", "age": 30},
			},
			expectEmptyArr: false,
		},
		{
			name: "nested record",
			records: []map[string]interface{}{
				{"name": "test", "address": map[string]string{"city": "NYC"}},
			},
			expectEmptyArr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := formatter.Format(tt.records)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expectEmptyArr {
				if string(data) != "[]" {
					t.Errorf("expected empty array, got %s", string(data))
				}
			}
			if len(data) == 0 && !tt.expectEmptyArr {
				t.Error("expected non-empty output")
			}
		})
	}
}

func TestJSONFormatter_FormatSingle(t *testing.T) {
	formatter := NewJSONFormatter()

	record := map[string]interface{}{
		"name": "test",
		"age":  25,
	}

	data, err := formatter.FormatSingle(record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output")
	}
}

func TestJSONFormatter_FormatSingle_Nil(t *testing.T) {
	formatter := NewJSONFormatter()

	data, err := formatter.FormatSingle(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Error("expected non-empty output for nil record")
	}
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
		{"", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ParseJSONFieldType(tt.input)
			if string(got) != tt.expected {
				t.Errorf("ParseJSONFieldType(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseJSONFieldType_CaseInsensitive(t *testing.T) {
	tests := []string{"String", "INTEGER", "FLOAT", "BOOLEAN", "EMAIL"}

	for _, tt := range tests {
		got := ParseJSONFieldType(tt)
		if got == "" {
			t.Errorf("ParseJSONFieldType(%q) returned empty", tt)
		}
	}
}
