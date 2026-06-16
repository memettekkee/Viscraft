package request

// BaseRequest is embedded by all request structs.
// The client-provided RequestId is overridden by the backend with a server-generated UUID.
type BaseRequest struct {
	RequestId string `json:"requestId"`
}
