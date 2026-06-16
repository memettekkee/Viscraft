package handler

import (
	"net/http"

	"viscraft-backend/model/response"
	"viscraft-backend/pkg/router"

	"github.com/gin-gonic/gin"
)

// HealthController handles HTTP requests for health check operations.
// It implements the router.Controller interface.
type HealthController struct{}

// NewHealthController creates a new HealthController.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Routes returns the route definitions for health check endpoints.
func (hc *HealthController) Routes() []router.Route {
	return []router.Route{
		{Path: "/health/check", Handler: hc.Check, Protected: false},
	}
}

// Check handles health check requests.
// Returns HTTP 200 with success=true and a message indicating the service is operational.
func (hc *HealthController) Check(c *gin.Context) {
	requestId := c.GetString("requestId")

	c.JSON(http.StatusOK, response.BaseResponse{
		RequestId: requestId,
		Success:   true,
		Message:   "Service is operational",
	})
}
