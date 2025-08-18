package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingApiSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithApiContent ManifestContent = ManifestContent(manifestWithApiContent)
	result, err := LintFromRoot(manifestWithApiContent, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(2, len(result.Structure.Root.Links))

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info.(Group).BlockType)
	assert.Equal("example", groupNode.Info.(Group).Id)
	assert.Equal(1, len(groupNode.Links)) // Root

	// Tags
	tagsNode := result.Structure.Root.Links[1]
	assert.Equal("Tag", tagsNode.Info.(Tag).BlockType)
	assert.Equal("internal", tagsNode.Info.(Tag).Id)
	assert.Equal(1, len(tagsNode.Links))
}
