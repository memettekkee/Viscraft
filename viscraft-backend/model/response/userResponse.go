package response

// UserData contains the user fields returned in responses.
type UserData struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

// CreateUserResponse is returned after successful user registration.
type CreateUserResponse struct {
	BaseResponse
	Token string    `json:"token"`
	Data  *UserData `json:"data,omitempty"`
}

// LoginResponse is returned after successful authentication.
type LoginResponse struct {
	BaseResponse
	Token string    `json:"token"`
	Data  *UserData `json:"data,omitempty"`
}

// DeleteUserResponse is returned after successful user deletion.
type DeleteUserResponse struct {
	BaseResponse
}
