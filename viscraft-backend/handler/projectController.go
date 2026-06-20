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

type ProjectController struct {
	projectService *service.ProjectService
}

func NewProjectController(projectService *service.ProjectService) *ProjectController {
	return &ProjectController{projectService: projectService}
}

func (pc *ProjectController) Routes() []router.Route {
	return []router.Route{
		{Path: "/projects/create", Handler: pc.Create, Protected: true},
		{Path: "/projects/get", Handler: pc.Get, Protected: true},
		{Path: "/projects/list", Handler: pc.List, Protected: true},
		{Path: "/projects/delete", Handler: pc.Delete, Protected: true},
	}
}

func (pc *ProjectController) Create(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := pc.projectService.CreateProject(requestId, userId, req)
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
	c.JSON(http.StatusCreated, res)
}

// Get handles requests for retrieving a single project.
func (pc *ProjectController) Get(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.GetProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := pc.projectService.GetProject(requestId, userId, req)
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

// List handles requests for listing all projects for the authenticated user.
func (pc *ProjectController) List(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	// Bind JSON body (optional fields in ListProjectsRequest)
	var req request.ListProjectsRequest
	_ = c.ShouldBindJSON(&req)

	res, appErr := pc.projectService.ListProjects(requestId, userId)
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

// Delete handles project deletion requests.
func (pc *ProjectController) Delete(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.DeleteProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := pc.projectService.DeleteProject(requestId, userId, req)
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
