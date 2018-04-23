package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/AlexsJones/vortex/abstraction"
	"github.com/AlexsJones/vortex/processor"
	"github.com/AlexsJones/vortex/validator"
)

/*********************************************************************************
*     File Name           :     main.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-09-26 18:35]
**********************************************************************************/
const (
	usage string = `
Vortex -- a simplified template parser

The desired usage is to read from a variables file (defined in yaml)
and template in the variables into the given templates.
Thus, the usage of the progam is:

vortex --template path --varpath path [--validate] [--output path]
`
)

var (
	templatePath string
	variablePath string
	outputPath   string
	validate     bool
)

func init() {
	const (
		blank = ""
	)
	flag.StringVar(&templatePath, "template", blank, "path to template to use")
	flag.StringVar(&variablePath, "varpath", blank, "path to var yaml to populate")
	flag.StringVar(&outputPath, "output", blank, "Output path for the rendered templates to be outputted")
	flag.BoolVar(&validate, "validate", false, "validate syntax and check for the required variables")
}

func main() {
	flag.Parse()
	var (
		vortex abstraction.TemplateProcessor
	)
	switch {
	case variablePath != "" && templatePath != "":
		if validate {
			vortex = validator.New()
		} else {
			vortex = processor.New()
		}
	default:
		fmt.Println(usage)
		flag.Usage()
		return
	}
	if err := vortex.LoadVariables(variablePath); err != nil {
		fmt.Println("Unable to load files due to:", err)
		os.Exit(1)
	}
	if err := vortex.ProcessTemplates(templatePath, outputPath); err != nil {
		fmt.Println("Unable to process templates due to:", err)
		os.Exit(1)
	}
}
