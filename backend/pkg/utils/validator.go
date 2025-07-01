package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	alphaSpaceRegex = regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
)

func SetupCustomValidators(v *validator.Validate) {
	v.RegisterValidation("alpha_space", validateAlphaSpace)
}

func validateAlphaSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}
	return alphaSpaceRegex.MatchString(field)
}
