package request

// CreateUserRequest represents the payload for creating a new user account.
type CreateUserRequest struct {
	BaseRequest
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name"`
}

// LoginRequest represents the payload for authenticating an existing user.
type LoginRequest struct {
	BaseRequest
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// DeleteUserRequest represents the payload for deleting the authenticated user.
// User identity is extracted from the JWT context.
type DeleteUserRequest struct {
	BaseRequest
}
