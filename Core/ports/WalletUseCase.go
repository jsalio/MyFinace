package ports

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
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
	CreateWallet(request dtos.CreateWalletRequest) (*db.Wallet, *response.ErrorResponse)

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
	UpdateWallet(request dtos.UpdateWalletRequest) (*db.Wallet, *response.ErrorResponse)

	// DeleteWallet removes a wallet by its ID
	//
	// Parameters:
	//   - walletID: ID of the wallet to delete
	//
	// Returns:
	//   - error: Error if deletion fails (e.g., wallet not found)
	DeleteWallet(walletID int) error

	// GetUserWallet retrieves wallet information for a specific user
	//
	// Parameters:
	//   - id:    The ID of the wallet to retrieve
	//   - email: Email of the user requesting the wallet (for authorization)
	//
	// Returns:
	//   - *response.UserWalletResponse: The wallet information including balance and transactions
	//   - error: Error if retrieval fails (e.g., wallet not found, unauthorized access)
	GetUserWallet(id int, email string) (*response.UserWalletResponse, *response.ErrorResponse)
}
