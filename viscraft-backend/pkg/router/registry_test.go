package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockController implements Controller for testing purposes.
type mockController struct {
	routes []Route
}

func (m *mockController) Routes() []Route {
	return m.routes
}

// dummyHandler is a simple handler that returns 200 OK.
func dummyHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// authMiddleware is a test middleware that sets a header to indicate it ran.
func testAuthMiddleware(c *gin.Context) {
	c.Set("authApplied", true)
	c.Next()
}

func TestRegister_PublicRoute(t *testing.T) {
	r := gin.New()

	ctrl := &mockController{
		routes: []Route{
			{Path: "/health/check", Handler: dummyHandler, Protected: false},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Test that the route is accessible via POST
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/health/check", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRegister_ProtectedRoute(t *testing.T) {
	r := gin.New()

	handler := func(c *gin.Context) {
		authApplied, exists := c.Get("authApplied")
		if !exists || authApplied != true {
			t.Error("expected auth middleware to have been applied")
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

	ctrl := &mockController{
		routes: []Route{
			{Path: "/images/generate", Handler: handler, Protected: true},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/images/generate", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRegister_PostOnly(t *testing.T) {
	r := gin.New()

	ctrl := &mockController{
		routes: []Route{
			{Path: "/users/create", Handler: dummyHandler, Protected: false},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// GET should not match the route
	methods := []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodPatch}
	for _, method := range methods {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(method, "/users/create", nil)
		r.ServeHTTP(w, req)

		if w.Code == http.StatusOK {
			t.Errorf("expected %s method to NOT return 200, got %d", method, w.Code)
		}
	}

	// POST should work
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users/create", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected POST to return 200, got %d", w.Code)
	}
}

func TestRegister_DuplicateRouteReturnsError(t *testing.T) {
	r := gin.New()

	ctrl1 := &mockController{
		routes: []Route{
			{Path: "/users/create", Handler: dummyHandler, Protected: false},
		},
	}
	ctrl2 := &mockController{
		routes: []Route{
			{Path: "/users/create", Handler: dummyHandler, Protected: true},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl1, ctrl2)
	if err == nil {
		t.Fatal("expected error for duplicate route, got nil")
	}

	expected := "duplicate route path detected: /users/create"
	if err.Error() != expected {
		t.Errorf("expected error message %q, got %q", expected, err.Error())
	}
}

func TestRegister_DuplicateWithinSameController(t *testing.T) {
	r := gin.New()

	ctrl := &mockController{
		routes: []Route{
			{Path: "/projects/create", Handler: dummyHandler, Protected: true},
			{Path: "/projects/create", Handler: dummyHandler, Protected: true},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl)
	if err == nil {
		t.Fatal("expected error for duplicate route within same controller, got nil")
	}
}

func TestRegister_MultipleControllers(t *testing.T) {
	r := gin.New()

	userCtrl := &mockController{
		routes: []Route{
			{Path: "/users/create", Handler: dummyHandler, Protected: false},
			{Path: "/users/login", Handler: dummyHandler, Protected: false},
		},
	}
	projectCtrl := &mockController{
		routes: []Route{
			{Path: "/projects/create", Handler: dummyHandler, Protected: true},
			{Path: "/projects/list", Handler: dummyHandler, Protected: true},
		},
	}

	err := Register(r, testAuthMiddleware, userCtrl, projectCtrl)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify all routes are accessible
	paths := []string{"/users/create", "/users/login", "/projects/create", "/projects/list"}
	for _, path := range paths {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, path, nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200 for %s, got %d", path, w.Code)
		}
	}
}

func TestRegister_NoControllers(t *testing.T) {
	r := gin.New()

	err := Register(r, testAuthMiddleware)
	if err != nil {
		t.Fatalf("expected no error with zero controllers, got: %v", err)
	}
}

func TestRegister_AuthMiddlewareNotAppliedToPublicRoute(t *testing.T) {
	r := gin.New()

	handler := func(c *gin.Context) {
		_, exists := c.Get("authApplied")
		if exists {
			t.Error("auth middleware should NOT be applied to public routes")
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}

	ctrl := &mockController{
		routes: []Route{
			{Path: "/health/check", Handler: handler, Protected: false},
		},
	}

	err := Register(r, testAuthMiddleware, ctrl)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/health/check", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
