package logger

import (
	"bytes"
	"errors"
	"log"
	"os"
	"strings"
	"testing"
)

// captureOutput captures log output during a function call.
func captureOutput(fn func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(os.Stderr)
	fn()
	return buf.String()
}

func TestInfoBasicMessage(t *testing.T) {
	output := captureOutput(func() {
		Info("req-123", "cache hit, returning existing image")
	})

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", output)
	}
	if !strings.Contains(output, "requestId=req-123") {
		t.Errorf("expected requestId=req-123 in output, got: %s", output)
	}
	if !strings.Contains(output, "cache hit, returning existing image") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestInfoWithKeyValueFields(t *testing.T) {
	output := captureOutput(func() {
		Info("req-456", "calling Gemini API", "prompt", "draw a dragon")
	})

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("expected [INFO] in output, got: %s", output)
	}
	if !strings.Contains(output, "requestId=req-456") {
		t.Errorf("expected requestId=req-456 in output, got: %s", output)
	}
	if !strings.Contains(output, "prompt=draw a dragon") {
		t.Errorf("expected prompt=draw a dragon in output, got: %s", output)
	}
}

func TestErrorWithError(t *testing.T) {
	err := errors.New("connection refused")
	output := captureOutput(func() {
		Error("req-789", "prompt validation failed", err)
	})

	if !strings.Contains(output, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", output)
	}
	if !strings.Contains(output, "requestId=req-789") {
		t.Errorf("expected requestId=req-789 in output, got: %s", output)
	}
	if !strings.Contains(output, "error=connection refused") {
		t.Errorf("expected error=connection refused in output, got: %s", output)
	}
}

func TestWarnBasicMessage(t *testing.T) {
	output := captureOutput(func() {
		Warn("req-warn-1", "rate limit approaching")
	})

	if !strings.Contains(output, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", output)
	}
	if !strings.Contains(output, "requestId=req-warn-1") {
		t.Errorf("expected requestId=req-warn-1 in output, got: %s", output)
	}
	if !strings.Contains(output, "rate limit approaching") {
		t.Errorf("expected message in output, got: %s", output)
	}
}

func TestInfoWithMultipleKeyValuePairs(t *testing.T) {
	output := captureOutput(func() {
		Info("req-multi", "image generation completed", "imageId", "img-001", "duration", "2.5s")
	})

	if !strings.Contains(output, "imageId=img-001") {
		t.Errorf("expected imageId=img-001 in output, got: %s", output)
	}
	if !strings.Contains(output, "duration=2.5s") {
		t.Errorf("expected duration=2.5s in output, got: %s", output)
	}
}

func TestInfoNoFields(t *testing.T) {
	output := captureOutput(func() {
		Info("req-none", "simple message")
	})

	if !strings.Contains(output, "[INFO] requestId=req-none simple message") {
		t.Errorf("unexpected output format: %s", output)
	}
}
