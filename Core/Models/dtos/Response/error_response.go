package dtos

// ErrorResponse representa una respuesta de error estándar
// swagger:model
// @name ErrorResponse
type ErrorResponse struct {
	Error string `json:"error"`
}
