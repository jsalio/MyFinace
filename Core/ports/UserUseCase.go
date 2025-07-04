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
	CreateAccount(nick string, email string, password string) (*response.SuccessResponse[*response.CreateAccountResponse], *[]response.ErrorResponse)

	// DestroyAccount permanently deletes a user account identified by email.
	//
	// Parameters:
	//   - email: The email of the account to be deleted
	//
	// Returns:
	//   - *response.ErrorResponse: Error response if account deletion fails (e.g., account not found, permission denied)
	DestroyAccount(email string) *response.ErrorResponse

	// UpdateAccount modifies an existing user's account information.
	//
	// Parameters:
	//   - user: A db.UpdateAccountRequest containing the fields to be updated
	//
	// Returns:
	//   - *response.SuccessResponse[*response.UpdateAccountResponse]: Wrapped success response containing the updated account details
	//   - *response.ErrorResponse: Error response if update fails (e.g., invalid data, user not found)
	UpdateAccount(user db.UpdateAccountRequest) (*response.SuccessResponse[*response.UpdateAccountResponse], *response.ErrorResponse)

	// Login authenticates a user with the provided credentials.
	//
	// Parameters:
	//   - auth: An AuthRequest containing the user's login credentials
	//
	// Returns:
	//   - *string: A JWT token if authentication is successful
	//   - error: Error if authentication fails (e.g., invalid credentials)
	Login(auth dtos.AuthRequest) (*string, error)
}
