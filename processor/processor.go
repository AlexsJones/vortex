package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Vortex container of information that are awesome and amazing
type Vortex struct {
	variables map[string]interface{}
	strict    bool
}

func New() *Vortex {
	return &Vortex{}
}

// LoadVariables will read from a file path and load Vortex with the variables ready
func (v *Vortex) LoadVariables(variablepath string) error {
	if _, err := os.Stat(variablepath); os.IsNotExist(err) {
		return fmt.Errorf("%v is not a valid path", variablepath)
	}
	buff, err := ioutil.ReadFile(variablepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buff, &(v.variables))
}

func (v *Vortex) EnableStrict() *Vortex {
	v.strict = true
	return v
}

// ProcessTemplates applys a DFS over the templateroot and will process the
// templates with the stored vortex variables
func (v *Vortex) ProcessTemplates(templateroot, outputroot string) error {
	// If the folder path doesn't exist, then say so
	// If the templateroot is a file, just process that
	root, err := os.Stat(templateroot)
	if os.IsNotExist(err) {
		return fmt.Errorf("%v does not exist", templateroot)
	}
	if !root.IsDir() {
		return v.processTemplate(templateroot, outputroot)
	}
	files, err := ioutil.ReadDir(templateroot)
	if err != nil {
		return err
	}
	for _, file := range files {
		readpath := path.Join(templateroot, file.Name())
		switch {
		case file.IsDir():
			newroot := path.Join(outputroot, file.Name())
			if err := v.ProcessTemplates(readpath, newroot); err != nil {
				return err
			}
		default:
			// If the file extension doesn't match what we expect then ignore it
			if err = v.processTemplate(readpath, outputroot); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Vortex) processTemplate(templatepath, outputpath string) error {
	if !strings.HasSuffix(templatepath, ".yaml") {
		return nil
	}
	// if the folder path doesn't exist, then we need to make it
	if _, err := os.Stat(outputpath); os.IsNotExist(err) && outputpath != "" {
		if err = os.MkdirAll(outputpath, 0755); err != nil {
			return err
		}
	}
	if f, err := os.Stat(outputpath); !os.IsNotExist(err) && !f.IsDir() {
		return fmt.Errorf("%v already exists, needs to be removed in order to process", outputpath)
	}
	buff, err := ioutil.ReadFile(templatepath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(path.Base(templatepath)).Parse(string(buff))
	if err != nil {
		return err
	}
	if v.strict {
		tmpl = tmpl.Option("missingkey=error")
	}
	writer := bytes.NewBuffer(nil)
	if err = tmpl.Execute(writer, v.variables); err != nil {
		return err
	}

	// Don't write the file if we have been told to validate only
	if !v.strict {
		filename := path.Join(outputpath, path.Base(templatepath))
		return ioutil.WriteFile(filename, writer.Bytes(), 0644)
	}
	// ensure that we have a valid yaml file at the end of it
	return yaml.UnmarshalStrict(writer.Bytes(), map[string]interface{}{})
}
