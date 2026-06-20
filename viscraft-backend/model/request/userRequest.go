package request

type CreateUserRequest struct {
	BaseRequest
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	BaseRequest
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type DeleteUserRequest struct {
	BaseRequest
}
