package ui

import (
	"testing"

	"github.com/luislve17/amauta/linter"
	"github.com/stretchr/testify/assert"
)

func TestRendersContentGraph(t *testing.T) {
	validManifest := linter.ManifestContent(linter.ValidManifest)
	data, err := linter.LintFromRoot(validManifest, true)

	_, err = Render(data)
	assert.NoError(t, err)
}
