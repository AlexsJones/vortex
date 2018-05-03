package abstraction

// TemplateProcessor allows to imply the use of
// the adapter pattern inside the vortex project
// to allow for cleaner and more concise code
type TemplateProcessor interface {
	// LoadVariables will read the file define by the variablePath
	// and store them inside the Processor
	// An issues doing so will return an error
	LoadVariables(variablePath string) error

	// ProcessTemplates will ready all files that are found inside
	// templatePath (and all subsequent directories) and output them
	// into the relative output path.
	ProcessTemplates(templatePath, outputPath string) error

	// EnableStrict will enforce that templates have all the required vars
	EnableStrict()
}
