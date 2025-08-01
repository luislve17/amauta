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

func TestMarkdownRendererInception(t *testing.T) {
	assert := assert.New(t)
	var loadedManifestWithComplexMdContent = ManifestContent(manifestWithComplexMdContent)
	resultGraph, err := generateGraph(loadedManifestWithComplexMdContent)

	assert.Nil(err)
	contentNode := resultGraph.Root.Links[0].Links[1]
	contentSummary := contentNode.Info.(Content).Summary

	expectedHTML := `<h1>Content</h1>
<pre><code class="language-toml">[[&lt;name&gt;@content#&lt;tag_ids&gt;]]
group: &lt;group-id&gt;
summary: &lt;summary&gt;

</code></pre>
<p>This is an inline code <code>&lt;md&gt;</code> and <code>&lt;/md&gt;</code></p>
<h3>Example</h3>
<pre><code class="language-toml">[[Tax management@content#internal]]
group: finances
summary: &lt;md&gt;
# Support for markdown!
## Subtitle
List:
1. One
2. Two
3. Three

&gt; This is a quote

| tables | also | work |
|--------|------|------|
| tables | also | work |
&lt;/md&gt;

</code></pre>
`
	assert.Equal(expectedHTML, string(contentSummary))
}
