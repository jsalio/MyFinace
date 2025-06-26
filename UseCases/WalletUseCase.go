package usecases

import (
	"Financial/Domains/ports"
	"Financial/Models/db"
	"Financial/Models/dtos"
	"Financial/infrastructure"
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
func (uc *WalletUseCase) CreateWallet(request dtos.CreateWalletRequest) (*db.Wallet, error) {
	// Input validations
	if strings.TrimSpace(request.Name) == "" {
		return nil, errors.New("wallet name cannot be empty")
	}

	if request.Balance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}

	if request.UserID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	// Check if wallet name already exists for this user
	existingWallets, err := uc.repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error checking wallet existence: %w", err)
	}

	for _, w := range existingWallets {
		if w.Name == request.Name && w.UserID == request.UserID {
			return nil, errors.New("a wallet with this name already exists for this user")
		}
	}

	wallet := &db.Wallet{
		Name:    request.Name,
		Type:    request.WalletType,
		Balance: request.Balance,
		UserID:  request.UserID,
	}

	return uc.repository.Create(wallet)
}

// UpdateWallet implements WalletUseCase.UpdateWallet
func (uc *WalletUseCase) UpdateWallet(request dtos.UpdateWalletRequest) (*db.Wallet, error) {
	// Input validation
	if request.WalletID <= 0 {
		return nil, errors.New("invalid wallet ID")
	}

	// Get existing wallet
	existingWallet, err := uc.repository.FindByField("id", request.WalletID)
	if err != nil {
		if err == infrastructure.ErrNotFound {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("error fetching wallet: %w", err)
	}

	// Update fields if provided
	updated := false

	if request.Name != "" && request.Name != existingWallet.Name {
		// Check if new name is already taken by another wallet of the same user
		existingWallets, err := uc.repository.GetAll()
		if err != nil {
			return nil, fmt.Errorf("error checking wallet names: %w", err)
		}

		for _, w := range existingWallets {
			if w.Name == request.Name && w.UserID == existingWallet.UserID && w.ID != request.WalletID {
				return nil, errors.New("a wallet with this name already exists for this user")
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
			return nil, errors.New("balance cannot be negative")
		}
		existingWallet.Balance = *request.Balance
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
