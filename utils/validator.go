package utils

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

var (
	phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)
	validate   = validator.New()
)

func init() {
	validate.RegisterValidation("phone", validatePhone)
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}
