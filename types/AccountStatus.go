package types

type AccountStatus string

const (
	Active   AccountStatus = "active"
	Inactive AccountStatus = "inactive"
	Pending  AccountStatus = "pending"
	Suspend  AccountStatus = "suspended"
)

type WalletType string

const (
	Debit  WalletType = "Debit"
	Credit WalletType = "Credit"
)
