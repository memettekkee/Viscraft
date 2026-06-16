package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"
)

func TestRecovery_CatchesPanicAndReturns500(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("requestId", "test-request-id-123")
		c.Next()
	})
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("something went wrong")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}

	var resp response.BaseResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.RequestId != "test-request-id-123" {
		t.Errorf("expected requestId 'test-request-id-123', got '%s'", resp.RequestId)
	}
	if resp.Success != false {
		t.Errorf("expected success=false, got %v", resp.Success)
	}
	if resp.ErrorCode != constant.ErrInternalServer.Code {
		t.Errorf("expected errorCode '%s', got '%s'", constant.ErrInternalServer.Code, resp.ErrorCode)
	}
	if resp.Message != constant.ErrInternalServer.Message {
		t.Errorf("expected message '%s', got '%s'", constant.ErrInternalServer.Message, resp.Message)
	}
}

func TestRecovery_DoesNotExposeInternalDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("requestId", "req-456")
		c.Next()
	})
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("secret internal error: database password is hunter2")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w, req)

	body := w.Body.String()

	// The response should NOT contain the panic value or any internal details
	if contains(body, "secret internal error") {
		t.Error("response body exposes panic value")
	}
	if contains(body, "hunter2") {
		t.Error("response body exposes internal details")
	}
	if contains(body, ".go") {
		t.Error("response body exposes file paths")
	}
	if contains(body, "goroutine") {
		t.Error("response body exposes stack trace")
	}
}

func TestRecovery_NoPanic_PassesThrough(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("requestId", "req-789")
		c.Next()
	})
	r.Use(Recovery())
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRecovery_ServerContinuesAfterPanic(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("requestId", "req-cont")
		c.Next()
	})
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// First request panics
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/panic", nil)
	r.ServeHTTP(w1, req1)

	if w1.Code != http.StatusInternalServerError {
		t.Errorf("first request: expected 500, got %d", w1.Code)
	}

	// Second request should work normally
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/ok", nil)
	r.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("second request: expected 200, got %d", w2.Code)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
