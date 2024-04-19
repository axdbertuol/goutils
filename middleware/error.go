package custom_middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CustomError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	InternalCode string `json:"internal_code"`
}

// Error returns the error message.
func (ce *CustomError) Error() string {
	return fmt.Sprintf("%s: %s", ce.InternalCode, ce.Message)
}

// Custom error handler middleware
func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			// Check if the error is of type CustomError
			if customErr, ok := err.(*CustomError); ok {
				// Return custom error response
				return c.JSON(customErr.Code, customErr)
			}
			// If the error is not a CustomError, return a generic internal server error
			customErr := &CustomError{
				Code:         http.StatusInternalServerError,
				Message:      "Internal Server Error",
				InternalCode: "unexpectedError",
			}
			return c.JSON(http.StatusInternalServerError, customErr)
		}

		return nil
	}
}
