package response

type BaseResponse struct {
	RequestId string `json:"requestId"`
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message"`
}
