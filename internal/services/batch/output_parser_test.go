package batch

import (
	"testing"
)

func TestOutputParser_Parse(t *testing.T) {
	parser := NewOutputParser()

	tests := []struct {
		name        string
		input       string
		wantValid   int
		wantInvalid bool
	}{
		{"valid json array", `[{"id": 1}, {"id": 2}]`, 2, false},
		{"empty array", `[]`, 0, false},
		{"invalid json", `not json`, 0, true},
		{"single record wrapped in array", `[{"id": 1}]`, 1, false},
		{"json with whitespace", `  [{"id": 1}]  `, 1, false},
		{"json with markdown", "```json\n[{\"id\": 1}]\n```", 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, invalid, _ := parser.Parse(tt.input)
			if len(valid) != tt.wantValid {
				t.Errorf("Parse() valid = %d, want %d", len(valid), tt.wantValid)
			}
			if tt.wantInvalid && len(invalid) == 0 {
				t.Error("expected invalid record on parse error")
			}
		})
	}
}

func TestOutputParser_IsEmpty(t *testing.T) {
	parser := NewOutputParser()

	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\n\t", true},
		{"{}", false},
		{"[]", false},
		{"valid content", false},
	}

	for _, tt := range tests {
		result := parser.IsEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}
