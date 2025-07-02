package dtos

// UpdateAccountResponse represents the response structure returned after successfully
// updating an existing account in the system. This DTO (Data Transfer Object) contains
// the updated account information that should be returned to the client.
//
// This response is typically used in account update endpoints to provide
// confirmation of the changes made to the account.
//
// Example response:
//
//	{
//	  "id": 123,
//	  "email": "updated.email@example.com"
//	}
//
// Note: The email field is optional in updates (omitempty) but must be a valid email when provided.
//
// swagger:model UpdateAccountResponse
// @name UpdateAccountResponse
// @description Response containing details of an updated account
// @property {integer} id - The unique identifier for the updated account
// @property {string} email - [Optional] The updated email address for the account
type UpdateAccountResponse struct {
	// ID is the unique identifier for the updated account.
	// This field is always required and cannot be empty.
	// Example: 123
	// Required: true
	// minimum: 1
	ID int `json:"id" binding:"required"`

	// Email is the updated email address for the account.
	// This field is optional during updates but must be a valid email format when provided.
	// Example: "updated.email@example.com"
	// Format: email
	Email string `json:"email" binding:"omitempty,email"`
}
