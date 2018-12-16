package processor

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/AlexsJones/vortex/secrets"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type vortex struct {
	variables map[string]interface{}
	strict    bool
	debug     bool
	validator string
	filter    *regexp.Regexp
}

func New() *vortex {
	return &vortex{
		variables: map[string]interface{}{},
		filter:    regexp.MustCompile(`.ya?ml$`),
		validator: "yaml",
	}
}

// Set allows the user to define variables as command line arguments
func (v *vortex) Set(input string) error {
	// Lets try Unmarshal a json like object and try key value pairs if all else fails
	if err := yaml.Unmarshal([]byte(input), &(v.variables)); err != nil {
		data := strings.Split(input, "=")
		// If we don't have a key value pair split by = then reject it
		if len(data) != 2 {
			return errors.New("Incorrect format, expect key=value")
		}
		v.variables[data[0]] = data[1]
	}
	return nil
}

func (v *vortex) String() string {
	return fmt.Sprintf("Vortex has %v loaded", len(v.variables))
}

// EnableDebug enables logging for Vortex
func (v *vortex) EnableDebug(enable bool) *vortex {
	v.debug = enable
	return v
}

func (v *vortex) SetValidator(validator string) *vortex {
	v.validator = validator
	return v
}

// LoadVariables will read from a file path and load Vortex with the variables ready
func (v *vortex) LoadVariables(variablepath string) error {
	if _, err := os.Stat(variablepath); os.IsNotExist(err) {
		// Possible that we have loaded variables already so
		// it is safe to continue
		if len(v.variables) != 0 && variablepath == "" {
			return nil
		}
		return fmt.Errorf("%v is not a valid path", variablepath)
	}
	buff, err := ioutil.ReadFile(variablepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buff, &(v.variables))
}

func (v *vortex) EnableStrict(enable bool) *vortex {
	v.strict = enable
	return v
}

func (v *vortex) SetFilter(filter string) *vortex {
	v.filter = regexp.MustCompile(filter)
	return v
}

func (v *vortex) canProcess(templatepath string) bool {
	return v.filter.MatchString(templatepath)
}

// ProcessTemplates applys a DFS over the templateroot and will process the
// templates with the stored vortex variables
func (v *vortex) ProcessTemplates(templateroot, outputroot string) error {
	// If the folder path doesn't exist, then say so
	// If the templateroot is a file, just process that
	v.logMessage("Output directory: ", outputroot)
	v.logMessage("Template directory: ", templateroot)

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
		// Ensure we don't automatically recurse down hidden files
		case file.IsDir():
			newroot := path.Join(outputroot, file.Name())
			if strings.HasPrefix(file.Name(), ".") {
				v.logMessage("Skipping hidden file: ", newroot)
				continue
			}
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

func (v *vortex) processTemplate(templatepath, outputpath string) error {
	if !v.canProcess(templatepath) {
		v.logMessage("Can not process: ", templatepath)
		return nil
	}
	// if the folder path doesn't exist, then we need to make it
	// and make sure we don't create a directory if we are just validating the contents
	if _, err := os.Stat(outputpath); os.IsNotExist(err) && !v.strict {
		v.logMessage("Creating the output directory as it doesn't exist yet")
		if err = os.MkdirAll(outputpath, 0755); err != nil {
			return err
		}
		v.logMessage(outputpath, "Directory now exists")
	}
	filename := path.Join(outputpath, path.Base(templatepath))
	if f, err := os.Stat(filename); !os.IsNotExist(err) && !f.IsDir() && !v.strict {
		return fmt.Errorf("%v already exists, needs to be removed in order to process", filename)
	}
	v.logMessage("Reading file: ", templatepath)
	buff, err := ioutil.ReadFile(templatepath)
	if err != nil {
		return err
	}
	tmpl, err := template.New(path.Base(templatepath)).
		Funcs(template.FuncMap{
			"vaultsecret": secrets.VaultFetchSecret,
			"getenv":      os.Getenv,
			"md5":         hashMd5,
			"base64Encode": base64Encode,
			"base64Decode": base64Decode,
		}).
		Parse(string(buff))
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
		v.logMessage("Attempting to write file to ", filename)
		if err := ioutil.WriteFile(filename, writer.Bytes(), 0644); err != nil {
			return err
		}
		v.logMessage("Successfully writen to: ", filename)
		return nil
	}
	return v.validate(writer.Bytes())
}

func (v *vortex) logMessage(args ...interface{}) {
	if v.debug {
		log.Info(args...)
	}
}

func (v *vortex) validate(content []byte) error {
	switch v.validator {
	case "yaml", "json":
		return yaml.Unmarshal(content, map[string]interface{}{})
	case "text":
		fallthrough
	default:
		v.logMessage("No additional file format validation set")
	}
	return nil
}
