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

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) Routes() []router.Route {
	return []router.Route{
		{Path: "/users/create", Handler: uc.Create, Protected: false},
		{Path: "/users/login", Handler: uc.Login, Protected: false},
		{Path: "/users/get", Handler: uc.Get, Protected: true},
		{Path: "/users/delete", Handler: uc.Delete, Protected: true},
		{Path: "/users/complete-tour", Handler: uc.CompleteTour, Protected: true},
	}
}

func (uc *UserController) Create(c *gin.Context) {
	requestId := c.GetString("requestId")

	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := uc.userService.CreateUser(requestId, req)
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
	c.JSON(http.StatusCreated, res)
}

func (uc *UserController) Login(c *gin.Context) {
	requestId := c.GetString("requestId")

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.BaseResponse{
			RequestId: requestId,
			Success:   false,
			ErrorCode: constant.ErrValidationFailed.Code,
			Message:   constant.ErrValidationFailed.Message,
		})
		return
	}

	res, appErr := uc.userService.Login(requestId, req)
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

func (uc *UserController) Get(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	res, appErr := uc.userService.GetUser(requestId, userId)
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

func (uc *UserController) Delete(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	res, appErr := uc.userService.DeleteUser(requestId, userId)
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

func (uc *UserController) CompleteTour(c *gin.Context) {
	requestId := c.GetString("requestId")
	userId := c.GetString("userId")

	res, appErr := uc.userService.CompleteTour(requestId, userId)
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
