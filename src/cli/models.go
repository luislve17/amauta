package cli

type cliOptions struct {
	version    *bool
	inputPath  *string
	outputPath *string
	lintOnly   *bool
	render     *bool
	theme      *string
}
