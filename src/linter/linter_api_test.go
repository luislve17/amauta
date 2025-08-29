package linter

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRunsLinterFindingAPISection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidAPI ManifestContent = ManifestContent(manifestWithValidAPI)
	result, err := LintFromRoot(manifestWithValidAPI, true)

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

	// API
	expectedAPIIds := []string{"Users", "Items"}
	foundAPIIds := []string{}

	for idx := 1; idx < len(groupNode.Links); idx++ {
		node := groupNode.Links[idx]
		switch info := node.Info.(type) {
		case *Root:
			continue // skip root
		case API:
			assert.Equal("API", info.BlockType)
			foundAPIIds = append(foundAPIIds, info.Id)
		default:
			t.Fatalf("unexpected node type: %T", info)
		}
	}

	assert.ElementsMatch(expectedAPIIds, foundAPIIds)
}

func TestRunsLinterLinkingAPIToTags(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidTaggedAPIs ManifestContent = ManifestContent(manifestWithValidTaggedAPIs)
	result, err := LintFromRoot(manifestWithValidTaggedAPIs, true)

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
	assert.Equal(3, len(groupNode.Links)) // 3 apis & root

	// APIs
	expectedAPIData := []map[string]interface{}{
		{"id": "Users", "_tagIds": []string{"public", "under-dev"}},
		{"id": "Items", "_tagIds": []string{"internal"}},
	}

	for _, expectedAPI := range expectedAPIData {
		for _, sectionNode := range groupNode.Links {
			switch info := sectionNode.Info.(type) {
			case *Group:
				if info.Id == expectedAPI["id"] {
					assert.Equal(expectedAPI["_tagIds"], info._tagIds)

					// APIs -> Tags link
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

func TestRunsLinterSkippingLinkingAPIToUnexistentTags(t *testing.T) {
	assert := assert.New(t)

	var manifestWithUnexistentTaggedAPIs ManifestContent = ManifestContent(manifestWithUnexistentTaggedAPIs)
	result, err := LintFromRoot(manifestWithUnexistentTaggedAPIs, true)

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
	assert.Equal(2, len(groupNode.Links)) // 1 api & root

	// APIs
	expectedAPIData := map[string]interface{}{"id": "Users", "_tagIds": []string{"public", "under-dev"}}

	for _, sectionNode := range result.Structure.Root.Links {
		switch info := sectionNode.Info.(type) {
		case *Group:
			if info.Id == expectedAPIData["id"] {
				assert.Equal(expectedAPIData["_tagIds"], sectionNode.Info.(Group)._tagIds)
				assert.Equal("public", sectionNode.Links[1].Info.(Group).Id)
			}
		}
	}
}

func TestRunsLinterCollectingApiSubsectionsInGraph(t *testing.T) {
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
	assert.Equal(2, len(groupNode.Links)) // Root + API

	// API
	APINode := groupNode.Links[1]
	assert.Equal("API", APINode.Info.(API).BlockType)
	assert.Equal("Endpoint", APINode.Info.(API).Id)
	assert.Equal(2, len(APINode.Links)) // Group + Endpoint

	// Endpoint
	EndpointNode := APINode.Links[0]
	assert.Equal("/v1/products", EndpointNode.Info.(Endpoint).Id)
	assert.Equal("Endpoint", EndpointNode.Info.(Endpoint).BlockType)
	assert.Equal(2, len(EndpointNode.Links)) // Htpp + API

	// HTTP Verb
	HttpVerbNode := EndpointNode.Links[0]
	assert.Equal("GET", HttpVerbNode.Info.(HTTPVerb).Id)
	assert.Equal("Http", HttpVerbNode.Info.(HTTPVerb).BlockType)
}
