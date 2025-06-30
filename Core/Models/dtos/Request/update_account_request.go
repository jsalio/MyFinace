package dtos

// UpdateAccountRequest representa la estructura de la solicitud para actualizar una cuenta
// swagger:model
// @name UpdateAccountRequest
type UpdateAccountRequest struct {
	ID        int    `json:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"omitempty,email"`
	Status    string `json:"status"`
	Password  string `json:"password"`
}
