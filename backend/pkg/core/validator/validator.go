package validator

import (
	"errors"
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/logger"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	govalidator "github.com/go-playground/validator/v10"
	"reflect"
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

		gov.RegisterCustomTypeFunc(ValidateUUID, uuid.OrderedUUID{})

		/*err := gov.RegisterValidation("ouuid", IsValidUUID)
		if err != nil {
			logger.Get().Error(err)
		}*/

		v = &defaultValidator{validate: gov}
	}

	return v

}

func (v *defaultValidator) Struct(s interface{}) (err error) {
	defer func() {
		if errRecovered := recover(); errRecovered != nil {
			logger.Get().Error(errRecovered)
			err = errors.New(fmt.Sprintf("unexpected core validator error: %v", errRecovered))
		}
	}()

	err = v.validate.Struct(s)

	return err
}

func ValidateUUID(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(uuid.OrderedUUID); ok {
		val, err := valuer.Value()
		if err == nil {
			return val
		}
	}
	return nil
}

func IsValidUUID(fl govalidator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uuid.OrderedUUID)
	if !ok {
		return false
	}

	if !id.Valid {
		return false
	}

	return true
}
