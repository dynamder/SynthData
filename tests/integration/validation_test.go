package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/anomalyco/synthdata/internal/services/parser"
	"github.com/anomalyco/synthdata/internal/services/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidation_Integration(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid description file",
			content: `{
				"name": "test_users",
				"description": "Test user dataset",
				"format": "json",
				"count": 100,
				"schema": {
					"name": "user",
					"type": "string",
					"required": true
				}
			}`,
			wantErr: false,
		},
		{
			name: "valid description file with description_file",
			content: `{
				"name": "test_users",
				"description_file": "test.md",
				"format": "json",
				"count": 100,
				"schema": {
					"name": "user",
					"type": "string",
					"required": true
				}
			}`,
			wantErr:     true,
			errContains: "description file not found",
		},
		{
			name: "missing required name",
			content: `{
				"name": "",
				"description": "Test dataset",
				"format": "json",
				"count": 100,
				"schema": {
					"name": "user",
					"type": "string"
				}
			}`,
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name: "invalid format",
			content: `{
				"name": "test",
				"description": "Test dataset",
				"format": "xml",
				"count": 100,
				"schema": {
					"name": "user",
					"type": "string"
				}
			}`,
			wantErr:     true,
			errContains: "format must be json or csv",
		},
		{
			name: "count exceeds maximum",
			content: `{
				"name": "test",
				"description": "Test dataset",
				"format": "json",
				"count": 50000,
				"schema": {
					"name": "user",
					"type": "string"
				}
			}`,
			wantErr:     true,
			errContains: "count must be between 1 and 10000",
		},
		{
			name: "missing schema name",
			content: `{
				"name": "test",
				"description": "Test dataset",
				"format": "json",
				"count": 100,
				"schema": {
					"name": "",
					"type": "string"
				}
			}`,
			wantErr:     true,
			errContains: "schema name is required",
		},
		{
			name:        "invalid JSON syntax",
			content:     `{ invalid json }`,
			wantErr:     true,
			errContains: "parse",
		},
		{
			name:        "missing both description and description_file",
			content:     `{"name": "test", "format": "json", "count": 100, "schema": {"name": "user", "type": "string"}}`,
			wantErr:     true,
			errContains: "either 'description' or 'description_file' must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			descFile := filepath.Join(tmpDir, "description.json")
			err := os.WriteFile(descFile, []byte(tt.content), 0644)
			require.NoError(t, err)

			parsed, err := parser.ParseDescriptionFile(descFile)
			if tt.wantErr {
				if err != nil {
					assert.Contains(t, err.Error(), tt.errContains)
					return
				}
				err = validator.ValidateDescription(parsed)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
				err = validator.ValidateDescription(parsed)
				require.NoError(t, err)
			}
		})
	}
}

func TestValidation_ErrorMessages(t *testing.T) {
	content := `{
		"name": "",
		"description": "Generate user data",
		"format": "invalid",
		"count": 0,
		"schema": {
			"name": "",
			"type": ""
		}
	}`

	tmpDir := t.TempDir()
	descFile := filepath.Join(tmpDir, "description.json")
	err := os.WriteFile(descFile, []byte(content), 0644)
	require.NoError(t, err)

	parsed, parseErr := parser.ParseDescriptionFile(descFile)
	require.NoError(t, parseErr)

	err = validator.ValidateDescription(parsed)
	require.Error(t, err)

	errMsg := err.Error()
	assert.Contains(t, errMsg, "name is required")
}
