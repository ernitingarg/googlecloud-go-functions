package entity

type ErrorCode int

const (
	AUTH_ERROR       ErrorCode = 40100
	VALIDATION_ERROR ErrorCode = 40000
	SERVER_ERROR     ErrorCode = 50000
)

type ErrorResponse struct {
	ID    ErrorCode `json:"id"`
	Error string    `json:"error"`
}

func NewErrorResponse(id ErrorCode, msg string) *ErrorResponse {
	return &ErrorResponse{ID: id, Error: msg}
}
