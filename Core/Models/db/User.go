package db

import (
	"Financial/Core/types"
	"encoding/json"
	"fmt"
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
	CreatedAt time.Time `json:"created_at,omitempty"`

	// created_at_str is used internally for JSON unmarshaling
	createdAtStr string `json:"-"`

	// Password is the hashed password for the user (never stored in plain text)
	Password string `json:"password"`
}

// UnmarshalJSON implements the json.Unmarshaler interface to handle custom timestamp parsing
func (u *User) UnmarshalJSON(data []byte) error {
	// Define an auxiliary type to avoid recursion
	type Alias User
	aux := &struct {
		*Alias
		CreatedAtStr string `json:"created_at"`
	}{
		Alias: (*Alias)(u),
	}

	// Unmarshal the data into the auxiliary struct
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the timestamp string into a time.Time value
	if aux.CreatedAtStr != "" {
		// Try parsing with the format from the error message first
		parsedTime, err := time.Parse("2006-01-02T15:04:05.999", aux.CreatedAtStr)
		if err != nil {
			// If that fails, try the RFC3339 format
			parsedTime, err = time.Parse(time.RFC3339, aux.CreatedAtStr)
			if err != nil {
				return fmt.Errorf("error parsing created_at timestamp: %v", err)
			}
		}
		u.CreatedAt = parsedTime
	}

	return nil
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
