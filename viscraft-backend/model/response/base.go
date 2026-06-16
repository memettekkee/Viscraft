package response

// BaseResponse is embedded by all response structs.
// Every response includes a server-generated RequestId, a Success flag,
// an optional ErrorCode for failures, and a human-readable Message.
type BaseResponse struct {
	RequestId string `json:"requestId"`
	Success   bool   `json:"success"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message"`
}
