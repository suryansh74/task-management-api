package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates any struct and returns field-error map
func ValidateStruct(data interface{}) map[string]string {
	errors := make(map[string]string)
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors[err.Field()] = MsgForTag(err)
		}
	}
	return errors
}

// GetValidator returns the validator instance for custom validations
func GetValidator() *validator.Validate {
	return validate
}
