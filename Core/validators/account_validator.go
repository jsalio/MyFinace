package validators

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	ErrEmailRequired      = "email cannot be empty"
	ErrEmailInvalid       = "invalid email format"
	ErrEmailExists        = "email already exists"
	ErrNickRequired       = "nickname cannot be empty"
	ErrNickExists         = "nickname already exists"
	ErrPasswordRequired   = "password cannot be empty"
	ErrAccountNotFound    = "account not found"
	ErrInvalidCredentials = "invalid credentials"
)

type Validator struct {
	errors []string
}

func NewValidator() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

func (v *Validator) Required(value, field string) {
	if strings.TrimSpace(value) == "" {
		v.errors = append(v.errors, fmt.Sprintf("%s cannot be empty", field))
	}
}

func (v *Validator) IsValid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Error() string {
	return strings.Join(v.errors, ", ")
}

func (v *Validator) AddError(message string) {
	v.errors = append(v.errors, message)
}

func IsValidEmail(email string) bool {
	// Basic email regex: allows alphanumeric, dots, and common domains
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
