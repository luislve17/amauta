package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingModuleSection(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(1, 1)

	var manifestWithValidModule ManifestContent = ManifestContent(manifestWithValidModule)
	result, err := LintFromRoot(manifestWithValidModule, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(2, len(result.Structure.Root.Links))

	// Module
	expectedModuleData := []map[string]interface{}{
		{"id": "Users", "description": "Endpoints related to user operations"},
		{"id": "Items", "description": "Endpoints related to items owned by users"},
	}

	for idx := 0; idx < len(result.Structure.Root.Links); idx++ {
		moduleNode := result.Structure.Root.Links[idx]
		assert.Equal("Module", moduleNode.Info["type"])
		assert.Equal(expectedModuleData[idx]["id"], moduleNode.Info["id"])
		assert.Equal(expectedModuleData[idx]["description"], moduleNode.Info["description"])
	}
}
