package validator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"sync"
)

type Option func(*validator.Validate) error

type Validator struct {
	validator *validator.Validate
}

var valid *Validator
var once sync.Once

func GetValidator(options ...Option) (*Validator, error) {
	if valid == nil && options == nil {
		return nil, fmt.Errorf("singleton hasn't created, options is missing")
	}

	once.Do(func() {
		valid = &Validator{validator: validator.New()}
		for _, option := range options {
			_ = option(valid.validator)
		}
	})
	return valid, nil
}

// Native return singleton validator.Validate, it's thread-safe.
// more see: https://pkg.go.dev/github.com/go-playground/validator/v10
func (v *Validator) Native() *validator.Validate {
	return v.validator
}
