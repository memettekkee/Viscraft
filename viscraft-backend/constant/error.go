package constant

// AppError represents a predefined application error with a machine-readable code,
// a human-readable message, and the associated HTTP status code.
type AppError struct {
	Code       string
	Message    string
	HttpStatus int
}

var (
	ErrImageNotFound     = AppError{Code: "ERR_01", Message: "Image not found", HttpStatus: 404}
	ErrTooManyRequest    = AppError{Code: "ERR_02", Message: "Too many requests, please wait", HttpStatus: 429}
	ErrGeminiTimeout     = AppError{Code: "ERR_03", Message: "Image generation timed out", HttpStatus: 504}
	ErrInvalidPrompt     = AppError{Code: "ERR_04", Message: "Invalid prompt", HttpStatus: 422}
	ErrGeminiBadResponse = AppError{Code: "ERR_05", Message: "Failed to generate image", HttpStatus: 502}
	ErrStorageFailed     = AppError{Code: "ERR_06", Message: "Failed to store image", HttpStatus: 500}
	ErrDatabaseFailed    = AppError{Code: "ERR_07", Message: "Database operation failed", HttpStatus: 500}
	ErrProjectNotFound   = AppError{Code: "ERR_08", Message: "Project not found", HttpStatus: 404}
	ErrUnauthorized      = AppError{Code: "ERR_09", Message: "Unauthorized", HttpStatus: 401}
	ErrInternalServer    = AppError{Code: "ERR_10", Message: "Internal server error", HttpStatus: 500}
	ErrValidationFailed  = AppError{Code: "ERR_11", Message: "Validation failed", HttpStatus: 422}
	ErrDuplicateResource = AppError{Code: "ERR_12", Message: "Resource already exists", HttpStatus: 409}
)
