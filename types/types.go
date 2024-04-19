package types

import "fmt"

type CustomError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	InternalCode string `json:"internal_code"`
}

// Error returns the error message.
func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", ce.InternalCode, ce.Message)
}
