package validator

import "github.com/AlexsJones/vortex/utils"

type Validator struct {
	variables map[string]interface{}
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) LoadVariables(variablepath string) error {
	content, err := utils.ProgramaticMarshall(variablepath)
	v.variables = content
	return err
}

func (v *Validator) ProcessTemplates(templatePath, outputPath string) error {
	return nil
}
