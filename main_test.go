package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/AlexsJones/vortex"
	"io/ioutil"
	"strings"
	"os"
)

var _ = Describe("Running the main application with valid parameters", func() {
	Context("When provide a template file, a variable file and an output directory", func() {
		It("It should output the rendered template in the output directory", func() {
			// Setup
			os.MkdirAll("test_output", 0700)

			// Given + When
			ParseSingleTemplate("test_files/test1.yaml", "test_output", "test_files/env1.yaml")
			// Then
			content, err := ioutil.ReadFile("test_output/test1.txt")
			if err != nil {
				Fail(err.Error())
			}

			expectedSubstring := "name: test-name"
			Expect(strings.Contains(string(content), expectedSubstring))

			// Cleanup
			os.RemoveAll("test_output")
		})
	})
})
