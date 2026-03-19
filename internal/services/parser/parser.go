package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/anomalyco/synthdata/internal/models"
)

var (
	ErrFileNotFound      = errors.New("description file not found")
	ErrParseFailed       = errors.New("failed to parse description file")
	ErrNoDescription     = errors.New("either 'description' or 'description_file' must be provided")
	ErrMarkdownFileEmpty = errors.New("markdown description file is empty")
)

type ParseError struct {
	Line    int
	Column  int
	Message string
}

func (e *ParseError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
	}
	return fmt.Sprintf("parse error: %s", e.Message)
}

func ParseDescriptionFile(path string) (*models.DescriptionFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var desc models.DescriptionFile
	if err := parseJSON(string(data), &desc); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrParseFailed, err.Error())
	}

	if desc.DescriptionFile != "" {
		markdownPath := desc.DescriptionFile
		if !filepath.IsAbs(markdownPath) {
			jsonDir := filepath.Dir(path)
			markdownPath = filepath.Join(jsonDir, markdownPath)
		}

		markdownData, err := os.ReadFile(markdownPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("%w: %s", ErrFileNotFound, markdownPath)
			}
			return nil, fmt.Errorf("failed to read description file: %w", err)
		}

		if len(markdownData) == 0 {
			return nil, ErrMarkdownFileEmpty
		}

		desc.Description = string(markdownData)
	}

	if desc.Description == "" && desc.DescriptionFile == "" {
		return nil, ErrNoDescription
	}

	return &desc, nil
}

func parseJSON(content string, desc *models.DescriptionFile) error {
	return parseJSONWithLineNumbers(content, desc)
}

func parseJSONWithLineNumbers(content string, desc *models.DescriptionFile) error {
	var lineCount int
	for _, c := range content {
		if c == '\n' {
			lineCount++
		}
	}

	return parseJSONSimple(content, desc)
}

func parseJSONSimple(content string, desc *models.DescriptionFile) error {
	return json.Unmarshal([]byte(content), desc)
}
