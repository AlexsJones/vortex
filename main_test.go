package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/AlexsJones/vortex"
	"io/ioutil"
	"os"
	"fmt"
)

var _ = Describe("Running the main application with valid parameters", func() {
	Context("When running with an output directory that doesn't exist", func() {
		It("Should create the output directory", func() {
			// Given I have a directory name that doesn't exist
			// When I try to run vortex
			// Then vortex will try to create the output directory for me
			testOutputDir := "this-dir-doesnt-exist"
			err := CreateOutputDirectoryIfDoesntExist(testOutputDir)
			if err != nil {
				Fail(err.Error())
			}

			Expect(testOutputDir).To(BeADirectory())

			// Cleanup
			os.RemoveAll(testOutputDir)
		})
	})

	Context("When provide a template file, a variable file and an output directory", func() {
		It("It should output the rendered template in the output directory", func() {
			// Setup
			testFiles := "test_files"
			templateFile := "test1.yaml"
			envFile := "env1.yaml"
			testOutput := "test_output"
			os.MkdirAll(testOutput, 0700)

			// Given + When

			ParseSingleTemplate(fmt.Sprint(testFiles, "/", templateFile), testOutput, fmt.Sprint(testFiles, "/", envFile))
			// Then
			content, err := ioutil.ReadFile(fmt.Sprint(testOutput, "/", templateFile))
			if err != nil {
				Fail(err.Error())
			}

			expectedSubstring := "name: test-name"
			Expect(string(content)).To(ContainSubstring(expectedSubstring))

			// Cleanup
			os.RemoveAll(testOutput)
		})
	})
})
