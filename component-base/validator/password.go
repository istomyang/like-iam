package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
	"unicode"
)

func WithCheckPassword() Option {
	return func(v *validator.Validate) error {
		return v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
			return checkPassword(fl.Field().String()) == nil
		})
	}
}

func CheckPassword(pw string) bool {
	return checkPassword(pw) == nil
}

func CheckPasswordErr(pw string) error {
	return checkPassword(pw)
}

// CheckPassword checks password directly.
func (v *Validator) CheckPassword(pw string) error {
	return checkPassword(pw)
}

const (
	minPassWordLength = 8
	maxPassWordLength = 16
)

func checkPassword(pw string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasPunctuation bool
	var hasSpace bool

	var length = 0
	var invalids string

	var appendInvalidMessage = func(message string) {
		msg := strings.TrimSpace(message)
		if invalids == "" {
			invalids = msg
		} else {
			invalids = invalids + "," + message
		}
	}

	for _, w := range pw {
		switch {
		case unicode.IsLower(w):
			hasLower = true
			length += 1
		case unicode.IsUpper(w):
			hasUpper = true
			length += 1
		case unicode.IsNumber(w):
			hasNumber = true
			length += 1
		case unicode.IsPunct(w):
			hasPunctuation = true
			length += 1
		case unicode.IsSpace(w):
			hasSpace = true
		}
	}

	if !hasNumber {
		appendInvalidMessage("number letter missing.")
	}
	if !hasUpper {
		appendInvalidMessage("upper letter missing.")
	}
	if !hasLower {
		appendInvalidMessage("lower letter missing.")
	}
	if !hasPunctuation {
		appendInvalidMessage("punctuation letter missing.")
	}
	if hasSpace {
		appendInvalidMessage("space latter must remove.")
	}
	if length < minPassWordLength && length > maxPassWordLength {
		appendInvalidMessage(fmt.Sprintf("password length must between %d and %d", minPassWordLength, maxPassWordLength))
	}

	if invalids != "" {
		return fmt.Errorf(invalids)
	}
	return nil
}
