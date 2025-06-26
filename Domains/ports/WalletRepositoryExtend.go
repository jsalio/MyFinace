package ports

import "Financial/Models/db"

type UserWallet struct {
	// ID is the unique identifier for the user
	ID int `json:"id"`

	// Nickname is the user's chosen display name (required, unique)
	Nickname string `json:"nick_name"`

	// Email is the user's email address (required, unique)
	Email string `json:"email"`

	Wallets []db.Wallet `json:"wallets"`
}

type ExtendedRepository[T any, ID comparable] interface {
	Repository[T, ID]                                        // Incrustamos la interfaz Repository
	GetUserWallet(id int, email string) (*UserWallet, error) // Añadimos un nuevo método
}
