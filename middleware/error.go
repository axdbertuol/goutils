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
				return echo.NewHTTPError(customErr.Code, &ErrorWrap{Error: customErr})
			} else if eHttpErr, ok := err.(*echo.HTTPError); ok {
				return echo.NewHTTPError(eHttpErr.Code, &ErrorWrap{Error: &types.CustomError{
					Code:         eHttpErr.Code,
					Message:      eHttpErr.Message.(string),
					InternalCode: "unexpectedHttpError",
				}})
			}

			customErr := &types.CustomError{
				Code:         http.StatusInternalServerError,
				Message:      "Internal Server Error",
				InternalCode: "unexpectedError",
			}

			return echo.NewHTTPError(http.StatusInternalServerError, &ErrorWrap{Error: customErr})
		}

		return nil
	}
}
