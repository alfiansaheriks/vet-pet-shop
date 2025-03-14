package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationErrors(errs validator.ValidationErrors) string {
	var errMessages []string

	for _, err := range errs {
		field := err.Field()
		tag := err.Tag()
		param := err.Param()

		switch tag {
		case "required":
			errMessages = append(errMessages, field+" is required")
		case "email":
			errMessages = append(errMessages, field+" must be a valid email address")
		case "gte":
			errMessages = append(errMessages, field+" must be greater than or equal to "+param)
		case "lte":
			errMessages = append(errMessages, field+" must be less than or equal to "+param)
		case "alphanum":
			errMessages = append(errMessages, field+" must be alphanumeric")
		case "eqfield":
			errMessages = append(errMessages, field+" must be equal to "+param)
		case "min":
			errMessages = append(errMessages, field+" must be at least "+param+" characters")
		case "oneof":
			errMessages = append(errMessages, field+" must be one of "+param)
		default:
			errMessages = append(errMessages, field+" is not valid")
		}
	}

	return strings.Join(errMessages, ", ")
}
