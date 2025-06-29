package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterIgnoringInlineComment(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(ValidManifestWithInlineComments)
	result, err := LintFromRoot(manifest, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(6, len(result.Structure.Root.Links))
}

func TestRunsLinterIgnoringMultilineComment(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(ValidManifestWithMultilineComments)
	result, err := LintFromRoot(manifest, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(6, len(result.Structure.Root.Links))
}
