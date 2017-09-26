package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

/*********************************************************************************
*     File Name           :     main.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-09-26 18:35]
**********************************************************************************/

func main() {
	var t = flag.String("template", "", "path to template to populate")
	var vars = flag.String("varpath", "", "path to var yaml to populate")
	var output = flag.String("output", "", "name of output file")
	flag.Parse()

	if *t == "" || *vars == "" || *output == "" {
		flag.Usage()
		return
	}
	//Parse template -------------------------------------
	tout, err := template.ParseFiles(*t)
	if err != nil {
		log.Print(err)
		return
	}
	//Create output file ---------------------------------
	f, err := os.Create(*output)
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()

	//Read YAML ------------------------------------------
	bytes, err := ioutil.ReadFile(*vars)
	if err != nil {
		log.Print(err)
		return
	}
	m := make(map[string]string)
	err = yaml.Unmarshal(bytes, m)
	if err != nil {
		log.Print("yaml: ", err)
		return
	}
	//Execute template ----------------------------------
	err = tout.Execute(f, m)
	if err != nil {
		log.Print("execute: ", err)
		return
	}

	color.Green("Done")
}
