package prompt

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ArgumentMetadata struct {
	Name         string
	Short        string
	Description  string
	DefaultValue interface{}
	IsRequired   bool
	Type         string
	ValidValues  []string
}

func ExtractArgumentMetadata(cmd *cobra.Command) ([]ArgumentMetadata, error) {
	var metadata []ArgumentMetadata

	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}

		argMeta := ArgumentMetadata{
			Name:        flag.Name,
			Short:       flag.Shorthand,
			Description: flag.Usage,
			IsRequired:  isRequiredFlag(flag, cmd),
		}

		defaultVal := flag.DefValue
		if defaultVal != "" && defaultVal != "false" && defaultVal != "0" {
			argMeta.DefaultValue = parseDefaultValue(defaultVal)
		}

		argMeta.Type = inferType(flag)

		if v, ok := flag.Annotations["valid-values"]; ok && len(v) > 0 {
			argMeta.ValidValues = strings.Split(v[0], ",")
		}

		metadata = append(metadata, argMeta)
	})

	return metadata, nil
}

func isRequiredFlag(flag *pflag.Flag, cmd *cobra.Command) bool {
	usage := strings.ToLower(flag.Usage)
	if strings.Contains(usage, "(required)") || strings.Contains(usage, "required") {
		return true
	}

	return false
}

func parseDefaultValue(defValue string) interface{} {
	var val interface{}

	fmt.Sscanf(defValue, "%v", &val)

	if val == "true" || val == "false" {
		if val == "true" {
			return true
		}
		return false
	}

	var intVal int
	if _, err := fmt.Sscanf(defValue, "%d", &intVal); err == nil {
		return intVal
	}

	return defValue
}

func inferType(flag *pflag.Flag) string {
	switch flag.Value.Type() {
	case "int", "int8", "int16", "int32", "int64":
		return "int"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "int"
	case "float32", "float64":
		return "float"
	case "bool":
		return "bool"
	case "string":
		return "string"
	default:
		return "string"
	}
}
