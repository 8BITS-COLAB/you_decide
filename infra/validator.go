package infra

import "github.com/go-playground/validator/v10"

type Validator struct {
	provider *validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		provider: validator.New(),
	}
}

func (v *Validator) ValidateStruct(params interface{}) error {
	if err := v.provider.Struct(params); err != nil {
		return err
	}

	return nil
}
