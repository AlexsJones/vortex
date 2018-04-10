package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"io/ioutil"
	"os"

	. "github.com/AlexsJones/vortex"
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

var _ = Describe("Running the vortex validator", func() {
	Context("With non-yaml files in a directory containing yaml template files", func() {
		It("Should only try to validate yaml files", func() {
			template := []byte(`a: {{.var}}`)
			vars := []byte(`var: some-var`)
			readme := []byte(`SOME MARKDOWN HERE * * VERY NICE * *`)

			err := os.MkdirAll("tmp", 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp/README.md", readme, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp/template.yaml", template, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("vars.yaml", vars, 0700)
			Expect(err).ToNot(HaveOccurred())

			areValid, err := InputFilesAreValid("tmp", "vars.yaml")
			Expect(areValid).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			os.RemoveAll("tmp")
			os.RemoveAll("vars.yaml")
		})
	})

	Context("With nested subdirectories of containing an invalid template", func() {
		It("Should fail validation", func() {
			validTemplate := []byte(`apiVersion: v1
kind: Pod
metadata:
 name: {{.name}}
spec:
 restartPolicy: Always
 containers:
   - name: test
     image: {{.image}}
`)

			invalidTemplate := []byte(`apiVersion: v1
kind: Pod
metadata:
 name: {{.name}}
spec:
 restartPolicy: Always
 containers:
   - name: test
     image: {{.anotherimage}}
`)

			vars := []byte(`name: some-name
image: some-image
`)

			err := os.MkdirAll("tmp/one", 0700)
			Expect(err).ToNot(HaveOccurred())
			err = os.MkdirAll("tmp/two", 0700)
			Expect(err).ToNot(HaveOccurred())

			template1 := "tmp/one/template1.yaml"
			template2 := "tmp/two/template2.yaml"
			varFile := "vars.yaml"

			err = ioutil.WriteFile(template1, validTemplate, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile(template2, invalidTemplate, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile(varFile, vars, 0700)
			Expect(err).ToNot(HaveOccurred())

			areValid, err := InputFilesAreValid("tmp", varFile)
			Expect(areValid).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())

			os.RemoveAll("tmp/one")
			os.RemoveAll("tmp/two")
			os.RemoveAll("tmp")
			os.Remove("vars.yaml")
		})
	})

	Context("With a preamble comment in the yaml template", func() {
		It("Should pass validation", func() {
			template := []byte(`{{.templatepreamble}}
apiVersion: v1
kind: Pod
metadata:
  name: {{.name}}
spec:
  restartPolicy: Always
  containers:
    - name: test
      image: {{.image}}
`)

			vars := []byte(`name: some-name
image: some-image
templatepreamble: # some preamble
`)

			err := os.Mkdir("tmp", 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp/template.yaml", template, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp/vars.yaml", vars, 0700)
			Expect(err).ToNot(HaveOccurred())

			templateFile := "tmp/template.yaml"
			varsFile := "tmp/vars.yaml"

			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			os.RemoveAll("tmp")
		})
	})

	Context("With invalid syntax in a variables file", func() {
		It("Should fail validation", func() {
			templateFile := "test_files/test1.yaml"
			varsFile := "test_files/badvars.yaml"
			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With {{.var}} syntax in a template file", func() {
		It("Should pass validation", func() {
			templateFile := "test_files/test1.yaml"
			varsFile := "test_files/vars.yaml"
			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With invalid syntax in a template file", func() {
		It("Should fail validation", func() {
			templateFile := "test_files/badtemplate.yaml"
			varsFile := "test_files/vars.yaml"
			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeFalse())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("With valid template and var files with fully corresponding variables", func() {
		It("Should pass validation", func() {
			templateFile := "test_files/test2.yaml"
			varsFile := "test_files/vars.yaml"
			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With a templated variable that doesn't exist in the given var file", func() {
		It("Should fail validation", func() {
			templateFile := "test_files/test3.yaml"
			varsFile := "test_files/vars.yaml"
			areValid, err := InputFilesAreValid(templateFile, varsFile)
			Expect(areValid).To(BeFalse())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("With a directory of valid templates and a valid var file with fully corresponding variables", func() {
		It("Should pass validation", func() {
			templateDir := "tmp1"
			varFileDir := "tmp2"
			validTemplate := []byte(`apiVersion: v1
kind: Pod
metadata:
  name: {{.name}}
spec:
  restartPolicy: Always
  containers:
    - name: test
      image: {{.image}}
`)

			validVarFile := []byte(`name: some-name
image: some-image
`)
			err := os.Mkdir(templateDir, 0700)
			Expect(err).ToNot(HaveOccurred())
			err = os.Mkdir(varFileDir, 0700)
			Expect(err).ToNot(HaveOccurred())

			err = ioutil.WriteFile("tmp1/template1.yaml", validTemplate, 0644)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp1/template2.yaml", validTemplate, 0644)
			Expect(err).ToNot(HaveOccurred())
			err = ioutil.WriteFile("tmp2/vars.yaml", validVarFile, 0644)
			Expect(err).ToNot(HaveOccurred())

			areValid, err := InputFilesAreValid(templateDir, "tmp2/vars.yaml")
			Expect(areValid).To(BeTrue())
			Expect(err).ToNot(HaveOccurred())

			os.RemoveAll(varFileDir)
			os.RemoveAll(templateDir)
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

	Context("With a template directory, a variable file and an output directory", func() {
		It("It should ignore nested subdirectories in the template directory", func() {
			// Setup
			testFiles := "test_files"

			varFile := "vars.yaml"
			testOutput := "test_output"
			subdirectoryToIgnore := "sub_directory"
			os.MkdirAll(testOutput, 0700)

			// When
			ParseDirectoryTemplates(testFiles, testOutput, fmt.Sprint(testFiles, "/", varFile))

			// Then
			if _, err := os.Stat(fmt.Sprint(testOutput, "/", subdirectoryToIgnore)); os.IsNotExist(err) {
				os.RemoveAll(testOutput)
				return
			} else {
				Fail("Subdirectory was not ignored")
			}
		})
	})
})
