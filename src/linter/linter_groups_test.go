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
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(2, len(result.Structure.Root.Links))

	// Groups
	expectedTagData := []map[string]interface{}{
		{"id": "getting-started", "description": "getting started"},
		{"id": "api", "description": "client api"},
	}
	for idx := 0; idx < len(result.Structure.Root.Links); idx++ {
		groupNode := result.Structure.Root.Links[idx]
		assert.Equal("Group", groupNode.Info["type"])
		assert.Equal(expectedTagData[idx]["id"], groupNode.Info["id"])
		assert.Equal(expectedTagData[idx]["color"], groupNode.Info["color"])
		assert.Equal(expectedTagData[idx]["description"], groupNode.Info["description"])
		assert.Equal(1, len(groupNode.Links))
		assert.Equal("root", groupNode.Links[0].Info["type"])
	}
}
