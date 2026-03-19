package synthdata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	logFile *os.File
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

var logger *Logger

func InitLogger() error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_150405")
	logFile, err := os.Create(filepath.Join(logDir, fmt.Sprintf("synthdata_%s.jsonl", timestamp)))
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}

	logger = &Logger{logFile: logFile}
	log.SetOutput(logFile)
	log.SetFlags(0)

	return nil
}

func GetLogger() *Logger {
	if logger == nil {
		if err := InitLogger(); err != nil {
			log.Printf("Failed to initialize logger: %v", err)
			logger = &Logger{logFile: os.Stderr}
		}
	}
	return logger
}

func (l *Logger) Info(msg string, context map[string]interface{}) {
	l.log("INFO", msg, context)
}

func (l *Logger) Error(msg string, context map[string]interface{}) {
	l.log("ERROR", msg, context)
}

func (l *Logger) Warn(msg string, context map[string]interface{}) {
	l.log("WARN", msg, context)
}

func (l *Logger) Debug(msg string, context map[string]interface{}) {
	l.log("DEBUG", msg, context)
}

func (l *Logger) log(level, msg string, context map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Context:   context,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		log.Printf("Failed to marshal log entry: %v", err)
		return
	}

	fmt.Fprintln(logger.logFile, string(data))
	if logger.logFile != os.Stderr && logger.logFile != os.Stdout {
		logger.logFile.Sync()
	}

	log.Printf("[%s] %s", level, msg)
}

func (l *Logger) Close() error {
	if logger != nil && logger.logFile != nil && logger.logFile != os.Stderr && logger.logFile != os.Stdout {
		return logger.logFile.Close()
	}
	return nil
}
