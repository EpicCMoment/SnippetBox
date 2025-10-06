package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}


func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	_, exists := v.FieldErrors[key]

	if !exists {
		v.FieldErrors[key] = message
	}

}

func (v *Validator) CheckField(ok bool, key, message string) {

	if !ok {
		v.AddFieldError(key, message)
	}

}

func NotBlank(field string) bool {
	return strings.TrimSpace(field) != ""
}

func BelowMaxChars(field string, n int) bool {
	return utf8.RuneCountInString(field) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {

	for _, permVal := range permittedValues {
		if value == permVal {
			return true
		}
	}

	return false

}