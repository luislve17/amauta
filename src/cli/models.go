package cli

type cliOptions struct {
	version    *bool
	inputPath  *string
	outputPath *string
	lintOnly   *bool
	render     *bool
	theme      *string
}

type regexLookupResult struct {
	Result     string
	LineNumber int
	FilePath   string
}
