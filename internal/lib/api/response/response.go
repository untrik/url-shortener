package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}
func Error(msq string) Response {
	return Response{
		Status: StatusError,
		Error:  msq,
	}
}
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsq []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "requered":
			errMsq = append(errMsq, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsq = append(errMsq, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsq = append(errMsq, fmt.Sprintf("field %s is not valid", err.Field()))
		}

	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsq, ", "),
	}
}
