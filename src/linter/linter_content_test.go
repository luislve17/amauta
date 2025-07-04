package linter

import (
	"github.com/stretchr/testify/assert"
	"html/template"

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
	assert.Equal("root", result.Structure.Root.Info["type"])
	assert.Equal("root", result.Structure.Root.Info["id"])
	assert.Equal(1, len(result.Structure.Root.Links))

	// Group
	groupNode := result.Structure.Root.Links[0]
	assert.Equal("Group", groupNode.Info["type"])
	assert.Equal(2, len(groupNode.Links))

	// Content
	htmlContentStart := "<h1>Amauta</h1>"
	expectedContentData := map[string]interface{}{
		"id":       "About amauta",
		"type":     "Content",
		"htmlBody": htmlContentStart,
	}

	contentNode := groupNode.Links[1]
	htmlContentFromNode := string(contentNode.Info["htmlContent"].(template.HTML))
	htmlContentFromNode = strings.ReplaceAll(htmlContentFromNode, "\t", "")
	htmlContentFromNode = strings.Split(htmlContentFromNode, "\n")[0]
	assert.Equal("Content", contentNode.Info["type"])
	assert.Equal(expectedContentData["id"], contentNode.Info["id"])
	assert.Equal(expectedContentData["type"], contentNode.Info["type"])
	assert.Equal(htmlContentStart, htmlContentFromNode)
	assert.Equal(1, len(contentNode.Links))
	assert.Equal("Group", contentNode.Links[0].Info["type"])
}
