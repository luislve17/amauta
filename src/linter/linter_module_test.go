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
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(1, len(result.Structure.Root.Links))

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info.(Group).BlockType)
	assert.Equal(3, len(groupNode.Links))

	// Module
	expectedModuleIds := []string{"Users", "Items"}
	foundModuleIds := []string{}

	for idx := 1; idx < len(groupNode.Links); idx++ {
		node := groupNode.Links[idx]
		switch info := node.Info.(type) {
		case *Root:
			continue // skip root
		case Module:
			assert.Equal("Module", info.BlockType)
			foundModuleIds = append(foundModuleIds, info.Id)
		default:
			t.Fatalf("unexpected node type: %T", info)
		}
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
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(5, len(result.Structure.Root.Links)) // 4 tags & 1 group

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info.(Group).BlockType)
	assert.Equal(3, len(groupNode.Links)) // 3 modules & root

	// Modules
	expectedModuleData := []map[string]interface{}{
		{"id": "Users", "_tagIds": []string{"public", "under-dev"}},
		{"id": "Items", "_tagIds": []string{"internal"}},
	}

	for _, expectedModule := range expectedModuleData {
		for _, sectionNode := range groupNode.Links {
			switch info := sectionNode.Info.(type) {
			case *Group:
				if info.Id == expectedModule["id"] {
					assert.Equal(expectedModule["_tagIds"], info._tagIds)

					// Modules -> Tags link
					var linkedIDs []string
					for _, linkedNode := range sectionNode.Links {
						if tag, ok := linkedNode.Info.(*Tag); ok {
							linkedIDs = append(linkedIDs, tag.Id)
						}
					}
					assert.ElementsMatch(info._tagIds, linkedIDs)
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
	assert.Equal("Root", result.Structure.Root.Info.(*Root).BlockType)
	assert.Equal("Root", result.Structure.Root.Info.(*Root).Id)
	assert.Equal(3, len(result.Structure.Root.Links)) // 2 tags & 1 group

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info.(Group).BlockType)
	assert.Equal(2, len(groupNode.Links)) // 1 module & root

	// Modules
	expectedModuleData := map[string]interface{}{"id": "Users", "_tagIds": []string{"public", "under-dev"}}

	for _, sectionNode := range result.Structure.Root.Links {
		switch info := sectionNode.Info.(type) {
		case *Group:
			if info.Id == expectedModuleData["id"] {
				assert.Equal(expectedModuleData["_tagIds"], sectionNode.Info.(Group)._tagIds)
				assert.Equal("public", sectionNode.Links[1].Info.(Group).Id)
			}
		}
	}
}
