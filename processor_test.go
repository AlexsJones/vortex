package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/DATA-DOG/godog"
	"log"
)

type ProcessorFeature struct {
	template        []byte
	vars            []byte
	outputFile      string
	outputDirectory string
}

func ProcessorFeatureContext(s *godog.Suite) {
	p := ProcessorFeature{
		template: []byte(`apiVersion: v1
kind: Pod
metadata:
  name: {{.name}}
spec:
  restartPolicy: {{.restartPolicy}}
  containers:
    - name: test
      image: us.gcr.io/test
`),
		vars: []byte(`name: vortex
restartPolicy: always
`),
		outputFile:      "output.yaml",
		outputDirectory: "deployment",
	}

	s.Step(`^a template file$`, p.aTemplateFile)
	s.Step(`^a variable file$`, p.aVariableFile)
	s.Step(`^vortex is run for a (template|directory)$`, p.vortexIsRun)
	s.Step(`^an output file should contain the interpolated variables$`, p.anOutputFileShouldContainTheInterpolatedVariables)
	s.Step(`^a template directory$`, p.aTemplateDirectory)
	s.Step(`^an output directory should contain the output files$`, p.anOutputDirectoryShouldContainTheOutputFiles)
	s.Step(`^the output files should contain the interpolated variables$`, p.theOutputFilesShouldContainTheInterpolatedVariables)

	s.BeforeScenario(func(interface{}) {
		os.RemoveAll(TemplateDir)
		log.SetOutput(ioutil.Discard)
	})

	s.AfterScenario(func(interface{}, error) {
		filesToCleanup := []string{TemplateFile, VarsFile, p.outputFile}

		for _, file := range filesToCleanup {
			os.Remove(file)
		}

		os.RemoveAll(TemplateDir)
		os.RemoveAll(p.outputDirectory)
	})
}

func (p *ProcessorFeature) aTemplateFile() error {
	if err := ioutil.WriteFile(TemplateFile, p.template, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (p *ProcessorFeature) aVariableFile() error {
	if err := ioutil.WriteFile(VarsFile, p.vars, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (p *ProcessorFeature) vortexIsRun(input string) error {
	if input == "template" {
		if err := ParseSingleTemplate(TemplateFile, p.outputFile, VarsFile); err != nil {
			return err
		}
	} else {
		if err := ParseDirectoryTemplates(TemplateDir, p.outputDirectory, VarsFile); err != nil {
			return err
		}
	}

	return nil
}

func (p *ProcessorFeature) anOutputFileShouldContainTheInterpolatedVariables() error {
	output, err := ioutil.ReadFile(p.outputFile)
	if err != nil {
		return err
	}

	file, err := os.Open(VarsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !strings.Contains(string(output), scanner.Text()) {
			return fmt.Errorf("expected output file to contain: %s", scanner.Text())
		}
	}

	return nil
}

func (p *ProcessorFeature) aTemplateDirectory() error {
	for _, dir := range directoryTreePaths {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	for _, template := range directoryTreeValidTemplatePaths {
		if err := ioutil.WriteFile(template, p.template, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (p *ProcessorFeature) anOutputDirectoryShouldContainTheOutputFiles() error {
	for _, template := range directoryTreeOutputFilePaths {
		if _, err := os.Stat(template); os.IsNotExist(err) {
			return fmt.Errorf("%s was not created as an output file", template)
		}
	}

	return nil
}

func (p *ProcessorFeature) theOutputFilesShouldContainTheInterpolatedVariables() error {
	file, err := os.Open(VarsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	var expected []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		expected = append(expected, scanner.Text())
	}

	for _, outputFile := range directoryTreeOutputFilePaths {
		output, err := ioutil.ReadFile(outputFile)
		if err != nil {
			return err
		}

		for _, line := range expected {
			if !strings.Contains(string(output), line) {
				return fmt.Errorf("expected output file to contain: %s", line)
			}
		}
	}

	return nil
}
