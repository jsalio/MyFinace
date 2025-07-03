package response

// CreateAccountResponse represents the response structure returned after successfully
// creating a new account in the system. This DTO (Data Transfer Object) contains
// the essential account information that should be returned to the client.
//
// This response is typically used in account creation endpoints to provide
// confirmation and details of the newly created account.
//
// Example response:
//
//	{
//	  "id": 123,
//	  "nick": "johndoe",
//	  "email": "john.doe@example.com"
//	}
//
// swagger:model CreateAccountResponse
// @name CreateAccountResponse
// @description Response containing details of a newly created account
// @property {integer} id - The unique identifier for the created account
// @property {string} nick - The display name or username of the account
// @property {string} email - The email address associated with the account
type CreateAccountResponse struct {
	// ID is the unique identifier for the created account.
	// Example: 123
	// Required: true
	// minimum: 1
	ID int `json:"id" binding:"required"`

	// Nick is the display name or username for the account.
	// Example: "johndoe"
	// Required: true
	// minLength: 3
	// maxLength: 50
	Nick string `json:"nick" binding:"required"`

	// Email is the email address associated with the account.
	// Example: "user@example.com"
	// Required: true
	// format: email
	Email string `json:"email" binding:"required"`
}
