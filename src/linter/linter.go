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
	root := initRoot()
	rawBlocks := extractRawBlocks(content)

	tags, tagsErr := getTags(rawBlocks)
	if tagsErr != nil {
		return nil, tagsErr
	}
	linkNodeOneToMany(root, tags)

	modules, modulesErr := getModules(rawBlocks)
	if modulesErr != nil {
		return nil, modulesErr
	}
	linkNodeOneToMany(root, modules)
	linkNodesManyToManyById("_tagIds", modules, tags)

	graph := StructureGraph{
		Root: root,
	}

	return &graph, nil
}

func initRoot() *Node {
	return &Node{
		Info: map[string]interface{}{
			"type": "root",
			"id":   "root",
		},
		Links: []*Node{},
	}
}

func linkNodeOneToOne(nodeA *Node, nodeB *Node) {
	nodeA.Links = append(nodeA.Links, nodeB)
	nodeB.Links = append(nodeB.Links, nodeA)
}

func linkNodeOneToMany(mainNode *Node, nodes []*Node) {
	for _, n := range nodes {
		n.Links = append(n.Links, mainNode)
	}
	mainNode.Links = append(mainNode.Links, nodes...)
}

func linkNodesManyToManyById(linkingKey string, nodesA []*Node, nodesB []*Node) {
	for _, nodeA := range nodesA {
		linkingIds := nodeA.Info[linkingKey].([]string)
		for _, nodeB := range nodesB {
			for _, linkId := range linkingIds {
				if nodeB.Info["id"] == linkId {
					linkNodeOneToOne(nodeA, nodeB)
				}
			}
		}
	}
}

func findSection(rawBlocks []RawBlock, rawRegex string, onlyFirst bool) []*RawBlock {
	headerPattern := regexp.MustCompile(rawRegex)
	var matches []*RawBlock

	for i := range rawBlocks {
		lines := strings.Split(rawBlocks[i].Content, "\n")
		if len(lines) == 0 {
			continue
		}
		header := strings.TrimSpace(lines[0])
		if headerPattern.MatchString(header) {
			if onlyFirst {
				return []*RawBlock{&rawBlocks[i]}
			}
			matches = append(matches, &rawBlocks[i])
		}
	}
	return matches
}

func getTags(rawBlocks []RawBlock) ([]*Node, error) {
	tagPattern := regexp.MustCompile(tagRegex)

	tagSection := findSection(rawBlocks, tagSectionRegex, true)
	if len(tagSection) == 0 {
		return nil, nil
	}

	lines := strings.Split(tagSection[0].Content, "\n")
	var nodes []*Node

	for i, tag := range lines[1:] {
		if tag == "" {
			continue
		}
		match := tagPattern.FindStringSubmatch(tag)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", tagSection[0].StartLine+i+1, tag)
		}

		node := &Node{
			Info: map[string]interface{}{
				"type":        "Tag",
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

func getModules(rawBlocks []RawBlock) ([]*Node, error) {
	moduleHeaderPatter := regexp.MustCompile(moduleSectionHeaderRegex)
	moduleSections := findSection(rawBlocks, moduleSectionHeaderRegex, false)
	if len(moduleSections) == 0 {
		return nil, nil
	}

	var nodes []*Node

	for i, moduleSection := range moduleSections {
		rawHeader := strings.Split(moduleSection.Content, "\n")[0]
		headerMatch := moduleHeaderPatter.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", moduleSection.StartLine+i+1, rawHeader)
		}
		node := &Node{
			Info: map[string]interface{}{
				"type":    "Module",
				"id":      headerMatch[1],
				"_tagIds": strings.Split(headerMatch[2], ","),
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func extractRawBlocks(content ManifestContent) []RawBlock {
	lines := strings.Split(string(content), "\n")
	var blocks []RawBlock
	var current []string
	var startLine int
	var inComment bool
	inBlock := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		isInlineCommentStart := strings.HasPrefix(trimmed, "--")
		isMultiLineCommentStart := strings.HasPrefix(trimmed, "<--")
		isMultiLineCommentEnd := strings.HasPrefix(trimmed, "-->")
		isNewModuleStart := strings.HasPrefix(trimmed, "[[")
		isNewInnerSectionStart := strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[")
		hasCurrentBlock := len(current) > 0

		if isMultiLineCommentEnd {
			inComment = false
			continue
		} else if isInlineCommentStart || inComment {
			continue
		} else if isMultiLineCommentStart {
			inComment = true
			continue
		} else if isNewModuleStart {
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
