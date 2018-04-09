package main

import (
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"

	"gopkg.in/yaml.v2"
	"path/filepath"
)

/*********************************************************************************
*     File Name           :     main.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-09-26 18:35]
**********************************************************************************/
var t *string
var vars *string
var output *string

func main() {

	t := flag.String("template", "", "path to template to populate")
	vars := flag.String("varpath", "", "path to var yaml to populate")
	output := flag.String("output", "", "name of output file")
	flag.Parse()

	if *t == "" || *vars == "" || *output == "" {
		fmt.Println("vortex is a simple program to combine a template with a yaml file of defined varibles it uses golang {{.var}} format with standard yaml")
		flag.Usage()
		return
	}

	if err := InputParametersCheck(t, output, vars); err != nil {
		log.Fatal(err)
	}

	if isDirectoryOfTemplates(*t) {
		if err := ParseDirectoryTemplates(*t, *output, *vars); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := ParseSingleTemplate(*t, *output, *vars); err != nil {
			log.Fatal(err)
		}
	}

}

func ParseSingleTemplate(tempName string, out string, vars string) error {
	tout, err := template.ParseFiles(tempName)
	if err != nil {
		return err
	}

	outputFileName := filepath.Base(tempName)

	if err := os.Chdir(out); err != nil {
		return err
	}

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	workDir := os.Getenv("PWD")
	os.Chdir(workDir)

	bytes, err := ioutil.ReadFile(vars)
	if err != nil {
		return err
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(bytes, m)
	if err != nil {
		return err
	}

	err = tout.Execute(outputFile, m)
	if err != nil {
		return err
	}

	color.Green("Done")

	return nil
}

func ParseDirectoryTemplates(tempDirectory string, outDirectory string, vars string) error {
	files, err := ioutil.ReadDir(tempDirectory)
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := ParseSingleTemplate(fmt.Sprint(tempDirectory, "/", f.Name()), outDirectory, vars); err != nil {
			return err
		}
	}

	return nil
}

func InputParametersCheck(t *string, output *string, v *string) error {
	if err := exists(*t); err != nil {
		return err
	}

	if err := exists(*v); err != nil {
		return err
	}
	err := createOutputDirectoryIfDoesntExist(*output)
	if err != nil {
		return err

	}

	return nil
}

func exists(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return err
	}

	return nil
}

func isDir(name string) (bool, error) {
	path, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	if !path.IsDir() {
		return false, nil
	}
	return true, nil
}

func isDirectoryOfTemplates(tempName string) bool {
	tempDir, err := isDir(tempName)
	if err != nil {
		log.Println(err)
		return false
	}
	return tempDir
}

func createOutputDirectoryIfDoesntExist(output string) error {
	if _, err := os.Stat(output); os.IsNotExist(err) {
		if err := os.Mkdir(output, 0700); err != nil {
			return err
		}
	}

	return nil
}
