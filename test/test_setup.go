package mocks

import "testing"

// TestSetup es una estructura genérica para configurar tests de casos de uso
// TModel: Tipo de la estructura de solicitud (request)
// TRepo: Tipo del modelo de base de datos
// ID: Tipo del ID del modelo (normalmente int, string, etc.)
type TestSetup[TModel any, TRepo any, ID comparable] struct {
	Name        string
	Req         TModel
	ExpectErr   bool
	ExpectedErr error
	SetupMock   func(*MockRepository[TRepo, ID])
	Verify      func(t *testing.T, result *TRepo, err error)
}
