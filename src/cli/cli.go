package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/luislve17/amauta/linter"
	"github.com/luislve17/amauta/ui"
)

func RunCLI() {
	inputPath := flag.String("i", "", "input path (currently unused)")
	outputPath := flag.String("o", "./dist/doc.html", "output HTML file")
	lintOnly := flag.Bool("lint", false, "only run linter (alias)")
	render := flag.Bool("render", false, "render HTML (alias)")
	flag.Parse()

	runLint := *lintOnly
	runRender := *render

	if !runLint && !runRender {
		fmt.Fprintf(os.Stderr, "%sError:%s %sYou must specify either --lint or --render%s\n", red, reset, bold, reset)
		flag.Usage()
		os.Exit(1)
	}

	if *inputPath == "" {
		fmt.Fprintf(os.Stderr, "%sError:%s %sYou must specify a non-empty input file path%s\n", red, reset, bold, reset)
		flag.Usage()
		os.Exit(1)
	}

	docManifestContent := loadFromFile(*inputPath)
	data, err := linter.LintFromRoot(docManifestContent, true)
	if err != nil {
		log.Fatalf("Linting failed: %v", err)
	}

	if runLint && !runRender {
		fmt.Println("Linting successful.")
		return
	}

	err = ui.RenderToFile(*outputPath, data)
	if err != nil {
		log.Fatalf("failed to render HTML: %v", err)
	}

	fmt.Printf("HTML rendered to %s\n", *outputPath)
}
