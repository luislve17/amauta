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
	cli := setupCli()
	verifyCliUsage(cli)

	resolvedDocManifestContent := resolveRefsInContent(*cli.inputPath)
	data, err := linter.LintFromRoot(resolvedDocManifestContent, true)
	if err != nil {
		log.Fatalf("Linting failed: %v", err)
	}
	if data.Status != linter.LintStatusOK {
		log.Fatalf("Linting failed: %v", data.Msg)
	}

	if *cli.lintOnly && !*cli.render {
		fmt.Println("Linting successful.")
		return
	}

	root := data.Structure.Root.Info.(*linter.Root)
	root.ThemeStyle, err = loadThemeFromFile(fmt.Sprintf("./ui/themes/%s.css", *cli.theme))
	if err != nil {
		fmt.Printf("Theme not found: %s. %s. Using default\n", *cli.theme, err)
	}
	err = ui.RenderToFile(*cli.outputPath, data.Structure.Root)
	if err != nil {
		log.Fatalf("failed to render HTML: %v", err)
	}

	fmt.Printf("HTML rendered to %s\n", *cli.outputPath)
}

func setupCli() cliOptions {
	version := flag.Bool("v", false, "Build version")
	inputPath := flag.String("i", "", "Input path")
	outputPath := flag.String("o", "./dist/doc.html", "Output HTML file path (defaults to './dist/doc.html')")
	lintOnly := flag.Bool("lint", false, "Lint doc manifest")
	render := flag.Bool("render", false, "Render HTML from doc manifest")
	theme := flag.String("theme", "default", "Name of the selected theme (available: 'default', 'dark')")

	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("Usage of Amauta CLI (%s):\n", buildVersion)
		fmt.Printf("-%s\tShow this help\n", styledString("h").style(italic))
		order := []string{"v", "i", "o", "lint", "render", "theme"}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			fmt.Printf("-%s\t%s\n", styledString(flag.Name).style(italic), flag.Usage)
		}
	}
	flag.Parse()

	return cliOptions{
		version:    version,
		inputPath:  inputPath,
		outputPath: outputPath,
		lintOnly:   lintOnly,
		render:     render,
		theme:      theme,
	}
}

func verifyCliUsage(cli cliOptions) {
	runLint := *cli.lintOnly
	runRender := *cli.render

	if *cli.version {
		fmt.Fprintf(
			os.Stdout,
			"%s version %s\n",
			styledString("Amauta:").style(bold), styledString(buildVersion).style(bold, green),
		)
		os.Exit(0)
	}

	if !runLint && !runRender {
		fmt.Fprintf(
			os.Stdout,
			"%s You must specify either %s or %s\n",
			styledString("Error:").style(red, bold),
			styledString("--lint").style(bold, italic),
			styledString("--render").style(bold, italic),
		)
		flag.Usage()
		os.Exit(1)
	}

	if *cli.inputPath == "" {
		fmt.Fprintf(
			os.Stderr,
			"%s You must specify a non-empty input file path (%s)\n",
			styledString("Error:").style(bold, red),
			styledString("-i").style(bold, italic),
		)
		flag.Usage()
		os.Exit(1)
	}
}
