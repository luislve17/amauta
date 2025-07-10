package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingModuleSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidModule ManifestContent = ManifestContent(manifestWithValidModule)
	result, err := LintFromRoot(manifestWithValidModule, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(1, len(result.Structure.Root.Links))

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info["type"])
	assert.Equal(3, len(groupNode.Links))

	// Module
	expectedModuleIds := []string{"Users", "Items"}
	foundModuleIds := []string{}

	for idx := 0; idx < len(groupNode.Links); idx++ {
		node := groupNode.Links[idx]
		if node.Info["type"] == "root" {
			continue
		}

		assert.Equal("Module", node.Info["type"])
		foundModuleIds = append(foundModuleIds, node.Info["id"].(string))
	}

	assert.ElementsMatch(expectedModuleIds, foundModuleIds)
}

func TestRunsLinterLinkingModuleToTags(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidTaggedModules ManifestContent = ManifestContent(manifestWithValidTaggedModules)
	result, err := LintFromRoot(manifestWithValidTaggedModules, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(5, len(result.Structure.Root.Links)) // 4 tags & 1 group

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info["type"])
	assert.Equal(3, len(groupNode.Links)) // 3 modules & root

	// Modules
	expectedModuleData := []map[string]interface{}{
		{"id": "Users", "_tagIds": []string{"public", "under-dev"}},
		{"id": "Items", "_tagIds": []string{"internal"}},
	}

	for _, expectedModule := range expectedModuleData {
		for _, sectionNode := range groupNode.Links {
			if sectionNode.Info["type"] == "Module" && sectionNode.Info["id"] == expectedModule["id"] {
				assert.Equal(expectedModule["_tagIds"], sectionNode.Info["_tagIds"])
				// Modules -> Tags link
				var linkedIDs []string
				for _, linkedNode := range sectionNode.Links {
					if sectionNode.Info["type"] == "Tag" {
						linkedIDs = append(linkedIDs, linkedNode.Info["id"].(string))
						assert.ElementsMatch(sectionNode.Info["_tagIds"], linkedIDs)
					}
				}
			}
		}
	}
}

func TestRunsLinterSkippingLinkingModuleToUnexistentTags(t *testing.T) {
	assert := assert.New(t)

	var manifestWithUnexistentTaggedModules ManifestContent = ManifestContent(manifestWithUnexistentTaggedModules)
	result, err := LintFromRoot(manifestWithUnexistentTaggedModules, true)

	assert.Nil(err)
	assert.Equal(LintStatusOK, result.Status)
	assert.Equal("success", result.Msg)

	// Root
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(3, len(result.Structure.Root.Links)) // 2 tags & 1 group

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info["type"])
	assert.Equal(2, len(groupNode.Links)) // 1 module & root

	// Modules
	expectedModuleData := map[string]interface{}{"id": "Users", "_tagIds": []string{"public", "under-dev"}}

	for _, sectionNode := range result.Structure.Root.Links {
		if sectionNode.Info["type"] == "Module" && sectionNode.Info["id"] == expectedModuleData["id"] {
			assert.Equal(expectedModuleData["_tagIds"], sectionNode.Info["_tagIds"])
			assert.Equal("public", sectionNode.Links[1].Info["id"])
		}
	}
}
