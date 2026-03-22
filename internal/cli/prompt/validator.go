package prompt

import (
	"fmt"
	"strconv"
	"strings"
)

type ValidationResult struct {
	IsValid      bool
	ErrorMessage string
	Value        interface{}
}

func ValidateInput(input string, arg ArgumentMetadata) ValidationResult {
	if input == "" {
		return handleEmptyInput(arg)
	}

	input = strings.TrimSpace(input)

	switch strings.ToLower(input) {
	case "help", "?":
		return ValidationResult{IsValid: true, Value: "help"}
	}

	switch arg.Type {
	case "int":
		return validateInt(input, arg)
	case "bool":
		return validateBool(input, arg)
	case "string":
		return validateString(input, arg)
	default:
		return ValidationResult{IsValid: true, Value: input}
	}
}

func handleEmptyInput(arg ArgumentMetadata) ValidationResult {
	if arg.IsRequired {
		return ValidationResult{
			IsValid:      false,
			ErrorMessage: "This field is required. Please enter a value.",
		}
	}

	if arg.DefaultValue != nil {
		return ValidationResult{
			IsValid: true,
			Value:   arg.DefaultValue,
		}
	}

	return ValidationResult{
		IsValid: true,
		Value:   nil,
	}
}

func validateInt(input string, arg ArgumentMetadata) ValidationResult {
	val, err := strconv.Atoi(input)
	if err != nil {
		return ValidationResult{
			IsValid:      false,
			ErrorMessage: "Invalid integer. Please enter a valid number.",
		}
	}

	if arg.ValidValues != nil && len(arg.ValidValues) > 0 {
		valid := false
		for _, v := range arg.ValidValues {
			intVal, _ := strconv.Atoi(v)
			if val == intVal {
				valid = true
				break
			}
		}
		if !valid {
			return ValidationResult{
				IsValid:      false,
				ErrorMessage: fmt.Sprintf("Invalid value. Allowed values: %s", strings.Join(arg.ValidValues, ", ")),
			}
		}
	}

	return ValidationResult{IsValid: true, Value: val}
}

func validateBool(input string, arg ArgumentMetadata) ValidationResult {
	input = strings.ToLower(input)

	boolValues := map[string]bool{
		"true":  true,
		"yes":   true,
		"y":     true,
		"1":     true,
		"false": false,
		"no":    false,
		"n":     false,
		"0":     false,
	}

	if val, ok := boolValues[input]; ok {
		return ValidationResult{IsValid: true, Value: val}
	}

	return ValidationResult{
		IsValid:      false,
		ErrorMessage: "Invalid boolean value. Please enter y/yes/true/1 or n/no/false/0.",
	}
}

func validateString(input string, arg ArgumentMetadata) ValidationResult {
	if arg.ValidValues != nil && len(arg.ValidValues) > 0 {
		valid := false
		for _, v := range arg.ValidValues {
			if strings.ToLower(input) == strings.ToLower(v) {
				valid = true
				input = v
				break
			}
		}
		if !valid {
			return ValidationResult{
				IsValid:      false,
				ErrorMessage: fmt.Sprintf("Invalid value. Allowed values: %s", strings.Join(arg.ValidValues, ", ")),
			}
		}
	}

	return ValidationResult{IsValid: true, Value: input}
}
