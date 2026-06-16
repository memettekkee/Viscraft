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

// UserController handles HTTP requests for user-related operations.
// It implements the router.Controller interface.
type UserController struct {
	userService *service.UserService
}

// NewUserController creates a new UserController with the given UserService dependency.
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// Routes returns the route definitions for user endpoints.
func (uc *UserController) Routes() []router.Route {
	return []router.Route{
		{Path: "/users/create", Handler: uc.Create, Protected: false},
		{Path: "/users/login", Handler: uc.Login, Protected: false},
		{Path: "/users/get", Handler: uc.Get, Protected: true},
		{Path: "/users/delete", Handler: uc.Delete, Protected: true},
	}
}

// Create handles user registration requests.
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

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusCreated, res)
}

// Login handles user authentication requests.
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

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

// Get handles requests for the authenticated user's profile.
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

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}

// Delete handles user account deletion requests.
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

	// Override requestId to ensure server-generated value is used
	res.RequestId = requestId
	c.JSON(http.StatusOK, res)
}
