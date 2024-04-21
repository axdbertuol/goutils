package types

import "strconv"

type CustomError struct {
	Code         int    `json:"code"`
	Message      string `json:"message"`
	InternalCode string `json:"internal_code"`
}

// Error returns the error message.
func (ce *CustomError) Error() string {
	json := `{
		"message": "` + ce.Message + `",
		"success": false,
		"internalCode": "` + ce.InternalCode + `",
		"code": ` + strconv.Itoa(ce.Code) + `
	}`
	return json
}
