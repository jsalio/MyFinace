package dtos

// SuccessResponse representa una respuesta de éxito estándar
// swagger:model
// @name SuccessResponse
type SuccessResponse[TObject any] struct {
	Message string  `json:"message"`
	Data    TObject `json:"data"`
}
