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

type SceneController struct {
	sceneService *service.SceneService
}

func NewSceneController(sceneService *service.SceneService) *SceneController {
	return &SceneController{sceneService: sceneService}
}

func (sc *SceneController) Routes() []router.Route {
	return []router.Route{
		{Path: "/scenes/get", Handler: sc.Get, Protected: true},
		{Path: "/scenes/list", Handler: sc.List, Protected: true},
		{Path: "/scenes/delete", Handler: sc.Delete, Protected: true},
	}
}

func (sc *SceneController) Generate(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.GenerateSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrInvalidPrompt.Code,
			Message:   "Missing required fields",
		})
		return
	}

	res, appErr := sc.sceneService.Generate(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	res.RequestId = requestId
	c.JSON(http.StatusAccepted, res)
}

func (sc *SceneController) Get(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.GetSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := sc.sceneService.GetScene(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

func (sc *SceneController) List(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.ListScenesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := sc.sceneService.ListScenes(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

func (sc *SceneController) Delete(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	var req request.DeleteSceneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := sc.sceneService.DeleteScene(requestId, userId, req)
	if appErr != nil {
		c.JSON(appErr.HttpStatus, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: appErr.Code,
			Message:   appErr.Message,
		})
		return
	}

	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}
