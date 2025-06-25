// Package models contains the data structures used throughout the application.
// This file defines the Wallet structure and its related operations.
package models

import (
	"Financial/types"
)

// Wallet represents a user's digital wallet in the financial application.
// It contains the basic information and balance of a wallet.
type Wallet struct {
	// ID is the unique identifier for the user
	ID int `json:"id"`

	// Name is the user-defined identifier for the wallet.
	// It must be unique per user.
	Name string `json:"name"`

	// Type represents the kind of wallet (e.g., checking, savings, credit).
	// It uses the WalletType type defined in the types package.
	Type types.WalletType `json:"type"`

	// Balance is the current monetary amount available in the wallet.
	// It's represented as a float64 to support decimal values.
	Balance float64 `json:"balance"`

	// UserID is the foreign key that references the user who owns this wallet.
	// This field is required and must reference a valid user ID.
	UserID int `json:"userId"`

	// User is the navigation property to access the user who owns this wallet.
	// This field should be populated manually when needed.
	User *User `json:"user,omitempty"`
}
