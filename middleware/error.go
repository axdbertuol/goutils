package custom_middleware

import (
	"net/http"

	"github.com/axdbertuol/goutils/types"
	"github.com/labstack/echo/v4"
)

type ErrorWrap struct {
	Error *types.CustomError `json:"error"`
}

// Custom error handler middleware
func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			// Check if the error is of type CustomError
			if customErr, ok := err.(*types.CustomError); ok {
				// Return custom error response
				return c.JSON(customErr.Code, &ErrorWrap{Error: customErr})
			}
			// If the error is not a CustomError, return a generic internal server error
			customErr := &types.CustomError{
				Code:         http.StatusInternalServerError,
				Message:      "Internal Server Error",
				InternalCode: "unexpectedError",
			}

			return c.JSON(http.StatusInternalServerError, &ErrorWrap{Error: customErr})
		}

		return nil
	}
}
