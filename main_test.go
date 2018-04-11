package main

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	TEMPLATE_FILE = "template.yaml"
	TEMPLATE_DIR  = "templates"
	VARS_FILE     = "vars.yaml"

	TEMPLATE       = "template"
	DIRECTORY      = "directory"
	DIRECTORY_TREE = "directory tree"
	VALID          = "valid"
	INVALID        = "invalid"
	MIXED          = "mixed"
	CONTAINS       = "contains"
	SUCCESSFUL     = "successful"
)

var (
	directoryValidTemplatePaths       = []string{"templates/template1.yaml", "templates/template2.yaml"}
	directoryInvalidTemplatePaths     = []string{"templates/template3.yaml"}

	directoryTreePaths                = []string{"templates/one", "templates/two"}
	directoryTreeValidTemplatePaths   = []string{"templates/one/template1.yaml", "templates/two/template2.yaml"}
	directoryTreeInvalidTemplatePaths = []string{"templates/template3.yaml"}
)

func (v *ValidationFeature) generateTemplates(validity, inputType string) error {
	switch inputType {
	case TEMPLATE:
		switch validity {
		case VALID:
			if err := ioutil.WriteFile(TEMPLATE_FILE, v.validTemplate, os.ModePerm); err != nil {
				return err
			}
		case INVALID:
			if err := ioutil.WriteFile(TEMPLATE_FILE, v.invalidTemplate, os.ModePerm); err != nil {
				return err
			}
		}
	case DIRECTORY:
		if err := os.MkdirAll(TEMPLATE_DIR, os.ModePerm); err != nil {
			return err
		}

		for _, validTemplate := range directoryValidTemplatePaths {
			if err := ioutil.WriteFile(validTemplate, v.validTemplate, os.ModePerm); err != nil {
				return err
			}
		}

		if validity == MIXED {
			for _, invalidTemplate := range directoryInvalidTemplatePaths {
				if err := ioutil.WriteFile(invalidTemplate, v.invalidTemplate, os.ModePerm); err != nil {
					return err
				}
			}
		}
	case DIRECTORY_TREE:
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

		if validity == MIXED {
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
	if validity == VALID {
		if err := ioutil.WriteFile(VARS_FILE, v.validVars, os.ModePerm); err != nil {
			return err
		}
	} else {
		if err := ioutil.WriteFile(VARS_FILE, v.invalidVars, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (v *ValidationFeature) validate(input string) error {
	var templateArg string

	if strings.Contains(input, DIRECTORY) {
		templateArg = TEMPLATE_DIR
	} else {
		templateArg = TEMPLATE_FILE
	}

	areValid, _ := InputFilesAreValid(templateArg, VARS_FILE)

	if !areValid {
		v.isValid = false
	}

	return nil
}

func (v *ValidationFeature) checkValidationResult(expected string) error {
	if expected == SUCCESSFUL {
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

type ValidationFeature struct {
	isValid         bool
	validTemplate   []byte
	invalidTemplate []byte
	validVars       []byte
	invalidVars     []byte
	missingVars     []byte
}

func (v *ValidationFeature) generateVarTemplatePairs(variables, input string) error {
	if variables == CONTAINS {
		if err := ioutil.WriteFile(VARS_FILE, v.validVars, os.ModePerm); err != nil {
			return err
		}
	} else {
		if err := ioutil.WriteFile(VARS_FILE, v.missingVars, os.ModePerm); err != nil {
			return err
		}
	}

	switch input {
	case DIRECTORY_TREE:
		v.generateTemplates(VALID, DIRECTORY)
	case DIRECTORY:
		v.generateTemplates(VALID, DIRECTORY_TREE)
	case TEMPLATE:
		v.generateTemplates(VALID, TEMPLATE)
	}
	return nil
}

func FeatureContext(s *godog.Suite) {
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
		os.RemoveAll(TEMPLATE_DIR)
		log.SetOutput(ioutil.Discard)
		v.isValid = true
	})

	s.AfterScenario(func(interface{}, error) {
		filesToCleanup := []string{TEMPLATE_FILE, VARS_FILE}

		for _, file := range filesToCleanup {
			os.Remove(file)
		}

		os.RemoveAll(TEMPLATE_DIR)
	})
}
