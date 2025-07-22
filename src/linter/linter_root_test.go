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

func TestRunsLinterLoadingRootData(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(manifestWithRootDetails)
	result, err := LintFromRoot(manifest, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	root := result.Structure.Root
	assert.Equal("Root", root.Info.(*Root).BlockType)
	assert.Equal("Root", root.Info.(*Root).Id)
	assert.Equal(0, len(root.Links))
	assert.Equal("https://raw.githubusercontent.com/luislve17/amauta/refs/heads/main/assets/amauta-logo.svg", root.Info.(*Root).LogoUrl)
	assert.Equal("https://github.com/luislve17/amauta", root.Info.(*Root).GithubUrl)
}
