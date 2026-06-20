package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Routes() []Route
}

// Route defines a single endpoint registration.
type Route struct {
	Path      string          // URL path in /{resource}/{action} format
	Handler   gin.HandlerFunc // Handler function for the route
	Protected bool            // If true, JWT auth middleware is applied
}

// Register registers all routes from the given controllers on the Gin engine.
// All routes use the POST HTTP method exclusively.
// If a route's Protected flag is true, the authMiddleware is applied to that route.
// Returns an error if duplicate route paths are detected.
func Register(r *gin.Engine, authMiddleware gin.HandlerFunc, controllers ...Controller) error {
	registered := make(map[string]bool)

	for _, ctrl := range controllers {
		for _, route := range ctrl.Routes() {
			if registered[route.Path] {
				return fmt.Errorf("duplicate route path detected: %s", route.Path)
			}
			registered[route.Path] = true

			if route.Protected {
				r.POST(route.Path, authMiddleware, route.Handler)
			} else {
				r.POST(route.Path, route.Handler)
			}
		}
	}

	return nil
}
