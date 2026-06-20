package response

type UserData struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
}

type CreateUserResponse struct {
	BaseResponse
	Token string    `json:"token"`
	Data  *UserData `json:"data,omitempty"`
}

type LoginResponse struct {
	BaseResponse
	Token string    `json:"token"`
	Data  *UserData `json:"data,omitempty"`
}

type GetUserResponse struct {
	BaseResponse
	Data *UserData `json:"data,omitempty"`
}

type DeleteUserResponse struct {
	BaseResponse
}
