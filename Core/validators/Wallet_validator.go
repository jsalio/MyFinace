package validators

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	"Financial/Core/ports"
	"Financial/Core/types"
	engine "Financial/Core/validators/Engine"
	"fmt"
)

// ValidateWallet validates the CreateWalletRequest and returns validation results.
// Returns true with nil errors if valid, or false with a slice of error messages.
func ValidateWallet(data dtos.CreateWalletRequest) (bool, *[]string) {
	var errors []string

	validator := engine.NewValidator()

	nameRule := []engine.PatialValidationRule{
		{Rule: engine.ShouldNotEmpty, Expected: nil, Message: ""},
		{Rule: engine.ShouldLength, Expected: 5, Message: ""},
	}

	validator.AddRules("Name", nameRule)
	validator.AddRule("Balance", engine.ShouldGreaterOrEqualThan, 0, "") // Fixed typo
	validator.AddRule("UserID", engine.ShouldGreaterOrEqualThan, 0, "")  // Fixed typo

	result := validator.Validate(data)

	if result.IsValid() {
		return true, nil
	}

	for _, err := range result.Errors {
		// Avoid shadowing err; use formatted string directly
		errorMsg := fmt.Sprintf("Field: %s, Rule: %s, Message: %s", err.Field, err.Rule, err.Message)
		errors = append(errors, errorMsg)
	}

	return false, &errors
}

// UpdateWalletValidator validates the UpdateWalletRequest and checks if the wallet exists.
// Returns a slice of error messages and the wallet if found, or nil if not found or invalid.
func UpdateWalletValidator(data dtos.UpdateWalletRequest, repository ports.Repository[db.Wallet, int]) (*[]string, *db.Wallet) {
	var errors []string

	validator := engine.NewValidator()
	validator.AddRule("Name", engine.ShouldNotEmpty, nil, "")
	validator.AddRule("Balance", engine.ShouldGreaterOrEqualThan, 0, "") // Fixed typo

	result := validator.Validate(data)

	if !result.IsValid() {
		for _, err := range result.Errors {
			// Avoid shadowing err; use formatted string directly
			errorMsg := fmt.Sprintf("Field: %s, Rule: %s, Message: %s", err.Field, err.Rule, err.Message)
			errors = append(errors, errorMsg)
		}
		return &errors, nil
	}

	wallet, err := repository.FindByField("id", data.WalletID)
	if err != nil {
		if err == types.ErrNotFound {
			errors = append(errors, "wallet not found")
		} else {
			errors = append(errors, fmt.Sprintf("error fetching wallet: %v", err))
		}
		return &errors, nil
	}

	return &errors, wallet
}
