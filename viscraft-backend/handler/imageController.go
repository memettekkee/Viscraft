package handler

import (
	"net/http"

	"viscraft-backend/constant"
	"viscraft-backend/model/request"
	"viscraft-backend/model/response"
	"viscraft-backend/pkg/router"
	"viscraft-backend/service"

	"github.com/gin-gonic/gin"
)

// ImageController handles HTTP requests for image-related operations.
// It implements the router.Controller interface.
type ImageController struct {
	imageService *service.ImageService
}

// NewImageController creates a new ImageController with the given ImageService dependency.
func NewImageController(imageService *service.ImageService) *ImageController {
	return &ImageController{imageService: imageService}
}

// Routes returns the route definitions for image endpoints.
// All image routes are protected (require JWT authentication).
// Note: /images/generate is excluded here and registered separately in main.go
// with rate limit middleware applied.
func (ic *ImageController) Routes() []router.Route {
	return []router.Route{
		{Path: "/images/get", Handler: ic.Get, Protected: true},
		{Path: "/images/list", Handler: ic.List, Protected: true},
		{Path: "/images/delete", Handler: ic.Delete, Protected: true},
	}
}

// Generate handles image generation requests.
// Returns HTTP 200 for cache hit, HTTP 202 for new generation, HTTP 422 for validation errors.
func (ic *ImageController) Generate(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrInvalidPrompt.Code,
			Message:   "Missing required fields",
		})
		return
	}

	res, appErr, cacheHit := ic.imageService.Generate(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId

	if cacheHit {
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusAccepted, res)
	}
}

// Get handles requests for retrieving a single image by ID.
func (ic *ImageController) Get(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.GetImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := ic.imageService.GetImage(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

// List handles requests for listing all images for a project.
func (ic *ImageController) List(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.ListImagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := ic.imageService.ListImages(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

// Delete handles image deletion requests.
func (ic *ImageController) Delete(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.DeleteImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := ic.imageService.DeleteImage(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}
