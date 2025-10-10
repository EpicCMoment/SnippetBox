package validator

import (
	"net/mail"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
}


func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
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

func (v *Validator) AddNonFieldError(message string) {

	v.NonFieldErrors = append(v.NonFieldErrors, message)

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

	return slices.Contains(permittedValues, value) 

}

func IsValidEmail(email string) bool {

	_, err := mail.ParseAddress(email)

	if err != nil {
		return false
	}

	return true


}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}