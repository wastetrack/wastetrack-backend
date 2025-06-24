package config

import (
	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

func NewValidator(viper *viper.Viper) *validator.Validate {
	return validator.New()
}
