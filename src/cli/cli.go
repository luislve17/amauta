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
	version := flag.Bool("v", false, "binary version")
	inputPath := flag.String("i", "", "input path (currently unused)")
	outputPath := flag.String("o", "./dist/doc.html", "output HTML file")
	lintOnly := flag.Bool("lint", false, "only run linter (alias)")
	render := flag.Bool("render", false, "render HTML (alias)")
	theme := flag.String("theme", "default", "select theme set as css")
	flag.Parse()

	runLint := *lintOnly
	runRender := *render

	if *version {
		fmt.Fprintf(os.Stdout, "%sAmauta:%s version %s%s%s%s\n", bold, reset, bold, green, buildVersion, reset)
		os.Exit(0)
	}

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

	docManifestContent := loadManifestFromFile(*inputPath)
	data, err := linter.LintFromRoot(docManifestContent, true)
	if err != nil {
		log.Fatalf("Linting failed: %v", err)
	}
	if data.Status != linter.LintStatusOK {
		log.Fatalf("Linting failed: %v", data.Msg)
	}

	if runLint && !runRender {
		fmt.Println("Linting successful.")
		return
	}

	data.Structure.Root.Info["themeStyle"], err = loadThemeFromFile(fmt.Sprintf("./ui/themes/%s.css", *theme))
	if err != nil {
		fmt.Printf("Theme not found: %s. %s. Using default\n", *theme, err)
	}
	err = ui.RenderToFile(*outputPath, data.Structure.Root)
	if err != nil {
		log.Fatalf("failed to render HTML: %v", err)
	}

	fmt.Printf("HTML rendered to %s\n", *outputPath)
}
