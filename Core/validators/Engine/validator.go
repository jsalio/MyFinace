package engine

import (
	"reflect"
)

// ValidationResult contiene los resultados de la validación
type ValidationResult struct {
	Errors []ValidationError
}

// IsValid indica si la validación fue exitosa
func (r ValidationResult) IsValid() bool {
	return len(r.Errors) == 0
}

// ValidatorEngine es el motor de validación
type ValidatorEngine struct {
	rules []ValidationRule
}

// NewValidator crea una nueva instancia del motor
func NewValidator() *ValidatorEngine {
	return &ValidatorEngine{}
}

// AddRule agrega una regla de validación para un campo
func (v *ValidatorEngine) AddRule(fieldName string, rule RuleType, expected interface{}) {
	v.rules = append(v.rules, ValidationRule{
		FieldName: fieldName,
		Rule:      rule,
		Expected:  expected,
	})
}

func (v *ValidatorEngine) AddRules(fieldName string, rules []PatialValidationRule) {
	for _, rule := range rules {
		v.AddRule(fieldName, rule.Rule, rule.Expected)
	}
}

// Validate valida un struct y devuelve los resultados
func (v *ValidatorEngine) Validate(data interface{}) ValidationResult {
	result := ValidationResult{}
	val := reflect.ValueOf(data)

	// Asegurarse de que es un struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "",
			Rule:    "",
			Message: ErrInvalidType.Error(),
		})
		return result
	}

	// Iterar sobre las reglas definidas
	for _, rule := range v.rules {
		field := val.FieldByName(rule.FieldName)
		if !field.IsValid() {
			result.Errors = append(result.Errors, ValidationError{
				Field:   rule.FieldName,
				Rule:    rule.Rule,
				Message: ErrFieldNotFound.Error(),
			})
			continue
		}

		// Aplicar la regla al campo
		if err := rule.Validate(field.Interface()); err != nil {
			result.Errors = append(result.Errors, *err)
		}
	}

	return result
}
