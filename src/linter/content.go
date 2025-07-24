package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func getContent(rawBlocks []RawBlock) ([]*Node, error) {
	contentHeaderPattern := regexp.MustCompile(contentSectionHeaderRegex)
	contentSections, sectionFindErr := findSection(rawBlocks, contentSectionHeaderRegex, false, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

	if len(contentSections) == 0 {
		return nil, nil
	}

	nodes, contentParseErr := getContents(contentSections, contentHeaderPattern)
	if contentParseErr != nil {
		return nil, contentParseErr
	}

	return nodes, nil
}

func getContents(contentSections []*RawBlock, contentHeaderPattern *regexp.Regexp) ([]*Node, error) {
	var nodes []*Node

	for i, contentSection := range contentSections {
		rawHeader := strings.Split(contentSection.Content, "\n")[0]
		headerMatch := contentHeaderPattern.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid content format: %q", contentSection.StartLine+i+1, rawHeader)
		}
		contentData := getContentData(contentSection)
		node := &Node{
			Info:  createContentNodeInfo(headerMatch, contentData),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func getContentData(contentSection *RawBlock) map[string]interface{} {
	contentData := make(map[string]interface{})
	lines := strings.Split(contentSection.Content, "\n")
	for ln, line := range lines {

		if line == "" {
			continue
		}

		rawFieldRegex := `^([-_\w]+):\s*(.*)`
		fieldRegex := regexp.MustCompile(rawFieldRegex)
		fieldInfo := fieldRegex.FindStringSubmatch(line)

		if len(fieldInfo) < 3 {
			continue
		}

		key := fieldInfo[1]
		switch key {
		case "group":
			value := fieldInfo[2]
			contentData[key] = value
		case "summary":
			value := extractSummary(strings.Join(lines[ln:], "\n"))
			contentData[key] = value
		}
	}
	return contentData
}

func createContentNodeInfo(headerMatch []string, contentData map[string]interface{}) Content {
	return Content{
		Identifiable: Identifiable{
			Id: headerMatch[1],
		},
		BlockType: "Content",
		Summary:   getHTMLContent(contentData["summary"].(string)),
		LinkFields: LinkFields{
			_tagIds:   strings.Split(headerMatch[2], ","),
			_groupIds: []string{contentData["group"].(string)},
		},
	}
}
