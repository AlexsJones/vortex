package processor_test

import (
	"os"
	"testing"

	"github.com/AlexsJones/vortex/processor"
)

func TestLoadingVars(t *testing.T) {
	vort := processor.New()
	if err := vort.LoadVariables("../test_files/vars.yaml"); err != nil {
		t.Fatal("Failed to load vars", err)
	}
	if err := vort.LoadVariables("wombat"); err == nil {
		t.Fatal("The file wombat does not exist", err)
	}
}

func TestProcessingTemplates(t *testing.T) {
	vort := processor.New()
	if err := vort.LoadVariables("../test_files/vars.yaml"); err != nil {
		t.Fatal("Failed to load variables", err)
	}
	defer os.RemoveAll("output/")
}
