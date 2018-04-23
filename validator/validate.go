package validator

import (
	"github.com/AlexsJones/vortex/abstraction"
)

type Validator struct {
	*abstraction.Base
}

func New() *Validator {
	return &Validator{
		Base: abstraction.New(),
	}
}

func (v *Validator) LoadVariables(variablepath string) error {
	return v.Base.LoadVariables(variablepath)
}

func (v *Validator) ProcessTemplates(templatePath, outputPath string) error {
	return nil
}
