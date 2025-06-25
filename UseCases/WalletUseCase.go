package usecases

import (
	"Financial/Domains/ports"
	models "Financial/Models"
	"Financial/infrastructure"
	"Financial/types"
	"errors"
	"fmt"
	"strings"
)

// WalletUseCase implements the WalletUseCase interface
type WalletUseCase struct {
	repository ports.Repository[models.Wallet, int]
}

// NewWalletUseCase creates a new instance of WalletUseCase
func NewWalletUseCase(repo ports.Repository[models.Wallet, int]) ports.WalletUseCase {
	return &WalletUseCase{
		repository: repo,
	}
}

// CreateWallet implements WalletUseCase.CreateWallet
func (uc *WalletUseCase) CreateWallet(name string, walletType types.WalletType, balance float64, userID int) (*models.Wallet, error) {
	// Input validations
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("wallet name cannot be empty")
	}

	if balance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Check if wallet name already exists for this user
	existingWallets, err := uc.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error checking wallet existence: %w", err)
	}

	for _, w := range existingWallets {
		if w.Name == name && w.UserID == userID {
			return nil, errors.New("a wallet with this name already exists for this user")
		}
	}

	wallet := &models.Wallet{
		Name:    name,
		Type:    walletType,
		Balance: balance,
		UserID:  userID,
	}

	return uc.repository.Create(wallet)
}

// UpdateWallet implements WalletUseCase.UpdateWallet
func (uc *WalletUseCase) UpdateWallet(walletID int, name string, walletType *types.WalletType, balance *float64) (*models.Wallet, error) {
	// Input validation
	if walletID <= 0 {
		return nil, errors.New("invalid wallet ID")
	}

	// Get existing wallet
	existingWallet, err := uc.repository.FindByField("id", walletID)
	if err != nil {
		if err == infrastructure.ErrNotFound {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("error fetching wallet: %w", err)
	}

	// Update fields if provided
	updated := false

	if name != "" && name != existingWallet.Name {
		// Check if new name is already taken by another wallet of the same user
		existingWallets, err := uc.repository.GetAll()
		if err != nil {
			return nil, fmt.Errorf("error checking wallet names: %w", err)
		}

		for _, w := range existingWallets {
			if w.Name == name && w.UserID == existingWallet.UserID && w.ID != walletID {
				return nil, errors.New("a wallet with this name already exists for this user")
			}
		}

		existingWallet.Name = name
		updated = true
	}

	if walletType != nil && *walletType != existingWallet.Type {
		existingWallet.Type = *walletType
		updated = true
	}

	if balance != nil && *balance != existingWallet.Balance {
		if *balance < 0 {
			return nil, errors.New("balance cannot be negative")
		}
		existingWallet.Balance = *balance
		updated = true
	}

	if !updated {
		return existingWallet, nil // No changes made
	}

	return uc.repository.Update(existingWallet)
}

// DeleteWallet implements WalletUseCase.DeleteWallet
func (uc *WalletUseCase) DeleteWallet(walletID int) error {
	if walletID <= 0 {
		return errors.New("invalid wallet ID")
	}

	// Check if wallet exists
	_, err := uc.repository.GetByID(walletID)
	if err != nil {
		if err == infrastructure.ErrNotFound {
			return errors.New("wallet not found")
		}
		return fmt.Errorf("error fetching wallet: %w", err)
	}

	// In a real application, you might want to check if the wallet has any transactions
	// before allowing deletion

	return uc.repository.Delete(walletID)
}
