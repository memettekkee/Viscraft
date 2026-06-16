package middleware

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/quick"
	"time"

	"github.com/gin-gonic/gin"
)

// **Validates: Requirements 11.1, 11.2, 11.3**

// TestProperty_RateLimitCorrectness verifies that for any sequence of N > 5
// requests within a 60s window with maxRequests=5, exactly 5 pass (HTTP 200)
// and the rest are rejected (HTTP 429).
func TestProperty_RateLimitCorrectness(t *testing.T) {
	gin.SetMode(gin.TestMode)

	property := func(n uint8) bool {
		// Generate request counts between 6 and 100
		requestCount := int(n)%95 + 6 // maps 0-255 → 6-100

		// Create a fresh middleware instance per test invocation
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("userId", "property-test-user")
			c.Set("requestId", "req-prop-test")
			c.Next()
		})
		router.Use(RateLimit(5, 60*time.Second))
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		passCount := 0
		rejectCount := 0

		for i := 0; i < requestCount; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)

			switch w.Code {
			case http.StatusOK:
				passCount++
			case http.StatusTooManyRequests:
				rejectCount++
			default:
				t.Logf("unexpected status code: %d", w.Code)
				return false
			}
		}

		// Property: exactly 5 requests pass, N-5 are rejected
		if passCount != 5 {
			t.Logf("requestCount=%d, passCount=%d (expected 5)", requestCount, passCount)
			return false
		}
		if rejectCount != requestCount-5 {
			t.Logf("requestCount=%d, rejectCount=%d (expected %d)", requestCount, rejectCount, requestCount-5)
			return false
		}
		return true
	}

	cfg := &quick.Config{
		MaxCount: 100,
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 7 (Rate Limit Correctness) failed: %v", err)
	}
}

// TestProperty_RateLimitTimestampPruning verifies that timestamps older than
// 60s are pruned and don't count against the limit. After the window expires,
// the user can make up to maxRequests new requests again.
func TestProperty_RateLimitTimestampPruning(t *testing.T) {
	gin.SetMode(gin.TestMode)

	property := func(n uint8) bool {
		// Generate initial request counts between 1 and 5 (fill up to limit)
		initialRequests := int(n)%5 + 1 // maps 0-255 → 1-5

		// Use a very short window so we can test pruning without long sleeps
		window := 50 * time.Millisecond

		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("userId", "prune-test-user")
			c.Set("requestId", "req-prune-test")
			c.Next()
		})
		router.Use(RateLimit(5, window))
		router.POST("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		// Phase 1: Fill up to maxRequests (5 requests)
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Logf("Phase 1: request %d should pass but got %d", i+1, w.Code)
				return false
			}
		}

		// Verify limit is enforced (6th request should fail)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusTooManyRequests {
			t.Logf("Phase 1: 6th request should be rejected but got %d", w.Code)
			return false
		}

		// Phase 2: Wait for window to expire so old timestamps are pruned
		time.Sleep(window + 10*time.Millisecond)

		// Phase 3: After pruning, user should be able to make up to initialRequests
		// (testing with variable count to exercise the property)
		for i := 0; i < initialRequests; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)
			if w.Code != http.StatusOK {
				t.Logf("Phase 3: request %d of %d should pass after pruning but got %d",
					i+1, initialRequests, w.Code)
				return false
			}
		}

		return true
	}

	cfg := &quick.Config{
		MaxCount: 50,
		Rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	if err := quick.Check(property, cfg); err != nil {
		t.Errorf("Property 7 (Rate Limit Timestamp Pruning) failed: %v", err)
	}
}
