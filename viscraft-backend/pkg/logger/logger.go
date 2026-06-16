package logger

import (
	"fmt"
	"log"
	"strings"
)

// Info logs an informational message with the requestId and optional key-value fields.
func Info(requestId string, message string, fields ...interface{}) {
	log.Printf("[INFO] requestId=%s %s%s", requestId, message, formatFields(fields))
}

// Warn logs a warning message with the requestId and optional key-value fields.
func Warn(requestId string, message string, fields ...interface{}) {
	log.Printf("[WARN] requestId=%s %s%s", requestId, message, formatFields(fields))
}

// Error logs an error message with the requestId and optional key-value fields.
// If the last field is an error, it is appended as "error=<value>".
func Error(requestId string, message string, fields ...interface{}) {
	log.Printf("[ERROR] requestId=%s %s%s", requestId, message, formatFields(fields))
}

// formatFields converts variadic fields into a formatted string.
// It supports two patterns:
//   - Key-value pairs: logger.Info(id, "msg", "key1", val1, "key2", val2)
//   - Single error: logger.Error(id, "msg", err)
func formatFields(fields []interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	var sb strings.Builder

	// If there's a single field and it's an error, format it directly.
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

	// Process as key-value pairs.
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		val := fmt.Sprintf("%v", fields[i+1])
		sb.WriteString(" ")
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(val)
	}

	// If there's an odd trailing field (likely an error), append it.
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
