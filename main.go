package main

import (
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"

	"path/filepath"

	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
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
	validate := flag.Bool("validate", false, "validate syntax and check for required variables")
	flag.Parse()

	if *validate {
		if *t == "" || *vars == "" {
			fmt.Println("To validate a file with vortex, pass a template file and variable file")
			flag.Usage()
			return
		}

		inputFilesAreValid, err := InputFilesAreValid(*t, *vars)
		if err != nil {
			log.Fatal(err)
		}

		if inputFilesAreValid {
			fmt.Println("template and var files are valid")
			return
		}

		os.Exit(1)
	}

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

func InputFilesAreValid(template string, varFile string) (bool, error) {
	if stat, err := os.Stat(template); err == nil && stat.IsDir() {
		isValid, err := varFileIsValid(varFile)
		if err != nil || !isValid {
			return false, err
		}

		var templates []string

		err = filepath.Walk(template, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			if strings.HasSuffix(info.Name(), ".yaml") {
				templates = append(templates, path)
				return nil
			}

			return nil
		})

		if err != nil {
			return false, err
		}

		for _, t := range templates {
			isValid, err := templateFileIsValid(t)
			if err != nil {
				return false, err
			}

			if !isValid {
				return false, nil
			}

			isValid, err = varFileHasExpectedVariables(t, varFile)
			if err != nil {
				return false, err
			}

			if !isValid {
				return false, nil
			}
		}

		return true, nil
	}

	isValid, err := varFileIsValid(varFile)
	if err != nil || !isValid {
		return false, err
	}

	isValid, err = templateFileIsValid(template)
	if err != nil {
		return false, err
	}

	if !isValid {
		return false, nil
	}

	isValid, err = varFileHasExpectedVariables(template, varFile)
	if err != nil {
		return false, err
	}

	if !isValid {
		return false, nil
	}

	return true, nil
}

func varFileHasExpectedVariables(templateFile string, varFile string) (bool, error) {
	varMap, err := readVars(varFile)
	if err != nil {
		return false, err
	}

	bytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return false, err
	}

	r := regexp.MustCompile(`{{\s{0,1}\.(.*?)\s{0,1}}}`)
	expectedVars := r.FindAllStringSubmatch(string(bytes), -1)

	for _, eVar := range expectedVars {
		if _, ok := varMap[eVar[1]]; !ok {
			log.Printf("Could not find variable: %s, required by %s, in %s", eVar[1], templateFile, varFile)
			return false, nil
		}
	}

	return true, nil
}

func templateFileIsValid(templateFile string) (bool, error) {
	bytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return false, err
	}

	r := regexp.MustCompile(`{{\s{0,1}\..*\s{0,1}}}`)
	bytes = r.ReplaceAll(bytes, []byte("# placeholder"))

	m := make(map[string]interface{})
	err = yaml.Unmarshal(bytes, m)
	if err != nil {
		log.Printf("Failed to validate syntax: %s", templateFile)
		return false, err
	}

	return true, nil
}

func varFileIsValid(varFile string) (bool, error) {
	if _, err := readVars(varFile); err != nil {
		log.Printf("Failed to validate syntax: %s", varFile)
		return false, err
	}

	return true, nil
}

func ParseSingleTemplate(templateFile string, outputFilePath string, varFilePath string) error {
	tout, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	workDir := os.Getenv("PWD")
	os.Chdir(workDir)

	varMap, err := readVars(varFilePath)

	err = tout.Execute(outputFile, varMap)
	if err != nil {
		return err
	}

	log.Printf("Creating: %s", outputFilePath)
	color.Green("Done")

	return nil
}

func readVars(varFile string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(varFile)
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(bytes, m)
	if err != nil {
		return nil, err
	}

	return m, err
}

func ParseDirectoryTemplates(tempDirectory string, outDirectory string, vars string) error {
	var templates []string
	var directories []string

	err := filepath.Walk(tempDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			directories = append(directories, path)
			return nil
		}

		if strings.HasSuffix(info.Name(), ".yaml") {
			templates = append(templates, path)
			return nil
		}

		return nil
	})

	if err != nil {
		return err
	}

	for _, dir := range directories {
		directoryPath := strings.Replace(dir, tempDirectory, outDirectory, -1)
		if err := createOutputDirectoryIfDoesntExist(directoryPath); err != nil {
			return err
		}
	}

	for _, t := range templates {
		outputFilePath := strings.Replace(t, tempDirectory, outDirectory, -1)
		if err := ParseSingleTemplate(t, outputFilePath, vars); err != nil {
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
		if err := os.MkdirAll(output, 0700); err != nil {
			return err
		}
	}

	return nil
}
