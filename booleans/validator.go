package booleans

import (
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	EPOCH_GT_NOW = "epoch-gt-now"
)

var _customValidator *validator.Validate

func NewCustomValidator() *validator.Validate {
	if _customValidator == nil {
		_customValidator = validator.New(validator.WithRequiredStructEnabled())

		if err := _customValidator.RegisterValidation(EPOCH_GT_NOW, validateEpochGreaterNow); err != nil {
			log.Println(err)

			log.Fatalf("failed to register custom validator: %s", EPOCH_GT_NOW)
		}
	}

	return _customValidator
}

func validateEpochGreaterNow(fl validator.FieldLevel) bool {
	var epoch int64

	now := time.Now().Unix()

	switch t := fl.Field().Interface().(type) {
	case int64:
		epoch = t
	default:
		return false
	}

	return epoch > now
}

func CustomValidateStruct(s interface{}) (err error) {
	if err = NewCustomValidator().Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)

			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			log.Printf("[%s] %s = %v\n", err.StructField(), err.Tag(), err.Value())

			return fmt.Errorf("[%s] invalid value \"%s\" for tag: %s", err.StructField(), err.Value(), err.Tag())
		}
	}

	return
}
