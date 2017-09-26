package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/dimiro1/banner"
	"github.com/fatih/color"

	yaml "gopkg.in/yaml.v2"
)

/*********************************************************************************
*     File Name           :     main.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-09-26 18:35]
**********************************************************************************/
const b string = `

                  _
                 | |
 __   _____  _ __| |_ _____  __
 \ \ / / _ \| '__| __/ _ \ \/ /
  \ V / (_) | |  | ||  __/>  <
   \_/ \___/|_|   \__\___/_/\_\


`

func main() {
	banner.Init(os.Stdout, true, true, bytes.NewBufferString(b))

	var t = flag.String("template", "", "path to template to populate")
	var vars = flag.String("varpath", "", "path to var yaml to populate")
	var output = flag.String("output", "", "name of output file")
	flag.Parse()

	if *t == "" || *vars == "" || *output == "" {
		fmt.Println("vortex is a simple program to combine a template with a yaml file of defined varibles it uses golang {{.var}} format with standard yaml")
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
