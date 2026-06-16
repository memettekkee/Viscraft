package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"viscraft-backend/constant"

	"github.com/gin-gonic/gin"
)

func TestTimeout_HandlerCompletesBeforeDeadline(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(Timeout(1 * time.Second))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestTimeout_DeadlineExceeded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(Timeout(50 * time.Millisecond))
	r.POST("/slow", func(c *gin.Context) {
		// Simulate a slow handler that exceeds the timeout.
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(200 * time.Millisecond):
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/slow", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("expected status 504, got %d", w.Code)
	}

	body := w.Body.String()
	if !contains(body, constant.ErrGeminiTimeout.Code) {
		t.Errorf("expected error code %s in body, got: %s", constant.ErrGeminiTimeout.Code, body)
	}
	if !contains(body, constant.ErrGeminiTimeout.Message) {
		t.Errorf("expected message %q in body, got: %s", constant.ErrGeminiTimeout.Message, body)
	}
}

func TestTimeout_ResponseIncludesRequestId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.Use(Timeout(50 * time.Millisecond))
	r.POST("/slow", func(c *gin.Context) {
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(200 * time.Millisecond):
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	req := httptest.NewRequest(http.MethodPost, "/slow", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := w.Body.String()
	if !contains(body, "requestId") {
		t.Errorf("expected requestId in response body, got: %s", body)
	}
}


