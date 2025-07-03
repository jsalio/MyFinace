package engine

import "errors"

var (
	ErrInvalidType      = errors.New("tipo de dato inválido")
	ErrFieldNotFound    = errors.New("campo no encontrado")
	ErrValidationFailed = errors.New("validación fallida")
)
