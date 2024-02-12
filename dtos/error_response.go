package dtos

import "strings"

type ErrorResponse struct {
	Error string `json:"error"`
}

var DefaultErrorResponse = &ErrorResponse{
	Error: "Unknown",
}

func NewErrorResponse(err ...string) *ErrorResponse {
	errJoin := strings.Join(err, ":")
	if errJoin != "" {
		return &ErrorResponse{
			Error: errJoin,
		}
	}
	return DefaultErrorResponse
}
