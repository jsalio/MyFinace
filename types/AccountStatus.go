package types

type AccountStatus string

const (
	Active   AccountStatus = "Active"
	Inactive AccountStatus = "Inactive"
	Pending  AccountStatus = "Pending"
)
