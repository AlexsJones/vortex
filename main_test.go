package main_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/AlexsJones/vortex/processor"
)

func TestVortexFunctionality(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vortex Suite")
}

var _ = Describe("Vortex loading variables", func() {
	var (
		vort          = processor.New()
		flatstructure = "example/vars.vortex"
		invalid_file  = "example/.mistyped/variables.yaml"
	)
	Context("With a valid template file", func() {
		It("Should not report an error", func() {
			err := vort.LoadVariables(flatstructure)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("With incorrectly formatted file", func() {
		It("Should report an error", func() {
			err := vort.LoadVariables(invalid_file)
			Expect(err).To(HaveOccurred())
		})
	})
	Context("With a file does not exist", func() {
		It("Should report an error", func() {
			err := vort.LoadVariables("wombat")
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Vortex validating templates", func() {
	var (
		variablePath = "example/vars.vortex"
		templatefile = "example/demo.yaml"
		badTemplate  = "example/.mistyped/template.yaml"
	)
	Context("With a valid template and variables defined", func() {
		It("should not error", func() {
			vort := processor.New()
			if err := vort.LoadVariables(variablePath); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
			err := vort.EnableStrict(true).ProcessTemplates(templatefile, "./")
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("With a valid template that has a invalid yaml definitions", func() {
		It("it should report an error", func() {
			vort := processor.New()
			if err := vort.LoadVariables(variablePath); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
			err := vort.EnableStrict(true).ProcessTemplates(badTemplate, "")
			Expect(err).To(HaveOccurred())
		})
	})
	Context("With a validate template but missing variables", func() {
		It("Should report an error", func() {
			vort := processor.New()
			err := vort.EnableStrict(true).ProcessTemplates(badTemplate, "")
			Expect(err).To(HaveOccurred())
		})
	})
	Context("with a valid vars file and a directory to validate", func() {
		It("should not report an error", func() {
			var (
				varPath     = "example/vars.vortex"
				templateDir = "example"
			)
			vort := processor.New()
			if err := vort.LoadVariables(varPath); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
			err := vort.EnableStrict(true).ProcessTemplates(templateDir, "")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})

var _ = Describe("Processing templates with vortex", func() {
	Context("Using a single valid template file as the template path", func() {
		It("should not report an error", func() {
			var (
				varPath     = "example/vars.vortex"
				templateDir = "example"
				vort        = processor.New()
			)
			dir, err := ioutil.TempDir("", "output")
			Expect(err).NotTo(HaveOccurred())

			// Ensure we clean up after ourselves
			defer os.RemoveAll(dir)

			err = vort.LoadVariables(varPath)
			Expect(err).NotTo(HaveOccurred())
			err = vort.ProcessTemplates(templateDir, dir)
			Expect(err).NotTo(HaveOccurred())
			_, err = os.Stat(path.Join(dir, "demo.yaml"))
			Expect(os.IsNotExist(err)).NotTo(Equal(true))
			_, err = os.Stat(path.Join(dir, "bar/example.yaml"))
			Expect(os.IsNotExist(err)).NotTo(Equal(true))
		})
	})
})
