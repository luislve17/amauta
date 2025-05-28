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
	if err := errors.Join(contentGraphErr); err != nil {
		lintStatus = LintStatusError
		resultMsg = "error"
	}

	return LintResult{Status: lintStatus, Msg: resultMsg, Structure: graph}, nil
}

func generateGraph(content ManifestContent) (*StructureGraph, error) {
	root := &Node{
		Info: map[string]interface{}{
			"type": "root",
			"name": "root",
		},
		Children: []*Node{},
	}

	rawBlocks := extractRawBlocks(content)
	modules := getModuleNodes(rawBlocks)
	root.Children = append(root.Children, modules...)

	graph := StructureGraph{
		Root: root,
	}

	graph.printStructure()

	return &graph, nil
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
			Children: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func extractRawBlocks(content ManifestContent) []string {
	lines := strings.Split(string(content), "\n")
	var blocks []string
	var current []string
	inBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		isNewModuleStart := strings.HasPrefix(trimmed, "[[")
		isNewSectionStart := strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[")
		hasCurrentBlock := len(current) > 0

		if isNewModuleStart {
			if inBlock && hasCurrentBlock {
				blocks = append(blocks, strings.Join(current, "\n"))
			}
			current = []string{line}
			inBlock = true
		} else if inBlock && isNewSectionStart {
			blocks = append(blocks, strings.Join(current, "\n"))
			current = nil
			inBlock = false
		} else if inBlock {
			current = append(current, line)
		}
	}

	if inBlock && len(current) > 0 {
		blocks = append(blocks, strings.Join(current, "\n"))
	}

	return blocks
}

func (g *StructureGraph) printStructure() {
	fmt.Println(": : : : : : : : : : : : : : : : ")
	g.Root.printTree("", true)
	fmt.Println(": : : : : : : : : : : : : : : : ")
}

func (n *Node) printTree(prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	fmt.Printf("%s%s'%s'(%s)\n", prefix, connector, n.Info["name"], n.Info["type"])

	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	for i, child := range n.Children {
		last := i == len(n.Children)-1
		child.printTree(newPrefix, last)
	}
}
