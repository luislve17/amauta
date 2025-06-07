package linter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func LintFromRoot(content ManifestContent, createStructure bool) (LintResult, error) {
	var contentGraphErr error = nil
	var graph *StructureGraph = nil
	lintStatus := LintStatusOK
	resultMsg := "success"
	if createStructure {
		graph, contentGraphErr = generateGraph(content)
	}
	err := errors.Join(contentGraphErr)
	if err != nil {
		lintStatus = LintStatusError
		resultMsg = "error"
	}

	return LintResult{Status: lintStatus, Msg: resultMsg, Structure: graph}, err
}

func generateGraph(content ManifestContent) (*StructureGraph, error) {
	root := &Node{
		Info: map[string]interface{}{
			"type": "root",
			"id":   "root",
		},
		Links: []*Node{},
	}

	rawBlocks := extractRawBlocks(content)
	tags, tagsErr := getTags(rawBlocks)
	if tagsErr != nil {
		return nil, tagsErr
	}
	linkNodeOneToMany(root, tags)

	// modules := getModuleNodes(rawBlocks)
	// linkNodeOneToMany(root, modules)

	graph := StructureGraph{
		Root: root,
	}

	return &graph, nil
}

func linkNodeOneToMany(mainNode *Node, nodes []*Node) {
	for _, n := range nodes {
		n.Links = append(n.Links, mainNode)
	}
	mainNode.Links = append(mainNode.Links, nodes...)
}

func findTagsSection(rawBlocks []RawBlock) *RawBlock {
	headerPattern := regexp.MustCompile(`^\[\[tags\]\]`)

	for _, section := range rawBlocks {
		lines := strings.Split(section.Content, "\n")
		if len(lines) == 0 {
			continue
		}
		header := strings.TrimSpace(lines[0])
		if headerPattern.MatchString(header) {
			return &section
		}
	}
	return nil
}

func getTags(rawBlocks []RawBlock) ([]*Node, error) {
	tagPattern := regexp.MustCompile(`^([-_\w]+)(#[A-F|\d]{6}):\s*(.*)`)

	tagSection := findTagsSection(rawBlocks)
	if tagSection == nil {
		return nil, nil
	}

	lines := strings.Split(tagSection.Content, "\n")
	var nodes []*Node

	for i, tag := range lines[1:] {
		if tag == "" {
			continue
		}
		match := tagPattern.FindStringSubmatch(tag)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", tagSection.StartLine+i+1, tag)
		}

		node := &Node{
			Info: map[string]interface{}{
				"type":        "tag",
				"id":          match[1],
				"color":       match[2],
				"description": match[3],
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getModuleNodes(rawBlocks []string) []*Node {
	headerPattern := regexp.MustCompile(`^\[\[(.+?)#(.*?)\]\]`)

	var nodes []*Node

	for _, section := range rawBlocks {
		lines := strings.Split(section, "\n")
		if len(lines) == 0 {
			continue
		}

		header := strings.TrimSpace(lines[0])
		matches := headerPattern.FindStringSubmatch(header)

		node := &Node{
			Info: map[string]interface{}{
				"type": "module",
				"name": matches[1],
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func extractRawBlocks(content ManifestContent) []RawBlock {
	lines := strings.Split(string(content), "\n")
	var blocks []RawBlock
	var current []string
	var startLine int
	inBlock := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		isNewModuleStart := strings.HasPrefix(trimmed, "[[")
		isNewInnerSectionStart := strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[")
		hasCurrentBlock := len(current) > 0

		if isNewModuleStart {
			if inBlock && hasCurrentBlock {
				blocks = append(blocks, RawBlock{Content: strings.Join(current, "\n"), StartLine: startLine})
			}
			current = []string{line}
			startLine = i + 1
			inBlock = true
		} else if inBlock && isNewInnerSectionStart {
			blocks = append(blocks, RawBlock{Content: strings.Join(current, "\n"), StartLine: startLine})
			current = nil
			inBlock = false
		} else if inBlock {
			current = append(current, line)
		}
	}

	if inBlock && len(current) > 0 {
		blocks = append(blocks, RawBlock{Content: strings.Join(current, "\n"), StartLine: startLine})
	}

	return blocks
}
