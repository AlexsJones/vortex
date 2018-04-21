package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func ProgramaticMarshall(variablepath string) (map[string]interface{}, error) {
	if _, err := os.Stat(variablepath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%v is not a valid path", variablepath)
	}
	buff, err := ioutil.ReadFile(variablepath)
	if err != nil {
		return nil, err
	}
	vars := map[string]interface{}{}
	if err := yaml.Unmarshal(buff, &vars); err != nil {
		return nil, err
	}
	return vars, nil
}
