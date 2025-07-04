package validators

import (
	repository "Financial/Core/Models/db"
	request "Financial/Core/Models/dtos/Request"
	contract "Financial/Core/ports"
	engine "Financial/Core/validators/Engine"
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

func CreateAccountValidator(data request.CreateAccountRequest, repo contract.Repository[repository.User, int]) *engine.ValidationResult {

	detectDuplicatedMail := func(value interface{}) (bool, string) {
		result, error := repo.FindByField("email", value)
		if error != nil {
			return false, "Error fectching data"
		}
		if result != nil {
			return false, "Duplicated Mail"
		}
		return true, ""
	}

	checkExistNick := func(value interface{}) (bool, string) {
		result, error := repo.FindByField("nick", value)
		if error != nil {
			return false, "Error fectching data"
		}
		if result != nil {
			return false, "Nickname already exists"
		}
		return true, ""
	}

	validator := engine.NewValidator()
	emailrules := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: "Email Is Empty"},
		{Rule: engine.ShouldMinLength, Expected: 12, Message: "Email not have length"},
		{Rule: engine.ShouldMatch, Expected: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, Message: "Email not match"},
		{Rule: engine.Must, Expected: engine.CustomValidatorFunc(detectDuplicatedMail), Message: "Email already exists"},
	}
	passwordRule := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: "Password Is Empty"},
		{Rule: engine.ShouldMatch, Expected: `^\S*$`, Message: "Password contains space"},
		{Rule: engine.ShouldMinLength, Expected: 8, Message: "Password not have length"},
	}
	nickNameRules := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: "Nickname Is Empty"},
		{Rule: engine.ShouldMinLength, Expected: 6, Message: "Nickname not have length"},
		{Rule: engine.ShouldMatch, Expected: `^\S*$`, Message: "Nickname contains space"},
		{Rule: engine.Must, Expected: engine.CustomValidatorFunc(checkExistNick), Message: "Nickname already exists"},
	}

	validator.AddRules("Email", emailrules)
	validator.AddRules("Password", passwordRule)
	validator.AddRules("Nick", nickNameRules)

	errors := validator.Validate(data)
	return &errors
}
