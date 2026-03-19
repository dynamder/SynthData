package validator

import (
	"errors"
	"fmt"

	"github.com/anomalyco/synthdata/internal/models"
)

var (
	ErrEmptyName        = errors.New("name is required")
	ErrInvalidFormat    = errors.New("format must be json or csv")
	ErrInvalidCount     = errors.New("count must be between 1 and 10000")
	ErrInvalidFieldType = errors.New("invalid field type")
	ErrEmptySchema      = errors.New("schema is required")
	ErrEmptySchemaName  = errors.New("schema name is required")
	ErrEmptySchemaType  = errors.New("schema type is required")
)

type ValidationError struct {
	Field   string
	Message string
	Line    int
}

func (e *ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s - %s", e.Line, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateDescription(desc *models.DescriptionFile) error {
	var errs []error

	if desc.Name == "" {
		errs = append(errs, ErrEmptyName)
	}

	if desc.Format != "json" && desc.Format != "csv" {
		errs = append(errs, ErrInvalidFormat)
	}

	if desc.Count < 1 || desc.Count > 10000 {
		errs = append(errs, ErrInvalidCount)
	}

	if desc.Schema.Name == "" {
		errs = append(errs, ErrEmptySchemaName)
	}

	if desc.Schema.Type == "" {
		errs = append(errs, ErrEmptySchemaType)
	}

	if err := validateFieldType(desc.Schema.Type); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func validateFieldType(fieldType models.FieldType) error {
	validTypes := map[models.FieldType]bool{
		models.FieldTypeString:  true,
		models.FieldTypeInteger: true,
		models.FieldTypeFloat:   true,
		models.FieldTypeBoolean: true,
		models.FieldTypeEmail:   true,
		models.FieldTypeDate:    true,
		models.FieldTypeUUID:    true,
		models.FieldTypeNested:  true,
	}

	if !validTypes[fieldType] {
		return fmt.Errorf("%w: %s", ErrInvalidFieldType, fieldType)
	}

	return nil
}
