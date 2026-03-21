package batch

import (
	"context"
	"testing"
)

func TestNewService(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	if service == nil {
		t.Error("expected non-nil service")
	}
	if service.client == nil {
		t.Error("expected client to be set")
	}
	if service.batcher == nil {
		t.Error("expected batcher to be set")
	}
	if service.executor == nil {
		t.Error("expected executor to be set")
	}
}

func TestService_Generate_InvalidTargetCount(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	req := GenerationRequest{
		TargetCount: 0,
		BatchSize:   10,
		Format:      "json",
		Output:      "/tmp/test.json",
	}

	_, err := service.Generate(context.Background(), req)
	if err == nil {
		t.Error("expected error for zero target count")
	}
}

func TestService_Generate_NegativeTargetCount(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	req := GenerationRequest{
		TargetCount: -1,
		BatchSize:   10,
		Format:      "json",
		Output:      "/tmp/test.json",
	}

	_, err := service.Generate(context.Background(), req)
	if err == nil {
		t.Error("expected error for negative target count")
	}
}

func TestService_Generate_ZeroBatchSize(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	req := GenerationRequest{
		TargetCount:     100,
		BatchSize:       0,
		DescriptionFile: "testdata/description.yaml",
		Format:          "json",
		Output:          t.TempDir() + "/test.json",
	}

	_, err := service.Generate(context.Background(), req)
	if err != nil && err.Error() == "validation error: description file not found: " {
		t.Skip("test data not available")
	}
}

func TestService_Generate_ZeroConcurrency(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 0, 3, true)

	req := GenerationRequest{
		TargetCount:     100,
		BatchSize:       10,
		Concurrency:     0,
		DescriptionFile: "testdata/description.yaml",
		Format:          "json",
		Output:          t.TempDir() + "/test.json",
	}

	_, err := service.Generate(context.Background(), req)
	if err != nil && err.Error() == "validation error: description file not found: " {
		t.Skip("test data not available")
	}
}

func TestService_Generate_InvalidDescriptionFile(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	req := GenerationRequest{
		TargetCount:     10,
		DescriptionFile: "nonexistent.yaml",
		Format:          "json",
		Output:          "/tmp/test.json",
	}

	_, err := service.Generate(context.Background(), req)
	if err == nil {
		t.Error("expected error for nonexistent description file")
	}
}

func TestService_GetFormatter(t *testing.T) {
	tests := []struct {
		format  string
		wantErr bool
	}{
		{"json", false},
		{"JSON", false},
		{"csv", false},
		{"CSV", false},
		{"unknown", true},
		{"", true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			_, err := GetFormatter(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFormatter(%q) error = %v, wantErr %v", tt.format, err, tt.wantErr)
			}
		})
	}
}

func TestService_writeOutput_Empty(t *testing.T) {
	client := &mockClient{generateFunc: func(p string) (string, error) { return "[]", nil }}
	service := NewService(client, 5, 3, true)

	session := &GenerationSession{
		BatchResults: []BatchResult{},
	}

	tmpFile := t.TempDir() + "/test.json"

	err := service.writeOutput(session, "json", tmpFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
