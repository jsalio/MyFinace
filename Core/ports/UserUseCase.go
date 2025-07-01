package ports

import (
	"Financial/Core/Models/db"
	dtos "Financial/Core/Models/dtos/Request"
	response "Financial/Core/Models/dtos/Response"
)

// UserUseCase defines the business logic operations for user account management.
// This interface serves as a contract for user-related use cases in the application.
type UserUseCase interface {
	// CreateAccount registers a new user account with the provided credentials.
	//
	// Parameters:
	//   - nick:     The user's nickname or display name
	//   - email:    The user's email address (must be unique)
	//   - password: The user's password (will be hashed before storage)
	//
	// Returns:
	//   - *response.SuccessResponse[*response.CreateAccountResponse]: Wrapped success response containing the created account details
	//   - *response.ErrorResponse: Error response if account creation fails (e.g., duplicate email, invalid input)
	CreateAccount(nick string, email string, password string) (*response.SuccessResponse[*response.CreateAccountResponse], *response.ErrorResponse)

	// DestroyAccount permanently deletes a user account identified by email.
	//
	// Parameters:
	//   - email: The email of the account to be deleted
	//
	// Returns:
	//   - error: Error if account deletion fails (e.g., account not found, permission denied)
	DestroyAccount(email string) *response.ErrorResponse

	// UpdateAccount modifies an existing user's account information.
	//
	// Parameters:
	//   - user: An UpdateAccountRequest containing the fields to be updated
	//
	// Returns:
	//   - *models.User: The updated user object
	//   - error:       Error if update fails (e.g., invalid data, user not found)
	UpdateAccount(user db.UpdateAccountRequest) (*response.SuccessResponse[*response.UpdateAccountResponse], *response.ErrorResponse)

	Login(auth dtos.AuthRequest) (*string, error)
}
