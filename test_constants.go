package main

const (
	TemplateFile = "template.yaml"
	TemplateDir  = "templates"
	VarsFile     = "vars.yaml"

	Template      = "template"
	Directory     = "directory"
	DirectoryTree = "directory tree"
	Valid         = "valid"
	Invalid       = "invalid"
	Mixed         = "mixed"
	Contains      = "contains"
	Successful    = "successful"
)

var (
	directoryValidTemplatePaths   = []string{"templates/template1.yaml", "templates/template2.yaml"}
	directoryInvalidTemplatePaths = []string{"templates/template3.yaml"}

	directoryTreePaths                = []string{"templates/one", "templates/two"}
	directoryTreeValidTemplatePaths   = []string{"templates/one/template1.yaml", "templates/two/template2.yaml"}
	directoryTreeInvalidTemplatePaths = []string{"templates/template3.yaml"}

	directoryTreeOutputFilePaths = []string{"deployment/one/template1.yaml", "deployment/two/template2.yaml"}
)
