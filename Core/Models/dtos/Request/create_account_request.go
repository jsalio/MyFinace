package dtos

// CreateAccountRequest representa la estructura de la solicitud para crear una cuenta
// swagger:model
// @name CreateAccountRequest
type CreateAccountRequest struct {
	Nick     string `json:"nick" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
