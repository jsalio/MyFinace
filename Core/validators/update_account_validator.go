package validators

import (
	"Financial/Core/Models/db"
	"Financial/Core/ports"
	"Financial/Core/types"
	engine "Financial/Core/validators/Engine"
)

func UpdateAccountValidator(data db.UpdateAccountRequest, repo ports.Repository[db.User, int]) *engine.ValidationResult {

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

	detectIsStatusInEnum := func(value interface{}) (bool, string) {
		// Convert the input value to AccountStatus
		status, ok := value.(types.AccountStatus)
		if !ok {
			return false, "Value is not of type AccountStatus"
		}

		// Check if the value matches any of the defined enum constants
		switch status {
		case types.Active, types.Inactive, types.Pending, types.Suspend:
			return true, ""
		default:
			return false, "Invalid AccountStatus value"
		}
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
	statuRules := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: "Value is not define"},
		{Rule: engine.Must, Expected: engine.CustomValidatorFunc(detectIsStatusInEnum), Message: "Value is not valid"},
	}

	validator.AddRules("Email", emailrules)
	validator.AddRules("Password", passwordRule)
	validator.AddRules("Status", statuRules)

	errors := validator.Validate(data)
	return &errors
}
