package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/AlexsJones/vortex/processor"
)

var _ = Describe("Vortex loading variables", func() {
	var (
		vort          = processor.New()
		flatstructure = "test_files/vars.yaml"
		invalid_file  = "test_files/badvars.yaml"
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
		variablePath = "test_files/vars.yaml"
		templatefile = "test_files/test1.yaml"
		badTemplate  = "test_files/badtemplate.yaml"
	)
	Context("With a valid template and variables defined", func() {
		It("should not error", func() {
			vort := processor.New()
			if err := vort.LoadVariables(variablePath); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
			err := vort.EnableStrict().ProcessTemplates(templatefile, "./")
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("With a valid template with a invalid yaml definitions", func() {
		It("it should report an error", func() {
			vort := processor.New()
			if err := vort.LoadVariables(variablePath); err != nil {
				Expect(err).NotTo(HaveOccurred())
			}
			err := vort.EnableStrict().ProcessTemplates(badTemplate, "")
			Expect(err).To(HaveOccurred())
		})
	})
})
