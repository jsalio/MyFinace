package ports

import (
	"Financial/Models/db"
	"Financial/Models/dtos"
)

// WalletUseCase defines the interface for wallet-related business logic operations.
type WalletUseCase interface {
	// CreateWallet creates a new wallet with the provided details
	//
	// Parameters:
	//   - name:      The name of the wallet (must be unique per user)
	//   - walletType: The type of wallet (e.g., checking, savings, credit)
	//   - balance:    Initial balance of the wallet (must be >= 0)
	//   - userID:     ID of the user who owns the wallet
	//
	// Returns:
	//   - *models.Wallet: The newly created wallet
	//   - error:         Error if creation fails (e.g., invalid data, duplicate name)
	CreateWallet(request dtos.CreateWalletRequest) (*db.Wallet, error)

	// UpdateWallet updates an existing wallet with new information
	//
	// Parameters:
	//   - walletID:  ID of the wallet to update
	//   - name:      New name for the wallet (optional)
	//   - walletType: New type for the wallet (optional)
	//   - balance:    New balance for the wallet (optional, must be >= 0)
	//
	// Returns:
	//   - *models.Wallet: The updated wallet
	//   - error:         Error if update fails (e.g., invalid data, wallet not found)
	UpdateWallet(request dtos.UpdateWalletRequest) (*db.Wallet, error)

	// DeleteWallet removes a wallet by its ID
	//
	// Parameters:
	//   - walletID: ID of the wallet to delete
	//
	// Returns:
	//   - error: Error if deletion fails (e.g., wallet not found)
	DeleteWallet(walletID int) error

	GetUserWallet(id int, email string) (*UserWallet, error)
}
