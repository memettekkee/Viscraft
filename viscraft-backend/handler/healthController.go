package handler

import (
	"net/http"

	"viscraft-backend/model/response"
	"viscraft-backend/pkg/router"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (hc *HealthController) Routes() []router.Route {
	return []router.Route{
		{Path: "/health/check", Handler: hc.Check, Protected: false},
	}
}

func (hc *HealthController) Check(c *gin.Context) {
	requestId := c.GetString("requestId")

	c.JSON(http.StatusOK, response.BaseResponse{
		RequestId: requestId,
		Success:   true,
		Message:   "Service is operational",
	})
}
