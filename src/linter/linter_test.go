package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterWithoutCreatingStructure(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(validManifest)
	result, err := LintFromRoot(manifest, false)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)
	assert.Nil(result.Structure)
}

func TestRunsLinterCreatingStructureGraph(t *testing.T) {
	assert := assert.New(t)

	var manifest ManifestContent = ManifestContent(validManifest)
	result, err := LintFromRoot(manifest, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)
}

func TestRunsLinterFindingTagSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidTags ManifestContent = ManifestContent(manifestWithValidTags)
	result, err := LintFromRoot(manifestWithValidTags, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(4, len(result.Structure.Root.Links))

	// Tags
	expectedTagData := []map[string]interface{}{
		{"id": "public", "color": "#00FF00", "description": "Public API"},
		{"id": "internal", "color": "#AAAAAA", "description": "Internal use only"},
		{"id": "deprecated", "color": "#FF6F61", "description": "Will be removed soon"},
		{"id": "under-dev", "color": "#FFD966", "description": "Still under development"},
	}
	for idx := 0; idx < len(result.Structure.Root.Links); idx++ {
		tagNode := result.Structure.Root.Links[idx]
		assert.Equal("tag", tagNode.Info["type"])
		assert.Equal(expectedTagData[idx]["id"], tagNode.Info["id"])
		assert.Equal(expectedTagData[idx]["color"], tagNode.Info["color"])
		assert.Equal(expectedTagData[idx]["description"], tagNode.Info["description"])
		assert.Equal(1, len(tagNode.Links))
		assert.Equal("root", tagNode.Links[0].Info["type"])
	}
}

func TestFailsLinterWhenTagsFailExpectedFormat(t *testing.T) {
	assert := assert.New(t)

	var manifestWithInvalidTags ManifestContent = ManifestContent(manifestWithInvalidTags)
	result, err := LintFromRoot(manifestWithInvalidTags, true)

	expectedErrMsg := "Error@line:8\n->Invalid tag format: \"internal@AAAAAA: Invalid tag format\""
	assert.Equal(expectedErrMsg, err.Error())
	assert.Equal(result.Status, LintStatusError)
	assert.Equal("error", result.Msg)

}
