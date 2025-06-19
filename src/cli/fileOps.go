package cli

import (
	"log"
	"os"

	"github.com/luislve17/amauta/linter"
)

func loadFromFile(path string) linter.ManifestContent {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	return linter.ManifestContent(content)
}
