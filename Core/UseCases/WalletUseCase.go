package usecases

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
	"Financial/Core/ports"
	"Financial/Core/types"
	"Financial/Core/validators"
	"errors"
	"fmt"
	"strings"
)

// WalletUseCase implements the WalletUseCase interface
type WalletUseCase struct {
	repository ports.Repository[db.Wallet, int]
}

// NewWalletUseCase creates a new instance of WalletUseCase
func NewWalletUseCase(repo ports.Repository[db.Wallet, int]) ports.WalletUseCase {
	return &WalletUseCase{
		repository: repo,
	}
}

// CreateWallet implements WalletUseCase.CreateWallet
func (uc *WalletUseCase) CreateWallet(request dtos.CreateWalletRequest) (*db.Wallet, *response.ErrorResponse) {
	// Input validations
	// if strings.TrimSpace(request.Name) == "" {
	// 	return nil, errors.New("wallet name cannot be empty")
	// }

	// if request.Balance < 0 {
	// 	return nil, errors.New("initial balance cannot be negative")
	// }

	// if request.UserID <= 0 {
	// 	return nil, errors.New("invalid user ID")
	// }

	success, error := validators.ValidateWallet(request)
	if !success {
		return nil, &response.ErrorResponse{
			Error: strings.Join(*error, " \n"),
		}
	}

	// Check if wallet name already exists for this user
	existingWallets, err := uc.repository.GetAll()
	if err != nil {
		return nil, &response.ErrorResponse{
			Error: fmt.Errorf("error checking wallet existence: %w", err).Error(),
		}
	}

	for _, w := range existingWallets {
		if w.Name == request.Name && w.UserID == request.UserID {
			return nil, &response.ErrorResponse{
				Error: errors.New("a wallet with this name already exists for this user").Error(),
			}
		}
	}

	wallet := db.Wallet{
		Name:    request.Name,
		Type:    request.WalletType,
		Balance: request.Balance,
		UserID:  request.UserID,
	}
	result, err := uc.repository.Create(&wallet)

	if err != nil {
		return nil, &response.ErrorResponse{
			Error: errors.New("a wallet with this name already exists for this user").Error(),
		}
	}

	return result, nil
}

// UpdateWallet implements WalletUseCase.UpdateWallet
func (uc *WalletUseCase) UpdateWallet(request dtos.UpdateWalletRequest) (*db.Wallet, *response.ErrorResponse) {
	// Input validation
	// if request.WalletID <= 0 {
	// 	return nil, errors.New("invalid wallet ID")
	// }

	errorsVal, existingWallet := validators.UpdateWalletValidator(request, uc.repository)

	if errorsVal != nil {
		return nil, &response.ErrorResponse{
			Error: strings.Join(*errorsVal, " \n"),
		}
	}

	// Get existing wallet
	//existingWallet, err := uc.repository.FindByField("id", request.WalletID)
	// if err != nil {
	// 	if err == types.ErrNotFound {
	// 		return nil, &response.ErrorResponse{
	// 			Error: errors.New("wallet not found").Error(),
	// 		}
	// 	}
	// 	return nil, &response.ErrorResponse{
	// 		Error: fmt.Errorf("error fetching wallet: %w", err).Error(),
	// 	}
	// }

	// Update fields if provided
	updated := false

	if request.Name != "" && request.Name != existingWallet.Name {
		// Check if new name is already taken by another wallet of the same user
		existingWallets, err := uc.repository.GetAll()
		if err != nil {
			return nil, &response.ErrorResponse{
				Error: fmt.Errorf("error checking wallet names: %w", err).Error(),
			}
		}

		for _, w := range existingWallets {
			if w.Name == request.Name && w.UserID == existingWallet.UserID && w.ID != request.WalletID {
				return nil, &response.ErrorResponse{
					Error: errors.New("a wallet with this name already exists for this user").Error(),
				}
			}
		}

		existingWallet.Name = request.Name
		updated = true
	}

	if request.WalletType != nil && request.WalletType != &existingWallet.Type {
		existingWallet.Type = *request.WalletType
		updated = true
	}

	if request.Balance != nil && *request.Balance != existingWallet.Balance {
		if *request.Balance < 0 {
			return nil, &response.ErrorResponse{
				Error: errors.New("balance cannot be negative").Error(),
			}
		}
		existingWallet.Balance = *request.Balance
		updated = true
	}

	if !updated {
		return existingWallet, nil // No changes made
	}

	result, errorUpdate := uc.repository.Update(existingWallet)

	if errorUpdate != nil {
		return nil, &response.ErrorResponse{
			Error: errorUpdate.Error(),
		}
	}
	return result, nil

}

// DeleteWallet implements WalletUseCase.DeleteWallet
func (uc *WalletUseCase) DeleteWallet(walletID int) error {
	if walletID <= 0 {
		return errors.New("invalid wallet ID")
	}

	// Check if wallet exists
	_, err := uc.repository.GetByID(walletID)
	if err != nil {
		if err == types.ErrNotFound {
			return errors.New("wallet not found")
		}
		return fmt.Errorf("error fetching wallet: %w", err)
	}

	// In a real application, you might want to check if the wallet has any transactions
	// before allowing deletion

	return uc.repository.Delete(walletID)
}

func (uc *WalletUseCase) GetUserWallet(id int, email string) (*response.UserWalletResponse, *response.ErrorResponse) {
	data, err := uc.repository.Query("id,name,type,balance,user:users!inner(email)", ports.QueryOptions{
		Filters: []ports.Filter{
			ports.Filter{
				Field:    "users.email",
				Operator: "eq",
				Value:    email,
			},
		},
	})
	if err != nil {
		return nil, &response.ErrorResponse{
			Error: err.Error(),
		}
	}

	var result response.UserWalletResponse

	// Type assert the result to *ports.UserWallet
	wallet, ok := data.([]db.Wallet)
	if !ok {
		return nil, &response.ErrorResponse{
			Error: errors.New("unexpected type returned from repository").Error(),
		}
	}

	result = response.UserWalletResponse{
		Email: email,
	}

	if len(wallet) == 0 {
		result.Wallets = []struct {
			Name    string           `json:"name"`
			Type    types.WalletType `json:"type"`
			Balance float64          `json:"balance"`
		}{}
		return &result, nil
	}

	result.Email = wallet[0].User.Email
	for _, w := range wallet {
		result.Wallets = append(result.Wallets, struct {
			Name    string           "json:\"name\""
			Type    types.WalletType "json:\"type\""
			Balance float64          "json:\"balance\""
		}{
			Name:    w.Name,
			Type:    w.Type,
			Balance: w.Balance,
		})
	}

	return &result, nil
}
