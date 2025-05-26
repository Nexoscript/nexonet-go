package logger

import (
	"fmt"
	"sync"
	"time"
)

type LoggingType string

const (
	INFO  LoggingType = "INFO"
	WARN  LoggingType = "WARN"
	ERROR LoggingType = "ERROR"
	DEBUG LoggingType = "DEBUG"
)

type NexonetLogger struct {
	enableLogging bool
	mu            sync.Mutex
}

func NewNexonetLogger(enable bool) *NexonetLogger {
	return &NexonetLogger{enableLogging: enable}
}

func (l *NexonetLogger) Log(logType LoggingType, message string) {
	if l.enableLogging {
		l.mu.Lock()
		defer l.mu.Unlock()
		fmt.Printf("[%s] [%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), logType, message)
	}
}
