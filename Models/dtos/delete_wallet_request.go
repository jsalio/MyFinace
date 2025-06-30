package dtos

// DeleteAccountRequest representa la estructura de la solicitud para eliminar una cuenta
// swagger:model
// @name DeleteAccountRequest
type DeleteWalletRequest struct {
	ID    int    `json:"id" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}
