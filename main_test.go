package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	. "github.com/AlexsJones/vortex"
	"io/ioutil"
	"os"
)

var _ = Describe("Running vortex with invalid parameters", func() {
	templateFile := "test_files/test1.yaml"
	invalidFile := "not-a-file"
	varsFile := "test_files/vars.yaml"
	testOutputDir := "this-dir-doesnt-exist"

	Context("Without a template file", func() {
		It("Should throw an error", func() {
			err := InputParametersCheck(&invalidFile, &testOutputDir, &varsFile)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Without a variable file", func() {
		It("Should throw an error", func() {
			err := InputParametersCheck(&templateFile, &testOutputDir, &invalidFile)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Running vortex with valid parameters", func() {
	Context("With an output directory that doesn't exist", func() {
		It("Should create the output directory", func() {
			templateFile := "test_files/test1.yaml"
			varsFile := "test_files/vars.yaml"
			testOutputDir := "this-dir-doesnt-exist"
			err := InputParametersCheck(&templateFile, &testOutputDir, &varsFile)
			if err != nil {
				Fail(err.Error())
			}

			Expect(testOutputDir).To(BeADirectory())

			os.RemoveAll(testOutputDir)
		})
	})

	Context("With a template file, a variable file and an output directory", func() {
		It("It should output the rendered template in the output directory", func() {
			testFiles := "test_files"
			templateFile := "test1.yaml"
			varFile := "vars.yaml"
			testOutput := "test_output"
			os.MkdirAll(testOutput, 0700)

			ParseSingleTemplate(fmt.Sprint(testFiles, "/", templateFile), testOutput, fmt.Sprint(testFiles, "/", varFile))

			content, err := ioutil.ReadFile(fmt.Sprint(testOutput, "/", templateFile))
			if err != nil {
				Fail(err.Error())
			}

			expectedSubstring := "name: test-name"
			Expect(string(content)).To(ContainSubstring(expectedSubstring))

			os.RemoveAll(testOutput)
		})
	})

	Context("With a template directory, a variable file and an output directory", func() {
		It("It should output the rendered templates in the output directory", func() {
			testFiles := "test_files"
			templateFiles := []string{"test1.yaml", "test2.yaml"}

			varFile := "vars.yaml"
			testOutput := "test_output"
			os.MkdirAll(testOutput, 0700)

			ParseDirectoryTemplates(testFiles, testOutput, fmt.Sprint(testFiles, "/", varFile))

			content1, err := ioutil.ReadFile(fmt.Sprint(testOutput, "/", templateFiles[0]))
			if err != nil {
				Fail(err.Error())
			}

			expectedSubstring1 := "name: test-name"
			Expect(string(content1)).To(ContainSubstring(expectedSubstring1))

			content2, err := ioutil.ReadFile(fmt.Sprint(testOutput, "/", templateFiles[1]))
			if err != nil {
				Fail(err.Error())
			}

			expectedSubstring2 := "image: image-test"
			Expect(string(content2)).To(ContainSubstring(expectedSubstring2))

			os.RemoveAll(testOutput)
		})
	})
})
