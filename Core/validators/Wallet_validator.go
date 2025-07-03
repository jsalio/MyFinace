package validators

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	"Financial/Core/ports"
	"Financial/Core/types"
	engine "Financial/Core/validators/Engine"
	"fmt"
)

func ValidateWallet(data dtos.CreateWalletRequest) (bool, *[]string) {
	var errors []string

	validator := engine.NewValidator()
	validator.AddRule("Name", engine.ShouldNotEmpty, nil)
	validator.AddRule("Balance", engine.ShouldGreatOrEqualThah, 0)
	validator.AddRule("UserID", engine.ShouldGreatOrEqualThah, 0)

	result := validator.Validate(data)

	if result.IsValid() {
		return true, nil
	} else {
		fmt.Println("Errores de validación:")
		for _, err := range result.Errors {
			err := fmt.Errorf("Campo: %s, Regla: %s, Mensaje: %s\n", err.Field, err.Rule, err.Message)
			errors = append(errors, err.Error())
		}
	}
	return false, &errors
}

func UpdateWalletValidator(data dtos.UpdateWalletRequest, repository ports.Repository[db.Wallet, int]) (*[]string, *db.Wallet) {
	var errors []string

	validator := engine.NewValidator()
	validator.AddRule("Name", engine.ShouldNotEmpty, nil)
	validator.AddRule("Balance", engine.ShouldGreatOrEqualThah, 0)

	result := validator.Validate(data)

	if result.IsValid() {
		wallet, err := repository.FindByField("id", data.WalletID)
		if err != nil {
			if err == types.ErrNotFound {
				errors = append(errors, "wallet not found")
			} else {
				errors = append(errors, fmt.Errorf("error fetching wallet: %w", err).Error())
			}
		}
		return &errors, wallet
	} else {
		fmt.Println("Errores de validación:")
		for _, err := range result.Errors {
			err := fmt.Errorf("Campo: %s, Regla: %s, Mensaje: %s\n", err.Field, err.Rule, err.Message).Error()
			errors = append(errors, err)
		}
	}

	return &errors, nil
}
