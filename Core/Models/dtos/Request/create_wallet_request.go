package dtos

import "Financial/Core/types"

type CreateWalletRequest struct {
	Name       string           `json:"name"`
	WalletType types.WalletType `json:"type"`
	Balance    float64          `json:"balance"`
	UserID     int              `json:"accoundId"`
}
