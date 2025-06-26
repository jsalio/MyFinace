package dtos

import "Financial/types"

type UpdateWalletRequest struct {
	WalletID   int               `json:"id"`
	Name       string            `json:"name"`
	WalletType *types.WalletType `json:"type"`
	Balance    *float64          `json:"balance"`
}
