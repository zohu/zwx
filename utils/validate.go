package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(i any) error {
	if err := validate.Struct(i); err != nil {
		return fmt.Errorf("validate error: %v", err)
	}
	return nil
}
