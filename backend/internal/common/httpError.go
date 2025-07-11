package common

import (
	"fmt"
	"net/http"
)

type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	// General errors
	ErrInvalidInput        = &APIError{Status: http.StatusBadRequest, Code: "INVALID_INPUT", Message: "Invalid input"}
	ErrUserAlreadyExists   = &APIError{Status: http.StatusConflict, Code: "USER_ALREADY_EXISTS", Message: "Email already exists"}
	ErrInternalServer      = &APIError{Status: http.StatusInternalServerError, Code: "INTERNAL_ERROR", Message: "Internal server error"}
	ErrUnauthorized        = &APIError{Status: http.StatusUnauthorized, Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrForbidden           = &APIError{Status: http.StatusForbidden, Code: "FORBIDDEN", Message: "Forbidden"}
	ErrNotFound            = &APIError{Status: http.StatusNotFound, Code: "NOT_FOUND", Message: "Not found"}
	ErrBadRequest          = &APIError{Status: http.StatusBadRequest, Code: "BAD_REQUEST", Message: "Bad request"}
	ErrUnprocessableEntity = &APIError{Status: http.StatusUnprocessableEntity, Code: "UNPROCESSABLE_ENTITY", Message: "Unprocessable entity"}
	ErrTooManyRequests     = &APIError{Status: http.StatusTooManyRequests, Code: "TOO_MANY_REQUESTS", Message: "Too many requests"}

	// Authentication specific errors
	ErrInvalidCredentials   = &APIError{Status: http.StatusUnauthorized, Code: "INVALID_CREDENTIALS", Message: "Invalid email or password"}
	ErrUserInactive         = &APIError{Status: http.StatusForbidden, Code: "USER_INACTIVE", Message: "User account is inactive"}
	ErrUserNotFound         = &APIError{Status: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: "User not found"}
	ErrTokenExpired         = &APIError{Status: http.StatusUnauthorized, Code: "TOKEN_EXPIRED", Message: "Token has expired"}
	ErrTokenInvalid         = &APIError{Status: http.StatusUnauthorized, Code: "TOKEN_INVALID", Message: "Token is invalid"}
	ErrTokenMalformed       = &APIError{Status: http.StatusUnauthorized, Code: "TOKEN_MALFORMED", Message: "Token is malformed"}
	ErrInvalidTokenFormat   = &APIError{Status: http.StatusUnauthorized, Code: "INVALID_TOKEN_FORMAT", Message: "Invalid authorization header format. Expected: Bearer <token>"}
	ErrInvalidTokenType     = &APIError{Status: http.StatusUnauthorized, Code: "INVALID_TOKEN_TYPE", Message: "Invalid token type. Access token required"}
	ErrTokenRequired        = &APIError{Status: http.StatusUnauthorized, Code: "TOKEN_REQUIRED", Message: "Token is required"}
	ErrTokenRefreshFailed   = &APIError{Status: http.StatusUnauthorized, Code: "TOKEN_REFRESH_FAILED", Message: "Failed to refresh token"}
	ErrRegistrationFailed   = &APIError{Status: http.StatusInternalServerError, Code: "REGISTRATION_FAILED", Message: "Failed to register user"}
	ErrAuthenticationFailed = &APIError{Status: http.StatusInternalServerError, Code: "AUTHENTICATION_FAILED", Message: "Failed to authenticate user"}
	ErrInvalidRequestBody   = &APIError{Status: http.StatusBadRequest, Code: "INVALID_REQUEST_BODY", Message: "Invalid request body"}

	// Validation specific errors
	ErrEmailAlreadyExists = &APIError{Status: http.StatusConflict, Code: "EMAIL_ALREADY_EXISTS", Message: "Email already exists"}
	ErrWeakPassword       = &APIError{Status: http.StatusBadRequest, Code: "WEAK_PASSWORD", Message: "Password must be at least 6 characters with uppercase, number and special character"}
	ErrValidationFailed   = &APIError{Status: http.StatusBadRequest, Code: "VALIDATION_FAILED", Message: "Validation failed"}
	ErrInvalidEmail       = &APIError{Status: http.StatusBadRequest, Code: "INVALID_EMAIL", Message: "Invalid email format"}
	ErrPasswordTooShort   = &APIError{Status: http.StatusBadRequest, Code: "PASSWORD_TOO_SHORT", Message: "Password must be at least 6 characters long"}
	ErrPasswordTooWeak    = &APIError{Status: http.StatusBadRequest, Code: "PASSWORD_TOO_WEAK", Message: "Password must contain uppercase, lowercase, number, and special character"}

	// User management errors
	ErrUserCreationFailed = &APIError{Status: http.StatusInternalServerError, Code: "USER_CREATION_FAILED", Message: "Failed to create user"}
	ErrUserUpdateFailed   = &APIError{Status: http.StatusInternalServerError, Code: "USER_UPDATE_FAILED", Message: "Failed to update user"}
	ErrUserDeleteFailed   = &APIError{Status: http.StatusInternalServerError, Code: "USER_DELETE_FAILED", Message: "Failed to delete user"}

	// Database transaction errors
	ErrTransactionFailed = &APIError{Status: http.StatusInternalServerError, Code: "TRANSACTION_FAILED", Message: "Database transaction failed"}

	// Workspace creation errors
	ErrWorkspaceCreateForbidden = &APIError{
		Status:  http.StatusForbidden,
		Code:    "WORKSPACE_CREATE_FORBIDDEN",
		Message: "Only super admin users can create workspaces",
	}
	ErrWorkspaceNameRequired = &APIError{
		Status:  http.StatusBadRequest,
		Code:    "WORKSPACE_NAME_REQUIRED",
		Message: "Workspace name is required",
	}
	ErrWorkspaceCreateFailed = &APIError{
		Status:  http.StatusInternalServerError,
		Code:    "WORKSPACE_CREATE_FAILED",
		Message: "Failed to create workspace",
	}
	ErrWorkspaceSlugGeneration = &APIError{
		Status:  http.StatusInternalServerError,
		Code:    "WORKSPACE_SLUG_GENERATION_FAILED",
		Message: "Failed to generate unique workspace slug",
	}
)
