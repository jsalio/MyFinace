package models

// AuthError represents an authentication or authorization error response.
// It is used to standardize error responses related to authentication failures,
// such as invalid credentials, expired tokens, or insufficient permissions.
//
// This struct implements the error interface, allowing it to be used as a standard error type
// while also providing structured error information for API responses.
//
// Example usage:
//
//	return &AuthError{Message: "invalid or expired authentication token"}
//
// swagger:model AuthError
// @name AuthError
// @description Standardized authentication error response
// @property {string} message - Human-readable error message describing the authentication failure
type AuthError struct {
	// Message contains a human-readable description of the authentication error.
	// This message should be user-friendly and can be displayed directly to end users.
	// Example: "invalid or expired authentication token"
	// Required: true
	Message string `json:"message"`
}

// Error implements the error interface for AuthError.
// This allows AuthError to be used anywhere a standard error is expected.
func (e *AuthError) Error() string {
	return e.Message
}
