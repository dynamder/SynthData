package prompt

import (
	"fmt"
	"strings"
)

func BuildPrompt(arg ArgumentMetadata, index, total int) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n▶ %s", formatArgName(arg.Name)))

	if arg.Short != "" {
		sb.WriteString(fmt.Sprintf(" (-%s, --%s)", arg.Short, arg.Name))
	} else {
		sb.WriteString(fmt.Sprintf(" (--%s)", arg.Name))
	}

	sb.WriteString(":\n")

	if arg.Description != "" {
		sb.WriteString(fmt.Sprintf("  %s\n", arg.Description))
	}

	if arg.DefaultValue != nil && !arg.IsRequired {
		sb.WriteString(fmt.Sprintf("  [default: %v]\n", arg.DefaultValue))
	}

	if arg.ValidValues != nil && len(arg.ValidValues) > 0 {
		sb.WriteString(fmt.Sprintf("  [%s]\n", strings.Join(arg.ValidValues, ", ")))
	}

	if arg.IsRequired {
		sb.WriteString("  > (required) ")
	} else {
		sb.WriteString("  > ")
	}

	return sb.String()
}

func formatArgName(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return strings.Title(name)
}

func ShowHelp(arg ArgumentMetadata) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n=== Help for %s ===\n", formatArgName(arg.Name)))

	if arg.Description != "" {
		sb.WriteString(fmt.Sprintf("\nDescription: %s\n", arg.Description))
	}

	sb.WriteString(fmt.Sprintf("\nType: %s\n", arg.Type))

	if arg.IsRequired {
		sb.WriteString("Required: yes\n")
	} else {
		sb.WriteString("Required: no\n")
	}

	if arg.DefaultValue != nil {
		sb.WriteString(fmt.Sprintf("Default: %v\n", arg.DefaultValue))
	}

	if arg.ValidValues != nil && len(arg.ValidValues) > 0 {
		sb.WriteString(fmt.Sprintf("Valid values: %s\n", strings.Join(arg.ValidValues, ", ")))
	}

	if arg.Short != "" {
		sb.WriteString(fmt.Sprintf("\nFlag: -%s, --%s\n", arg.Short, arg.Name))
	} else {
		sb.WriteString(fmt.Sprintf("\nFlag: --%s\n", arg.Name))
	}

	return sb.String()
}

func ShowWelcome() string {
	return `
Welcome to Interactive Mode!
Press Ctrl+C at any time to exit.
`
}

func ShowProgress(index, total int) string {
	return fmt.Sprintf("\n[%d/%d] ", index+1, total)
}

func ShowCollectedArgs(args map[string]interface{}) string {
	var sb strings.Builder

	sb.WriteString("\nCollected arguments:\n")

	for k, v := range args {
		if v != nil {
			sb.WriteString(fmt.Sprintf("  --%s %v\n", k, v))
		}
	}

	return sb.String()
}

func ShowError(message string) string {
	return fmt.Sprintf("\n✗ %s\n", message)
}

func ShowSuccess(message string) string {
	return fmt.Sprintf("\n✓ %s\n", message)
}
