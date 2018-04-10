package processor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

// Vortex container of information that are awesome and amazing
type Vortex struct {
	variables map[string]interface{}
}

// LoadVariables will read from a file path and load Vortex with the variables ready
func (v *Vortex) LoadVariables(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("%v is not a valid path", filepath)
	}
	buff, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buff, &(v.variables))
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
		readpath := path.Join(outputroot, file.Name())
		switch {
		case file.IsDir():
			outputroot = path.Join(outputroot, file.Name())
			if err := v.ProcessTemplates(readpath, outputroot); err != nil {
				return err
			}
		default:
			// If the file extension doesn't match what we expect then ignore it
			if !regexp.MustCompile("$[a-zA-Z0-9]+\\.ya?ml^").MatchString(file.Name()) {
				continue
			}
			if err = v.processTemplate(readpath, outputroot); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Vortex) processTemplate(templatepath, outputpath string) error {
	// if the folder path doesn't exist, then we need to make it
	if _, err := os.Stat(path.Dir(outputpath)); os.IsNotExist(err) {
		if err = os.MkdirAll(path.Dir(outputpath), 0755); err != nil {
			return err
		}
	}
	if _, err := os.Stat(outputpath); !os.IsNotExist(err) {
		return fmt.Errorf("%v already exists, needs to be removed in order to process", outputpath)
	}
	buff, err := ioutil.ReadFile(templatepath)
	if err != nil {
		return err
	}
	tmpl, err := template.New("templated_var").Parse(string(buff))
	if err != nil {
		return err
	}
	writer := bytes.NewBuffer(nil)
	if err = tmpl.Execute(writer, v.variables); err != nil {
		return err
	}
	return ioutil.WriteFile(path.Join(outputpath, path.Base(templatepath)), writer.Bytes(), 0644)
}
