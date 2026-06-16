package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"viscraft-backend/constant"
	"viscraft-backend/model/response"

	"github.com/gin-gonic/gin"
)

func setupRateLimitRouter(maxRequests int, window time.Duration) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userId", "user-1")
		c.Set("requestId", "req-123")
		c.Next()
	})
	r.Use(RateLimit(maxRequests, window))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func TestRateLimit_AllowsRequestsUnderLimit(t *testing.T) {
	r := setupRateLimitRouter(5, 60*time.Second)

	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}
}

func TestRateLimit_BlocksRequestsOverLimit(t *testing.T) {
	r := setupRateLimitRouter(5, 60*time.Second)

	// Make 5 allowed requests
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
	}

	// 6th request should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w.Code)
	}

	var resp response.BaseResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.ErrorCode != constant.ErrTooManyRequest.Code {
		t.Errorf("expected error code %s, got %s", constant.ErrTooManyRequest.Code, resp.ErrorCode)
	}
	if resp.Success != false {
		t.Error("expected success=false")
	}
	if resp.RequestId != "req-123" {
		t.Errorf("expected requestId 'req-123', got '%s'", resp.RequestId)
	}
	if resp.Message != constant.ErrTooManyRequest.Message {
		t.Errorf("expected message '%s', got '%s'", constant.ErrTooManyRequest.Message, resp.Message)
	}
}

func TestRateLimit_PrunesExpiredTimestamps(t *testing.T) {
	// Use a very short window so timestamps expire quickly
	r := setupRateLimitRouter(2, 50*time.Millisecond)

	// Make 2 requests (fills the limit)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
	}

	// Wait for window to expire
	time.Sleep(60 * time.Millisecond)

	// Next request should succeed because old timestamps were pruned
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 after window expired, got %d", w.Code)
	}
}

func TestRateLimit_IsolatesUsersByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	var currentUser string
	r.Use(func(c *gin.Context) {
		c.Set("userId", currentUser)
		c.Set("requestId", "req-456")
		c.Next()
	})
	r.Use(RateLimit(2, 60*time.Second))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// User A makes 2 requests (hits limit)
	currentUser = "user-a"
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		r.ServeHTTP(w, req)
	}

	// User A's 3rd request should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("user-a: expected 429, got %d", w.Code)
	}

	// User B should still be allowed
	currentUser = "user-b"
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("user-b: expected 200, got %d", w.Code)
	}
}
