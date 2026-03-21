package batch

import (
	"encoding/json"
	"strings"

	"github.com/dynamder/synthdata/internal/services/llm"
)

type OutputParser struct{}

func NewOutputParser() *OutputParser {
	return &OutputParser{}
}

func (p *OutputParser) Parse(response string) (valid []map[string]interface{}, invalid []FailedRecord, err error) {
	cleaned := cleanJSONResponse(response)

	var records []map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &records); err != nil {
		return nil, []FailedRecord{{
				OriginalOutput: response,
				Error:          llm.NewLLmCallError(llm.InvalidFormat, err),
			}},
			nil
	}

	for _, record := range records {
		if p.isValidRecord(record) {
			valid = append(valid, record)
		} else {
			recordJSON, _ := json.Marshal(record)
			invalid = append(invalid, FailedRecord{
				OriginalOutput: string(recordJSON),
				Error:          llm.NewLLmCallError(llm.GenerationFailed, err),
			})
		}
	}

	return valid, invalid, nil
}

func (p *OutputParser) isValidRecord(record map[string]interface{}) bool {
	//TODO: add schema validation
	if record == nil {
		return false
	}
	return len(record) > 0
}

func (p *OutputParser) IsEmpty(response string) bool {
	cleaned := strings.TrimSpace(response)
	return len(cleaned) == 0
}
