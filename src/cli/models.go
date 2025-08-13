package cli

import "github.com/luislve17/amauta/linter"

type cliOptions struct {
	version    *bool
	inputPath  *string
	outputPath *string
	lintOnly   *bool
	render     *bool
	theme      *string
}

type regexLookupResult struct {
	Result    string
	FilePath  string
	LineRange linter.LineRange
}
