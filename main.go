package main

import (
	"flag"
	"fmt"

	"os"

	"github.com/AlexsJones/vortex/processor"
	log "github.com/sirupsen/logrus"
)

/*********************************************************************************
*     File Name           :     main.go
*     Created By          :     jonesax
*     Creation Date       :     [2017-09-26 18:35]
**********************************************************************************/
const (
	usage string = `%s -- a simplified template parser

The desired usage is to read from a variables file (defined in yaml)
and template in the variables into the given templates.
Thus, the usage of the program is:

%s --template path --varpath path [--validate] [--output path] [--verbose]

The flags being used are:
`
)

var (
	templatePath string
	variablePath string
	outputPath   string
	debug        bool
	validate     bool
)

func init() {
	const (
		blank = ""
	)
	flag.StringVar(&templatePath, "template", blank, "path to the the directory or file to process")
	flag.StringVar(&variablePath, "varpath", blank, "path to the variable config to use while processing")
	flag.StringVar(&outputPath, "output", "./", "Output path for the rendered templates to be outputted")
	flag.BoolVar(&validate, "validate", false, "validate syntax and check for the required variables")
	flag.BoolVar(&debug, "verbose", false, "enable verbose logging")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	var (
		vortex = processor.New()
	)
	switch {
	case variablePath != "" && templatePath != "":
		if validate {
			vortex.EnableStrict()
		}
	default:
		flag.Usage()
		return
	}
	if debug {
		vortex.EnableDebug()
	}
	if err := vortex.LoadVariables(variablePath); err != nil {
		log.Warn("Unable to load variables due to %v", err)
		os.Exit(1)
	}
	if err := vortex.ProcessTemplates(templatePath, outputPath); err != nil {
		log.Warn("Unable to process templates due to %v", err)
		os.Exit(1)
	}
}
