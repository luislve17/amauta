package linter

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
	goldmark.WithRendererOptions(html.WithUnsafe()),
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

func generateGraph(manifestContent ManifestContent) (*StructureGraph, error) {
	rawBlocks := extractRawBlocks(manifestContent)
	root, rootErr := initRoot(rawBlocks)
	if rootErr != nil {
		return nil, rootErr
	}

	groups, groupsErr := getGroups(rawBlocks)
	if groupsErr != nil {
		return nil, groupsErr
	}
	linkNodeOneToMany(root, groups)

	contents, contentErr := getContent(rawBlocks)
	if contentErr != nil {
		return nil, contentErr
	}

	tags, tagsErr := getTags(rawBlocks)
	if tagsErr != nil {
		return nil, tagsErr
	}
	linkNodeOneToMany(root, tags)

	modules, modulesErr := getModules(rawBlocks)
	if modulesErr != nil {
		return nil, modulesErr
	}

	linkNodesManyToManyById("_tagIds", modules, tags)
	linkNodesManyToManyById("_tagIds", groups, tags)

	linkNodesManyToManyById("_groupIds", modules, groups)
	linkNodesManyToManyById("_groupIds", contents, groups)

	graph := StructureGraph{
		Root: root,
	}

	return &graph, nil
}

func initRoot(rawBlocks []RawBlock) (*Node, error) {
	root := &Node{
		Info: map[string]interface{}{
			"type": "root",
			"id":   "root",
		},
		Links: []*Node{},
	}

	rootSection, sectionFindErr := findSection(rawBlocks, rootSectionRegex, true, false)
	if sectionFindErr != nil {
		return root, nil
	}

	fields := strings.Split(rootSection[0].Content, "\n")
	rootFieldPattern := regexp.MustCompile(rootFieldRegex)

	for i, field := range fields[1:] {
		if field == "" {
			continue
		}
		match := rootFieldPattern.FindStringSubmatch(field)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid root field format: %q", rootSection[0].StartLine+i+1, field)
		}
		key := match[1]
		value := match[2]
		root.Info[key] = value
	}

	return root, nil
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

func findSection(rawBlocks []RawBlock, rawRegex string, onlyOne bool, allowNone bool) ([]*RawBlock, error) {
	headerPattern := regexp.MustCompile(rawRegex)
	var matches []*RawBlock

	for i := range rawBlocks {
		lines := strings.Split(rawBlocks[i].Content, "\n")
		if len(lines) == 0 {
			continue
		}
		header := strings.TrimSpace(lines[0])
		if headerPattern.MatchString(header) {
			matches = append(matches, &rawBlocks[i])
		}
	}

	if onlyOne {
		if len(matches) == 1 {
			return []*RawBlock{matches[0]}, nil
		}
		if len(matches) > 1 {
			return nil, fmt.Errorf("multiple sections matched regex: %s", rawRegex)
		}
		if allowNone {
			return nil, nil
		} else {
			return nil, fmt.Errorf("no section matched regex: %s", rawRegex)
		}
	}

	return matches, nil
}

func getGroups(rawBlocks []RawBlock) ([]*Node, error) {
	groupPattern := regexp.MustCompile(groupRegex)

	groupSection, sectionFindErr := findSection(rawBlocks, groupSectionRegex, true, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

	if len(groupSection) == 0 {
		return nil, nil
	}

	lines := strings.Split(groupSection[0].Content, "\n")
	var nodes []*Node

	for i, group := range lines[1:] {
		if group == "" {
			continue
		}
		match := groupPattern.FindStringSubmatch(group)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid group format: %q", groupSection[0].StartLine+i+1, group)
		}

		node := &Node{
			Info: map[string]interface{}{
				"type":        "Group",
				"id":          match[1],
				"description": match[3],
				"_tagIds":     strings.Split(match[2], ","),
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getTags(rawBlocks []RawBlock) ([]*Node, error) {
	tagPattern := regexp.MustCompile(tagRegex)

	tagSection, sectionFindErr := findSection(rawBlocks, tagSectionRegex, true, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

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
	moduleHeaderPattern := regexp.MustCompile(moduleSectionHeaderRegex)
	moduleSections, sectionFindErr := findSection(rawBlocks, moduleSectionHeaderRegex, false, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

	if len(moduleSections) == 0 {
		return nil, nil
	}

	var nodes []*Node

	for i, moduleSection := range moduleSections {
		rawHeader := strings.Split(moduleSection.Content, "\n")[0]
		headerMatch := moduleHeaderPattern.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", moduleSection.StartLine+i+1, rawHeader)
		}
		node := &Node{
			Info: map[string]interface{}{
				"type":        "Module",
				"id":          headerMatch[1],
				"htmlContent": getHTMLContent(moduleSection.Content),
				"_tagIds":     strings.Split(headerMatch[2], ","),
				"_groupIds":   strings.Split(getGroupsInSection(moduleSection.Content), ","),
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getContent(rawBlocks []RawBlock) ([]*Node, error) {
	contentHeaderPattern := regexp.MustCompile(contentSectionHeaderRegex)
	contentSections, sectionFindErr := findSection(rawBlocks, contentSectionHeaderRegex, false, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

	if len(contentSections) == 0 {
		return nil, nil
	}

	var nodes []*Node

	for i, contentSection := range contentSections {
		rawHeader := strings.Split(contentSection.Content, "\n")[0]
		headerMatch := contentHeaderPattern.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid content format: %q", contentSection.StartLine+i+1, rawHeader)
		}
		node := &Node{
			Info: map[string]interface{}{
				"type":        "Content",
				"id":          headerMatch[1],
				"htmlContent": getHTMLContent(contentSection.Content),
				"_tagIds":     strings.Split(headerMatch[2], ","),
				"_groupIds":   strings.Split(getGroupsInSection(contentSection.Content), ","),
			},
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getGroupsInSection(raw string) string {
	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		prefix := "group:"
		trimmed := strings.ReplaceAll(line, " ", "")
		if strings.HasPrefix(trimmed, prefix) {
			return trimmed[len(prefix):]
		}
	}
	return ""
}

func getHTMLContent(raw string) template.HTML {
	lines := strings.Split(raw, "\n")
	var inMarkdown bool
	var mdLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "summary: <md>" {
			inMarkdown = true
			continue
		}

		if trimmed == "</md>" {
			inMarkdown = false
			break
		}

		if strings.HasPrefix(trimmed, "summary:") && !strings.Contains(trimmed, "<md>") {
			summary := strings.TrimPrefix(trimmed, "summary:")
			return renderMarkdown(strings.TrimSpace(summary))
		}

		if inMarkdown {
			mdLines = append(mdLines, line)
		}
	}

	return renderMarkdown(strings.Join(mdLines, "\n"))
}

func renderMarkdown(content string) template.HTML {
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return template.HTML(content)
	}
	return template.HTML(buf.String())
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
