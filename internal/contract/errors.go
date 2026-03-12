package contract

const (
	AuthRequired    = "AUTH_REQUIRED"
	AuthInvalid     = "AUTH_INVALID"
	NotFound        = "NOT_FOUND"
	ValidationError = "VALIDATION_ERROR"
	Conflict        = "CONFLICT"
	RateLimited     = "RATE_LIMITED"
	HTTPError       = "HTTP_ERROR"
	FileNotFound    = "FILE_NOT_FOUND"
	Unsupported     = "UNSUPPORTED"
	UnknownError    = "UNKNOWN_ERROR"
)

// ContractError represents a structured CLI error with a machine-readable code.
type ContractError struct {
	Code    string
	Message string
}

func (e *ContractError) Error() string {
	return e.Code + ": " + e.Message
}

// NewError creates a new ContractError.
func NewError(code, message string) error {
	return &ContractError{Code: code, Message: message}
}
