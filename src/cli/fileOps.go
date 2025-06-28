package cli

import (
	"html/template"
	"log"
	"os"

	"github.com/luislve17/amauta/linter"
)

func loadManifestFromFile(path string) linter.ManifestContent {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	return linter.ManifestContent(content)
}

func loadThemeFromFile(path string) (template.CSS, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return template.CSS(""), err
	}
	return template.CSS(content), nil
}
