package synthdata

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLogger_InitLogger(t *testing.T) {
	tmpDir := t.TempDir()
	logDir := filepath.Join(tmpDir, "logs")
	os.MkdirAll(logDir, 0755)

	origLogger := logger
	defer func() { logger = origLogger }()

	logFile, err := os.CreateTemp(tmpDir, "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp log file: %v", err)
	}
	defer logFile.Close()

	logger = &Logger{logFile: logFile}

	l := GetLogger()
	if l == nil {
		t.Error("expected non-nil logger")
	}
}

func TestLogger_Info(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: tmpFile}

	l := GetLogger()
	l.Info("test message", map[string]interface{}{"key": "value"})
	logger.logFile.Sync()
}

func TestLogger_Error(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: tmpFile}

	l := GetLogger()
	l.Error("error message", map[string]interface{}{"code": 500})
	logger.logFile.Sync()
}

func TestLogger_Warn(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: tmpFile}

	l := GetLogger()
	l.Warn("warning message", nil)
	logger.logFile.Sync()
}

func TestLogger_Debug(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: tmpFile}

	l := GetLogger()
	l.Debug("debug message", map[string]interface{}{"debug": true})
	logger.logFile.Sync()
}

func TestLogger_Close(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_log_*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: tmpFile}

	err = logger.Close()
	if err != nil {
		t.Errorf("unexpected error closing logger: %v", err)
	}
}

func TestLogger_Close_NilFile(t *testing.T) {
	origLogger := logger
	defer func() { logger = origLogger }()

	logger = &Logger{logFile: nil}

	err := logger.Close()
	if err != nil {
		t.Errorf("unexpected error closing nil file: %v", err)
	}
}

func TestGetLogger_Default(t *testing.T) {
	origLogger := logger
	defer func() { logger = origLogger }()
	logger = nil

	l := GetLogger()
	if l == nil {
		t.Error("expected non-nil logger")
	}
}
