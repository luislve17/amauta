package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingTagSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidTags ManifestContent = ManifestContent(manifestWithValidTags)
	result, err := LintFromRoot(manifestWithValidTags, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
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
		assert.Equal("Tag", tagNode.Info.(Tag).BlockType)
		assert.Equal(expectedTagData[idx]["id"], tagNode.Info.(Tag).Id)
		assert.Equal(expectedTagData[idx]["color"], tagNode.Info.(Tag).color)
		assert.Equal(expectedTagData[idx]["description"], tagNode.Info.(Tag).Description)
		assert.Equal(1, len(tagNode.Links))
		assert.Equal("Root", tagNode.Links[0].Info.(*Root).BlockType)
	}
}

func TestRunsLinterFindingEmptyTagSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithEmptyTags ManifestContent = ManifestContent(manifestWithEmptyTags)
	result, err := LintFromRoot(manifestWithEmptyTags, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(0, len(result.Structure.Root.Links))
}

func TestFailsLinterWhenTagsFailExpectedFormat(t *testing.T) {
	assert := assert.New(t)

	var manifestWithInvalidTags ManifestContent = ManifestContent(manifestWithInvalidTags)
	result, err := LintFromRoot(manifestWithInvalidTags, true)

	expectedErrMsg := "Error@line:11\n->Invalid tag format: \"internal@AAAAAA: Invalid tag format\""
	assert.Equal(expectedErrMsg, err.Error())
	assert.Equal(result.Status, LintStatusError)
	assert.Equal("error", result.Msg)

}
