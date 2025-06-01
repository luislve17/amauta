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

func getTags(rawBlocks []RawBlock) ([]*Node, error) {
	headerPattern := regexp.MustCompile(`^\[\[tags\]\]`)
	tagPatter := regexp.MustCompile(`^([-_\w]+)(#[A-F|\d]{6}):\s*(.*)`)

	var nodes []*Node

	for _, module := range rawBlocks {
		lines := strings.Split(module.Content, "\n")
		if len(lines) == 0 {
			continue
		}

		header := strings.TrimSpace(lines[0])
		foundTagSection := headerPattern.MatchString(header)
		if !foundTagSection {
			return nodes, nil
		}

		for lineNumber, tag := range lines[1:] {
			if tag == "" {
				break
			}

			tagMatch := tagPatter.FindStringSubmatch(tag)
			if len(tagMatch) == 0 {
				return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", module.StartLine+lineNumber+1, tag)
			}

			node := &Node{
				Info: map[string]interface{}{
					"type":        "tag",
					"id":          tagMatch[1],
					"color":       tagMatch[2],
					"description": tagMatch[3],
				},
				Links: []*Node{},
			}
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func getModuleNodes(rawBlocks []string) []*Node {
	headerPattern := regexp.MustCompile(`^\[\[(.+?)#(.*?)\]\]`)

	var nodes []*Node

	for _, module := range rawBlocks {
		lines := strings.Split(module, "\n")
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
		isNewSectionStart := strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[")
		hasCurrentBlock := len(current) > 0

		if isNewModuleStart {
			if inBlock && hasCurrentBlock {
				blocks = append(blocks, RawBlock{Content: strings.Join(current, "\n"), StartLine: startLine})
			}
			current = []string{line}
			startLine = i + 1
			inBlock = true
		} else if inBlock && isNewSectionStart {
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
