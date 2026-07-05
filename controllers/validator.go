package controllers

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func validateStruct(v any) error {
	return validate.Struct(v)
}
