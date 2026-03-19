package models

type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeInteger FieldType = "integer"
	FieldTypeFloat   FieldType = "float"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeEmail   FieldType = "email"
	FieldTypeDate    FieldType = "date"
	FieldTypeUUID    FieldType = "uuid"
	FieldTypeNested  FieldType = "nested"
)

type SchemaField struct {
	Name     string            `json:"name"`
	Type     FieldType         `json:"type"`
	Required bool              `json:"required"`
	Options  map[string]string `json:"options,omitempty"`
	Children []SchemaField     `json:"children,omitempty"`
}

type DescriptionFile struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Format      string      `json:"format"`
	Count       int         `json:"count"`
	Schema      SchemaField `json:"schema"`
}
