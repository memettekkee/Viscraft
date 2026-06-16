package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func generateTestToken(claims jwt.MapClaims, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func setupAuthRouter(secret string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("requestId", "test-request-id")
		c.Next()
	})
	r.Use(JWTAuth(secret))
	r.GET("/protected", func(c *gin.Context) {
		userId, _ := c.Get("userId")
		c.JSON(http.StatusOK, gin.H{"userId": userId})
	})
	return r
}

func TestJWTAuth_ValidToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	claims := jwt.MapClaims{
		"userId": "user-123",
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := generateTestToken(claims, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestJWTAuth_MissingAuthorizationHeader(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_InvalidPrefix(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Basic some-token")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	claims := jwt.MapClaims{
		"userId": "user-123",
		"exp":    jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}
	token := generateTestToken(claims, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_WrongSecret(t *testing.T) {
	router := setupAuthRouter(testSecret)

	claims := jwt.MapClaims{
		"userId": "user-123",
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := generateTestToken(claims, "wrong-secret")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_MalformedToken(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_MissingUserIdClaim(t *testing.T) {
	router := setupAuthRouter(testSecret)

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := generateTestToken(claims, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_EmptyUserIdClaim(t *testing.T) {
	router := setupAuthRouter(testSecret)

	claims := jwt.MapClaims{
		"userId": "",
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := generateTestToken(claims, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
}

func TestJWTAuth_SetsUserIdInContext(t *testing.T) {
	r := gin.New()
	r.Use(JWTAuth(testSecret))

	var capturedUserId string
	r.GET("/protected", func(c *gin.Context) {
		capturedUserId = c.GetString("userId")
		c.JSON(http.StatusOK, nil)
	})

	claims := jwt.MapClaims{
		"userId": "user-456",
		"exp":    jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := generateTestToken(claims, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if capturedUserId != "user-456" {
		t.Errorf("expected userId 'user-456', got '%s'", capturedUserId)
	}
}

func TestJWTAuth_ResponseContainsErrorCode(t *testing.T) {
	router := setupAuthRouter(testSecret)

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	body := w.Body.String()
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", w.Code)
	}
	// Check that ERR_09 is in the response body
	if !contains(body, "ERR_09") {
		t.Errorf("expected response to contain ERR_09, got: %s", body)
	}
	// Check requestId is included
	if !contains(body, "test-request-id") {
		t.Errorf("expected response to contain requestId, got: %s", body)
	}
}


