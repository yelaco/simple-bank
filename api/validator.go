package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/yelaco/simple-bank/util"
)

var ErrBindingError = errors.New("coulnd't bind any validators")

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}

func (server Server) bindValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)
		if err != nil {
			return fmt.Errorf("api.bindValidators: %w", err)
		}

		return nil
	}
	return fmt.Errorf("api.bindValidators: %w", ErrBindingError)
}
