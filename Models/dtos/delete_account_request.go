package dtos

// DeleteAccountRequest representa la estructura de la solicitud para eliminar una cuenta
// swagger:model
// @name DeleteAccountRequest
type DeleteAccountRequest struct {
	ID    int    `json:"id" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}
