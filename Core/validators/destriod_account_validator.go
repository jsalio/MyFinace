package validators

import (
	"Financial/Core/Models/db"
	"Financial/Core/ports"
	engine "Financial/Core/validators/Engine"
)

func DestroidAccountValidator(email string, repo ports.Repository[db.User, int]) *engine.ValidationResult {

	detectValidMail := func(value interface{}) (bool, string) {
		result, error := repo.FindByField("email", value)
		if error != nil {
			return false, "Error fectching data"
		}
		if result == nil {
			return false, "Invalid Mail"
		}
		return true, ""
	}

	validator := engine.NewValidator()
	emailrules := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: "Email Is Empty"},
		{Rule: engine.ShouldMinLength, Expected: 12, Message: "Email not have length"},
		{Rule: engine.ShouldMatch, Expected: `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, Message: "Email not match"},
		{Rule: engine.Must, Expected: engine.CustomValidatorFunc(detectValidMail), Message: "Email not exists"},
	}
	validator.AddRules("Email", emailrules)
	errors := validator.Validate(email)
	return &errors
}
