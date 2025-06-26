package types

type WalletType string

const (
	Debit  WalletType = "Debit"
	Credit WalletType = "Credit"
)

func WalletTypePtr(wt WalletType) *WalletType {
	return &wt
}
