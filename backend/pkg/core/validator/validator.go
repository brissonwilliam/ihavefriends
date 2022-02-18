package validator

import (
	"errors"
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	govalidator "github.com/go-playground/validator/v10"
)

var v *defaultValidator

// Validator is for validating usecase data entering in
type Validator interface {
	Struct(interface{}) error
}

type defaultValidator struct {
	validate *govalidator.Validate
}

func Get() Validator {
	if v == nil {
		gov := govalidator.New()

		// register custom validations on gov here

		v = &defaultValidator{validate: gov}
	}

	return v

}

func (v defaultValidator) Struct(s interface{}) (err error) {
	defer func() {
		if errRecovered := recover(); errRecovered != nil {
			logger.Get().Error(errRecovered)
			err = errors.New(fmt.Sprintf("unexpected core validator error: %v", errRecovered))
		}
	}()

	err = v.validate.Struct(s)

	return err
}
