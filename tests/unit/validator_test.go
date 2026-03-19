package unit

import (
	"testing"

	"github.com/anomalyco/synthdata/internal/models"
	"github.com/anomalyco/synthdata/internal/services/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_InvalidSyntax(t *testing.T) {
	tests := []struct {
		name    string
		desc    models.DescriptionFile
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty name",
			desc: models.DescriptionFile{
				Name:        "",
				Description: "Test dataset",
				Format:      "json",
				Count:       10,
				Schema: models.SchemaField{
					Name: "user",
					Type: models.FieldTypeString,
				},
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid format",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "xml",
				Count:       10,
				Schema: models.SchemaField{
					Name: "user",
					Type: models.FieldTypeString,
				},
			},
			wantErr: true,
			errMsg:  "format must be json or csv",
		},
		{
			name: "invalid count zero",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       0,
				Schema: models.SchemaField{
					Name: "user",
					Type: models.FieldTypeString,
				},
			},
			wantErr: true,
			errMsg:  "count must be between 1 and 10000",
		},
		{
			name: "invalid count exceeds max",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       50000,
				Schema: models.SchemaField{
					Name: "user",
					Type: models.FieldTypeString,
				},
			},
			wantErr: true,
			errMsg:  "count must be between 1 and 10000",
		},
		{
			name: "invalid field type",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       10,
				Schema: models.SchemaField{
					Name: "user",
					Type: "invalid_type",
				},
			},
			wantErr: true,
			errMsg:  "invalid field type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDescription(&tt.desc)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_MissingRequiredFields(t *testing.T) {
	tests := []struct {
		name    string
		desc    models.DescriptionFile
		wantErr bool
		errMsg  string
	}{
		{
			name: "missing schema",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       10,
				Schema:      models.SchemaField{},
			},
			wantErr: true,
			errMsg:  "schema name is required",
		},
		{
			name: "missing schema name",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       10,
				Schema: models.SchemaField{
					Name: "",
					Type: models.FieldTypeString,
				},
			},
			wantErr: true,
			errMsg:  "schema name is required",
		},
		{
			name: "missing schema type",
			desc: models.DescriptionFile{
				Name:        "test",
				Description: "Test dataset",
				Format:      "json",
				Count:       10,
				Schema: models.SchemaField{
					Name: "user",
					Type: "",
				},
			},
			wantErr: true,
			errMsg:  "schema type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDescription(&tt.desc)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidator_ValidDescription(t *testing.T) {
	desc := models.DescriptionFile{
		Name:        "test",
		Description: "Test dataset",
		Format:      "json",
		Count:       100,
		Schema: models.SchemaField{
			Name:     "user",
			Type:     models.FieldTypeString,
			Required: true,
		},
	}

	err := validator.ValidateDescription(&desc)
	require.NoError(t, err)
}
