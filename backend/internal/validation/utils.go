package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"guthub.com/gopeshwark/go-boilerplate/internal/errs"
)

type Validatable interface {
	Validate() error
}

type CustomValidationError struct {
	Field   string
	Message string
}

type CustomValidationErrors []CustomValidationError

func (c CustomValidationErrors) Error() string {
	return "Validation failed"
}

func BindAndValidate(c echo.Context, payload Validatable) error {
	if err := c.Bind(payload); err != nil {
		message := strings.Split(strings.Split(err.Error(), ",")[1], "message=")[1]
		return errs.NewBadRequestError(message, false, nil, nil, nil)
	}

	if msg, fieldErrors := validateStruct(payload); fieldErrors != nil {
		return errs.NewBadRequestError(msg, true, nil, fieldErrors, nil)
	}

	return nil
}

func BindAndValidateQuery(c echo.Context, payload Validatable) error {
	if err := c.Bind(payload); err != nil {
		return errs.NewBadRequestError("Invalid query parameters", false, nil, nil, nil)
	}

	if msg, fieldErrors := validateStruct(payload); fieldErrors != nil {
		return errs.NewBadRequestError(msg, true, nil, fieldErrors, nil)
	}

	return nil
}

func validateStruct(v Validatable) (string, []errs.FieldError) {
	if err := v.Validate(); err != nil {
		return extractValidationErrors(err)
	}

	return "", nil
}

func extractValidationErrors(err error) (string, []errs.FieldError) {
	var fieldErrors []errs.FieldError
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		CustomValidationErrors := err.(CustomValidationErrors)
		for _, err := range CustomValidationErrors {
			fieldErrors = append(fieldErrors, errs.FieldError{
				Field: err.Field,
				Error: err.Message,
			})
		}
	}

	for _, err := range validationErrors {
		field := strings.ToLower(err.Field())
		var msg string

		switch err.Tag() {
		case "required":
			msg = "is required"
		case "min":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must be at least %s characters", err.Params())
			} else {
				msg = fmt.Sprintf("must be at least %s", err.Params())
			}
		case "max":
			if err.Type().Kind() == reflect.String {
				msg = fmt.Sprintf("must not exceed %s characters", err.Params())
			} else {
				msg = fmt.Sprintf("must not exceed %s", err.Params())
			}
		case "oneof":
			msg = fmt.Sprintf("must be one of: %s", err.Params())
		case "email":
			msg = "must be a valid email address"
		case "e164":
			msg = "must be a valid phone number with country code"
		case "uuid":
			msg = "must be a valid UUID"
		case "uuidList":
			msg = "must be comma-seperated list of valid UUIDs"
		case "dive":
			msg = "some items are invalid"
		default:
			if err.Params() != "" {
				msg = fmt.Sprintf("%s: %s:%s", field, err.Tag(), err.Params())
			} else {
				msg = fmt.Sprintf("%s: %s", field, err.Tag())
			}
		}

		fieldErrors = append(fieldErrors, errs.FieldError{
			Field: strings.ToLower(err.Field()),
			Error: msg,
		})
	}

	return "Validation failed", fieldErrors
}

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

func IsValidUUID(uuid string) bool {
	return uuidRegex.MatchString(uuid)
}
