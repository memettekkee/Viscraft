package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequestID_GeneratesValidUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	headerVal := w.Header().Get("X-Request-ID")
	if headerVal == "" {
		t.Fatal("expected X-Request-ID header to be set")
	}

	_, err := uuid.Parse(headerVal)
	if err != nil {
		t.Fatalf("X-Request-ID is not a valid UUID: %s", headerVal)
	}
}

func TestRequestID_SetsGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())

	var contextRequestId string
	router.POST("/test", func(c *gin.Context) {
		val, exists := c.Get("requestId")
		if !exists {
			t.Fatal("expected requestId to exist in Gin context")
		}
		contextRequestId = val.(string)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	headerVal := w.Header().Get("X-Request-ID")
	if contextRequestId != headerVal {
		t.Fatalf("context requestId (%s) does not match header (%s)", contextRequestId, headerVal)
	}
}

func TestRequestID_IgnoresClientHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	clientId := "client-supplied-id-12345"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("X-Request-ID", clientId)
	router.ServeHTTP(w, req)

	headerVal := w.Header().Get("X-Request-ID")
	if headerVal == clientId {
		t.Fatal("middleware should ignore client-supplied X-Request-ID")
	}

	_, err := uuid.Parse(headerVal)
	if err != nil {
		t.Fatalf("X-Request-ID should be a server-generated UUID, got: %s", headerVal)
	}
}

func TestRequestID_UniqueBetweenRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestID())
	router.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w1, req1)

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w2, req2)

	id1 := w1.Header().Get("X-Request-ID")
	id2 := w2.Header().Get("X-Request-ID")

	if id1 == id2 {
		t.Fatal("expected unique request IDs between requests")
	}
}
