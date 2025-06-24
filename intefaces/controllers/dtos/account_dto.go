package dtos

// CreateAccountRequest representa la estructura de la solicitud para crear una cuenta
// swagger:model
// @name CreateAccountRequest
type CreateAccountRequest struct {
	Nick     string `json:"nick" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

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

// DeleteAccountRequest representa la estructura de la solicitud para eliminar una cuenta
// swagger:model
// @name DeleteAccountRequest
type DeleteAccountRequest struct {
	ID    int    `json:"id" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

// ErrorResponse representa una respuesta de error estándar
// swagger:model
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse representa una respuesta de éxito estándar
// swagger:model
// @name SuccessResponse
type SuccessResponse struct {
	Message string `json:"message"`
}
