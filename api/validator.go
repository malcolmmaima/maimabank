package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/malcolmmaima/maimabank/util"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	currency, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return util.IsSupportedCurrency(currency)
}