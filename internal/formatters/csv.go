package formatters

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/dynamder/synthdata/internal/models"
)

type CSVFormatter struct{}

func NewCSVFormatter() *CSVFormatter {
	return &CSVFormatter{}
}

func (f *CSVFormatter) Format(records []map[string]interface{}) ([]byte, error) {
	if len(records) == 0 {
		return []byte{}, nil
	}

	headers := getHeaders(records)
	var builder strings.Builder
	writer := csv.NewWriter(&builder)

	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	for _, record := range records {
		row := make([]string, len(headers))
		for i, header := range headers {
			val := record[header]
			if val != nil {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV flush error: %w", err)
	}

	return []byte(builder.String()), nil
}

func getHeaders(records []map[string]interface{}) []string {
	headerSet := make(map[string]bool)
	for _, record := range records {
		for key := range record {
			headerSet[key] = true
		}
	}

	headers := make([]string, 0, len(headerSet))
	for header := range headerSet {
		headers = append(headers, header)
	}
	return headers
}

func ParseCSVFieldType(fieldType string) models.FieldType {
	return ParseJSONFieldType(fieldType)
}
