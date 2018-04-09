package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestVortex(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Vortex Suite")
}
