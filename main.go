package main

import (
	"flag"
	"fmt"

	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
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

func exists(name string) error {
	if _, err := os.Stat(name); err != nil {
		panic(err)
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

//template -- template/ output outout/ vars example/vars.yml
//template is A DIRECTORY, and output IS A DIRECTORY
// func parseDirectoryTemplates(tempName string, outName, string, vars string) {
//
// 	//read each file and parse with output file
// 	//iterate through tempName files
//
// 	if files, err := ioutil.ReadDir(tempName); err != nil {
// 		log.Println(err)
// 		return
// 	}else {
// 		//loop each file and parse
// 		for _, file := range files {
// 			//
// 			if tout, err := template.ParseFiles(file.Name()); err != nil {
// 				log.Println(err)
// 				return
// 			}else {
// 				//do something with tout
// 				//create output file
// 				var fileTo = file.Name()
//
// 				//apend .txt to fileTo
// 				newFile := strings.Join(fileTo, ".txt")
//
// 				os.Chdir(outName)
//
// 				f, err := os.Create(newFile)
// 				if err != nil {
// 					log.Println("Create file: ", err)
// 					return
// 				}
// 				defer f.Close()
//
// 				workDir := os.Getenv("PWD")
//
// 				os.Chdir(workDir)
//
// 			//Read YAML ------------------------------------------
// 				bytes, err := ioutil.ReadFile(vars)
// 				if err != nil {
// 					log.Print(err)
// 					return
// 				}
// 				m := make(map[string]interface{})
// 				err = yaml.Unmarshal(bytes, m)
// 				if err != nil {
// 					log.Print("yaml: ", err)
// 					return
// 				}
// 				//Execute template ----------------------------------
// 				err = tout.Execute(f, m)
// 				if err != nil {
// 					log.Print("execute: ", err)
// 					return
// 				}
// 				color.Green("Done")
// 			}
// 		}
// 	}
// }

//template -- template/demo.tmpl output outout/ vars example/vars.yml
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

//function to return if template is a directory - isTemplateDir(*t)
func isTemplateDir(tempName string) bool {
	tempDir, err := isDir(tempName)
	if err != nil {
		log.Println(err)
		return false
	}
	return tempDir
}

//function to return if output is a directory - isOutputDir(*output)
func isOutputDir(outName string) bool {
	outDir, err := isDir(outName)
	if err != nil {
		log.Println(err)
		return false
	}
	return outDir
}

func CreateOutputDirectoryIfDoesntExist(output string) error {
	if _, err := os.Stat(output); os.IsNotExist(err) {
		if err := os.Mkdir(output, 0700); err != nil {
			return err
		}
	}

	return nil
}

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

	if err := checkForExistenceOfRequiredDirectories(t, output); err != nil {
		log.Fatal(err)
	}

	// check if template(*t) is a file or a directory
	tempDir := isTemplateDir(*t)

	//check if *output is a directory
	outDir := isOutputDir(*output)

	//PARSE SINGLE FILE
	switch tempDir {
	case true:
		if outDir {
			//T is a directory, O is a directory;\s
			//call function to parse template directory files and create each file for output
			//parseDirectoryTemplates(*t, *output, *vars)
			log.Println("Do nothing for now")
		} else {
			//output is not a directory, and template is a file- do nothing for now
			log.Println("O file, and T is a file.Do nothing for now")
		}
	default:
		//output is a directory, T is a file,  parse it
		if outDir {
			//call function to parse single template file
			if err := ParseSingleTemplate(*t, *output, *vars); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkForExistenceOfRequiredDirectories(t *string, output *string) error {
	if err := exists(*t); err != nil {
		return err
	}

	err := CreateOutputDirectoryIfDoesntExist(*output)
	if err != nil {
		return err

	}

	return nil
}
