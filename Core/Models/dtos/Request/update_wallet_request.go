package dtos

import "Financial/Core/types"

type UpdateWalletRequest struct {
	WalletID   int               `json:"id"`
	Name       string            `json:"name"`
	WalletType *types.WalletType `json:"type"`
	Balance    *float64          `json:"balance"`
}
