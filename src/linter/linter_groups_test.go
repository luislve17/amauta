package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingGroupSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidGroups ManifestContent = ManifestContent(manifestWithValidGroup)
	result, err := LintFromRoot(manifestWithValidGroups, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(2, len(result.Structure.Root.Links))

	// Groups
	expectedGroupData := []map[string]interface{}{
		{"id": "getting-started", "description": "getting started"},
		{"id": "api", "description": "client api"},
	}
	for idx := 0; idx < len(result.Structure.Root.Links); idx++ {
		groupNode := result.Structure.Root.Links[idx]
		assert.Equal("Group", groupNode.Info.(Group).BlockType)
		assert.Equal(expectedGroupData[idx]["id"], groupNode.Info.(Group).Id)
		assert.Equal(expectedGroupData[idx]["description"], groupNode.Info.(Group).Description)
		assert.Equal(1, len(groupNode.Links))
		assert.Equal("Root", groupNode.Links[0].Info.(*Root).Id)
	}
}
