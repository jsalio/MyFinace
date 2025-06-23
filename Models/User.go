package models

import (
	"Financial/types"
	"time"
)

// User represents a user account in the financial application.
// This struct is used to store and manage user information throughout the system.
type User struct {
	// ID is the unique identifier for the user
	ID int `json:"id"`

	// Nickname is the user's chosen display name (required, unique)
	Nickname string `json:"nick_name"`

	// FirstName is the user's first name (required)
	FirstName string `json:"first_name"`

	// Lastname is the user's last name (required)
	Lastname string `json:"last_name"`

	// Email is the user's email address (required, unique)
	Email string `json:"email"`

	// Status represents the current state of the user's account
	Status types.AccountStatus `json:"status"`

	// CreatedAt is the timestamp when the user account was created
	CreatedAt time.Time `json:"created_at"`

	// Password is the hashed password for the user (never stored in plain text)
	Password string `json:"password"`
}

// UpdateAccountRequest represents the data that can be updated for a user account.
// This struct is used as a DTO (Data Transfer Object) for account updates.
type UpdateAccountRequest struct {
	// ID is the unique identifier of the user to update (required)
	ID int

	// FirstName is the updated first name (optional)
	FirstName string

	// Lastname is the updated last name (optional)
	Lastname string

	// Email is the updated email address (must remain unique)
	Email string

	// Status is the updated account status (optional)
	Status types.AccountStatus

	// Password is the new password (will be hashed before storage, optional)
	Password string
}
