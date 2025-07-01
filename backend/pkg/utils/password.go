package utils

import (
	"fmt"
	"go-backend-v2/internal/common"

	"golang.org/x/crypto/bcrypt"
)

const (
	BcryptCost = 10
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidatePassword(password string) error {
	if len(password) < common.PasswordMinLength {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_TOO_SHORT",
			Message: common.ErrMsgPasswordTooShort,
		}
	}

	if len(password) > common.PasswordMaxLength {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_TOO_LONG",
			Message: common.ErrMsgPasswordTooLong,
		}
	}

	if common.PasswordRequiredUpper && !common.PasswordUpperRegex.MatchString(password) {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_NO_UPPER",
			Message: common.ErrMsgPasswordNoUpper,
		}
	}

	if common.PasswordRequiredLower && !common.PasswordLowerRegex.MatchString(password) {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_NO_LOWER",
			Message: common.ErrMsgPasswordNoLower,
		}
	}

	if common.PasswordRequiredDigit && !common.PasswordDigitRegex.MatchString(password) {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_NO_DIGIT",
			Message: common.ErrMsgPasswordNoDigit,
		}
	}

	if common.PasswordRequiredSpecial && !common.PasswordSpecialRegex.MatchString(password) {
		return &common.APIError{
			Status:  400,
			Code:    "PASSWORD_NO_SPECIAL",
			Message: common.ErrMsgPasswordNoSpecial,
		}
	}

	return nil
}
