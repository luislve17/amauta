package linter

import (
	"github.com/stretchr/testify/assert"

	"strings"
	"testing"
)

func TestRunsLinterFindingContentSection(t *testing.T) {
	assert := assert.New(t)

	var manifestWithValidTags ManifestContent = ManifestContent(ValidManifestWithContentSection)
	result, err := LintFromRoot(manifestWithValidTags, true)

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
	assert.Equal(2, len(groupNode.Links))

	// Content
	htmlContentStart := "<h1>Amauta</h1>"
	expectedContentData := map[string]interface{}{
		"Id":        "About amauta",
		"BlockType": "Content",
		"Summary":   htmlContentStart,
	}

	contentNode := groupNode.Links[1]
	htmlContentFromNode := string(contentNode.Info.(Content).Summary)
	htmlContentFromNode = strings.ReplaceAll(htmlContentFromNode, "\t", "")
	htmlContentFromNode = strings.Split(htmlContentFromNode, "\n")[0]
	assert.Equal("Content", contentNode.Info.(Content).BlockType)
	assert.Equal(expectedContentData["Id"], contentNode.Info.(Content).Id)
	assert.Equal(expectedContentData["BlockType"], contentNode.Info.(Content).BlockType)
	assert.Equal(htmlContentStart, htmlContentFromNode)
	assert.Equal(1, len(contentNode.Links))
	assert.Equal("Group", contentNode.Links[0].Info.(Group).BlockType)
}
