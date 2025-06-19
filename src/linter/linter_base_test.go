package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterWithoutCreatingStructure(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(ValidManifest)
	result, err := LintFromRoot(manifest, false)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)
	assert.Nil(result.Structure)
}

func TestRunsLinterCreatingStructureGraph(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(ValidManifest)
	result, err := LintFromRoot(manifest, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)
}
