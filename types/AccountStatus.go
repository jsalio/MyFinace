package types

type AccountStatus string

const (
	Active   AccountStatus = "active"
	Inactive AccountStatus = "inactive"
	Pending  AccountStatus = "pending"
	Suspend  AccountStatus = "suspended"
)
