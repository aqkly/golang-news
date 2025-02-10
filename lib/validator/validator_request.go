package validatorLib

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	var errorMessages []string

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "email":
				errorMessages = append(errorMessages, "Invalid email format")
			case "required":
				errorMessages = append(errorMessages, "Field "+err.Field()+" is required")
			case "min":
				if err.Field() == "Password" {
					errorMessages = append(errorMessages, "Password must be at least 8 characters")
				}
			case "eqfield":
				errorMessages = append(errorMessages, "Field "+err.Field()+" must be equal to "+err.Param())
			default:
				errorMessages = append(errorMessages, "Field "+err.Field()+" is invalid")
			}
		}

		return errors.New("Validate error: " + joinMessages(errorMessages))
	}

	return nil
}

func joinMessages(messages []string) string {
	result := ""

	for i, messages := range messages {
		if i > 0 {
			result += ", "
		}
		result += messages
	}

	return result
}
