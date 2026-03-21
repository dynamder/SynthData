package models

type DatasetConfig struct {
	Name        string `json:"name" mapstructure:"name"`
	Description string `json:"description" mapstructure:"description"`
	Format      string `json:"format" mapstructure:"format"`
	Scale       int    `json:"scale" mapstructure:"scale"`
	Schema      Schema `json:"schema" mapstructure:"schema"`
	Output      string `json:"output" mapstructure:"output"`
}

type Schema struct {
	Fields []Field `json:"fields" mapstructure:"fields"`
}

type Field struct {
	Name     string      `json:"name" mapstructure:"name"`
	Type     string      `json:"type" mapstructure:"type"`
	Required bool        `json:"required" mapstructure:"required"`
	Enum     []string    `json:"enum,omitempty" mapstructure:"enum"`
	Min      interface{} `json:"min,omitempty" mapstructure:"min"`
	Max      interface{} `json:"max,omitempty" mapstructure:"max"`
	Format   string      `json:"format,omitempty" mapstructure:"format"`
	Nested   *Schema     `json:"nested,omitempty" mapstructure:"nested"`
}
