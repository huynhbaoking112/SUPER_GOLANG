package common

import "regexp"

const (
	PasswordMinLength       = 6
	PasswordMaxLength       = 128
	PasswordRequiredUpper   = true
	PasswordRequiredLower   = true
	PasswordRequiredDigit   = true
	PasswordRequiredSpecial = true
)

var (
	PasswordUpperRegex   = regexp.MustCompile(`[A-Z]`)
	PasswordLowerRegex   = regexp.MustCompile(`[a-z]`)
	PasswordDigitRegex   = regexp.MustCompile(`[0-9]`)
	PasswordSpecialRegex = regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`)

	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	NameRegex = regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
)

const (
	ErrMsgPasswordTooShort  = "Password must be at least 6 characters long"
	ErrMsgPasswordTooLong   = "Password must not exceed 128 characters"
	ErrMsgPasswordNoUpper   = "Password must contain at least one uppercase letter"
	ErrMsgPasswordNoLower   = "Password must contain at least one lowercase letter"
	ErrMsgPasswordNoDigit   = "Password must contain at least one digit"
	ErrMsgPasswordNoSpecial = "Password must contain at least one special character"
)
