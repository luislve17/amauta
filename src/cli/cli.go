package cli

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/luislve17/amauta/linter"
	"github.com/luislve17/amauta/ui"
)

func Run() {
	// CLI flags
	inputPath := flag.String("i", ".", "input path (currently unused)")
	outputPath := flag.String("o", "./dist/doc.html", "output HTML file")
	lintOnly := flag.Bool("l", false, "only run linter")
	lintOnlyAlias := flag.Bool("lint", false, "only run linter (alias)")
	renderHTML := flag.Bool("r", false, "render HTML")
	renderAlias := flag.Bool("render", false, "render HTML (alias)")
	flag.Parse()

	runLint := *lintOnly || *lintOnlyAlias
	runRender := *renderHTML || *renderAlias

	if !runLint && !runRender {
		fmt.Fprintln(os.Stderr, "Error: You must specify either --lint or --render")
		flag.Usage()
		os.Exit(1)
	}

	docManifestContent := loadFromFile(*inputPath)
	data, err := linter.LintFromRoot(docManifestContent, true)
	if err != nil {
		log.Fatalf("linting failed: %v", err)
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
