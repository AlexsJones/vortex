package abstraction

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Base struct {
	Variables map[string]interface{}
}

func New() *Base {
	return &Base{
		Variables: map[string]interface{}{},
	}
}

func (b *Base) LoadVariables(variablepath string) error {
	if _, err := os.Stat(variablepath); os.IsNotExist(err) {
		return fmt.Errorf("%v is not a valid path", variablepath)
	}
	buff, err := ioutil.ReadFile(variablepath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buff, &(b.Variables))
}

func (b *Base) ProcessTemplates(templatePath, outputPath string) error {
	return nil
}
