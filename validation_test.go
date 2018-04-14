package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/DATA-DOG/godog"
	"log"
)


type ValidationFeature struct {
	isValid         bool
	validTemplate   []byte
	invalidTemplate []byte
	validVars       []byte
	invalidVars     []byte
	missingVars     []byte
}

func ValidationFeatureContext(s *godog.Suite) {
	v := &ValidationFeature{
		isValid:         true,
		validTemplate:   []byte("valid: {{.template}}"),
		invalidTemplate: []byte("invalid; {{.template}}"),
		validVars:       []byte("template: valid"),
		invalidVars:     []byte("template; invalid"),
		missingVars:     []byte("not-template: something else"),
	}

	s.Step(`^an? (valid|invalid|mixed) (template|directory|directory tree)$`, v.generateTemplates)
	s.Step(`^an? (valid|invalid) vars.yaml file$`, v.generateVarsFile)
	s.Step(`^vortex is run with the -validate flag for a (template|directory|directory tree)$`, v.validate)
	s.Step(`^validation should be (successful|unsuccessful)$`, v.checkValidationResult)
	s.Step(`^a vars.yaml that (contains|doesn't contain) all expected variables for a (template|directory|directory tree)$`, v.generateVarTemplatePairs)

	s.BeforeScenario(func(interface{}) {
		os.RemoveAll(TemplateDir)
		log.SetOutput(ioutil.Discard)
		v.isValid = true
	})

	s.AfterScenario(func(interface{}, error) {
		filesToCleanup := []string{TemplateFile, VarsFile}

		for _, file := range filesToCleanup {
			os.Remove(file)
		}

		os.RemoveAll(TemplateDir)
	})
}

func (v *ValidationFeature) generateTemplates(validity, inputType string) error {
	switch inputType {
	case Template:
		switch validity {
		case Valid:
			if err := ioutil.WriteFile(TemplateFile, v.validTemplate, os.ModePerm); err != nil {
				return err
			}
		case Invalid:
			if err := ioutil.WriteFile(TemplateFile, v.invalidTemplate, os.ModePerm); err != nil {
				return err
			}
		}
	case Directory:
		if err := os.MkdirAll(TemplateDir, os.ModePerm); err != nil {
			return err
		}

		for _, validTemplate := range directoryValidTemplatePaths {
			if err := ioutil.WriteFile(validTemplate, v.validTemplate, os.ModePerm); err != nil {
				return err
			}
		}

		if validity == Mixed {
			for _, invalidTemplate := range directoryInvalidTemplatePaths {
				if err := ioutil.WriteFile(invalidTemplate, v.invalidTemplate, os.ModePerm); err != nil {
					return err
				}
			}
		}
	case DirectoryTree:
		for _, dir := range directoryTreePaths {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}
		}

		for _, validTemplate := range directoryTreeValidTemplatePaths {
			if err := ioutil.WriteFile(validTemplate, v.validTemplate, os.ModePerm); err != nil {
				return err
			}
		}

		if validity == Mixed {
			for _, invalidTemplate := range directoryTreeInvalidTemplatePaths {
				if err := ioutil.WriteFile(invalidTemplate, v.invalidTemplate, os.ModePerm); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (v *ValidationFeature) generateVarsFile(validity string) error {
	if validity == Valid {
		if err := ioutil.WriteFile(VarsFile, v.validVars, os.ModePerm); err != nil {
			return err
		}
	} else {
		if err := ioutil.WriteFile(VarsFile, v.invalidVars, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (v *ValidationFeature) validate(input string) error {
	var templateArg string

	if strings.Contains(input, Directory) {
		templateArg = TemplateDir
	} else {
		templateArg = TemplateFile
	}

	areValid, _ := InputFilesAreValid(templateArg, VarsFile)

	if !areValid {
		v.isValid = false
	}

	return nil
}

func (v *ValidationFeature) checkValidationResult(expected string) error {
	if expected == Successful {
		if v.isValid {
			return nil
		} else {
			return fmt.Errorf("should have been valid")
		}
	} else {
		if v.isValid {
			return fmt.Errorf("should not have been valid")
		} else {
			return nil
		}
	}
}

func (v *ValidationFeature) generateVarTemplatePairs(variables, input string) error {
	if variables == Contains {
		if err := ioutil.WriteFile(VarsFile, v.validVars, os.ModePerm); err != nil {
			return err
		}
	} else {
		if err := ioutil.WriteFile(VarsFile, v.missingVars, os.ModePerm); err != nil {
			return err
		}
	}

	switch input {
	case DirectoryTree:
		v.generateTemplates(Valid, Directory)
	case Directory:
		v.generateTemplates(Valid, DirectoryTree)
	case Template:
		v.generateTemplates(Valid, Template)
	}
	return nil
}
