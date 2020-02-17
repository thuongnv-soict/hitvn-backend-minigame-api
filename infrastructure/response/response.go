package response

import (
	"g-tech.com/gerror"
)

/**
 * Defines a response object
 */
type Response struct {
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

/**
 * Defines an error response object
 */
type ErrorResponse struct {
	ErrorCode int	 `json:"ErrorCode"`
	Message   string `json:"Message"`
	Exception string `json:"Exception"`
}

/**
 * Returns a new error response object
 */
func NewErrorResponse(errorCode int, message string, exception string) (string, ErrorResponse) {
	msg := gerror.T(errorCode)
	error := ErrorResponse{
		ErrorCode: errorCode,
		Message:   message,
		Exception: exception,
	}

	return msg, error
}
