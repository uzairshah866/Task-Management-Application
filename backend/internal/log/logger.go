package log

import (
	"log"
	"os"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel = INFO
	logger       = log.New(os.Stderr, "", log.LstdFlags)
)

func init() {
	if os.Getenv("LOG_LEVEL") == "debug" {
		currentLevel = DEBUG
	}
}

func Debug(msg string, args ...interface{}) {
	if currentLevel <= DEBUG {
		logger.Printf("[DEBUG] "+msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	if currentLevel <= INFO {
		logger.Printf("[INFO] "+msg, args...)
	}
}

func Warn(msg string, args ...interface{}) {
	if currentLevel <= WARN {
		logger.Printf("[WARN] "+msg, args...)
	}
}

func Error(msg string, args ...interface{}) {
	if currentLevel <= ERROR {
		logger.Printf("[ERROR] "+msg, args...)
	}
}

func ErrorWithStack(err error, context string) {
	logger.Printf("[ERROR] %s: %v (time=%s)", context, err, time.Now().Format(time.RFC3339))
}
