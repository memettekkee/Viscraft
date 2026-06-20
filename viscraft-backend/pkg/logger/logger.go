package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	once      sync.Once
	logWriter io.Writer
)

// init sets up the logger to write to both stdout and a daily log file.
func initLogger() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.SetOutput(os.Stdout)
		return
	}

	today := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logDir, "app-"+today+".log")

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.SetOutput(os.Stdout)
		return
	}

	logWriter = io.MultiWriter(os.Stdout, f)
	log.SetOutput(logWriter)
	log.SetFlags(0) 
}

func setup() {
	once.Do(initLogger)
}

func timestamp() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func Info(requestId string, message string, fields ...interface{}) {
	setup()
	log.Printf("%s [INFO] requestId=%s %s%s", timestamp(), requestId, message, formatFields(fields))
}

func Warn(requestId string, message string, fields ...interface{}) {
	setup()
	log.Printf("%s [WARN] requestId=%s %s%s", timestamp(), requestId, message, formatFields(fields))
}

func Error(requestId string, message string, fields ...interface{}) {
	setup()
	log.Printf("%s [ERROR] requestId=%s %s%s", timestamp(), requestId, message, formatFields(fields))
}

func formatFields(fields []interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var sb strings.Builder

	if len(fields) == 1 {
		if err, ok := fields[0].(error); ok {
			sb.WriteString(" error=")
			sb.WriteString(err.Error())
			return sb.String()
		}
		sb.WriteString(" ")
		sb.WriteString(fmt.Sprintf("%v", fields[0]))
		return sb.String()
	}

	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		val := fmt.Sprintf("%v", fields[i+1])
		sb.WriteString(" ")
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(val)
	}

	if len(fields)%2 != 0 {
		last := fields[len(fields)-1]
		if err, ok := last.(error); ok {
			sb.WriteString(" error=")
			sb.WriteString(err.Error())
		} else {
			sb.WriteString(" ")
			sb.WriteString(fmt.Sprintf("%v", last))
		}
	}

	return sb.String()
}
