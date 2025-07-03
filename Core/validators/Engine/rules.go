package engine

import (
	"fmt"
	"reflect"
	"regexp"
)

type RuleType string

const (
	ShouldEqual            RuleType = "Equal"
	ShouldNotEqual         RuleType = "NotEqual"
	ShouldGreatThah        RuleType = "GreatThat"
	ShouldGreatOrEqualThah RuleType = "GreatOrEqualThat"
	ShouldLessThat         RuleType = "LessThat"
	ShouldLessOrEqualThat  RuleType = "LessOrEqualThat"
	ShouldEmpty            RuleType = "Empty"
	ShouldNotEmpty         RuleType = "NotEmpty"
	ShouldMatch            RuleType = "Match"
	ShouldLength           RuleType = "Length"
)

type ValidationError struct {
	Field   string
	Rule    RuleType
	Message string
}

type ValidationRule struct {
	FieldName string
	Rule      RuleType
	Expected  interface{} // Valor esperado para la regla (puede ser string, int, regex, etc.)
}

func (r ValidationRule) Validate(value interface{}) *ValidationError {
	v := reflect.ValueOf(value)
	expected := reflect.ValueOf(r.Expected)

	switch r.Rule {
	case ShouldEqual:
		if !reflect.DeepEqual(value, r.Expected) {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("el campo %s debe ser igual a %v", r.FieldName, r.Expected),
			}
		}
	case ShouldNotEqual:
		if reflect.DeepEqual(value, r.Expected) {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("el campo %s no debe ser igual a %v", r.FieldName, r.Expected),
			}
		}
	case ShouldGreatThah:
		if v.Kind() == reflect.Int && expected.Kind() == reflect.Int {
			if v.Int() <= expected.Int() {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s debe ser mayor que %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	case ShouldGreatOrEqualThah:
		if v.Kind() == reflect.Int && expected.Kind() == reflect.Int {
			if v.Int() < expected.Int() {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s debe ser mayor o igual a %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	case ShouldLessThat:
		if v.Kind() == reflect.Int && expected.Kind() == reflect.Int {
			if v.Int() >= expected.Int() {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s debe ser menor que %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	case ShouldLessOrEqualThat:
		if v.Kind() == reflect.Int && expected.Kind() == reflect.Int {
			if v.Int() > expected.Int() {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s debe ser menor o igual a %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	case ShouldEmpty:
		if v.Kind() == reflect.String && v.String() != "" {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("el campo %s debe estar vacío", r.FieldName),
			}
		}
	case ShouldNotEmpty:
		if v.Kind() == reflect.String && v.String() == "" {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("el campo %s no debe estar vacío", r.FieldName),
			}
		}
	case ShouldMatch:
		if v.Kind() == reflect.String && expected.Kind() == reflect.String {
			if !regexp.MustCompile(expected.String()).MatchString(v.String()) {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s no coincide con el patrón %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	case ShouldLength:
		if v.Kind() == reflect.String && expected.Kind() == reflect.Int {
			if len(v.String()) != int(expected.Int()) {
				return &ValidationError{
					Field:   r.FieldName,
					Rule:    r.Rule,
					Message: fmt.Sprintf("el campo %s debe tener longitud %v", r.FieldName, r.Expected),
				}
			}
		} else {
			return &ValidationError{
				Field:   r.FieldName,
				Rule:    r.Rule,
				Message: fmt.Sprintf("tipo inválido para %s", r.FieldName),
			}
		}
	}
	return nil
}
