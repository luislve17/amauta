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
			return nil, fmt.Errorf("Error@line:%d\n->Invalid content format: %q", contentSection.From+i+1, rawHeader)
		}
		contentData, contentErr := getContentData(contentSection)
		if contentErr != nil {
			return nil, contentErr
		}
		node := &Node{
			Info:  createContentNodeInfo(headerMatch, contentData),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func getContentData(contentSection *RawBlock) (map[string]interface{}, error) {
	contentData := make(map[string]interface{})
	lines := strings.Split(contentSection.Content, "\n")
	for ln := 0; ln < len(lines); ln++ {
		line := lines[ln]
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
		case "summary":
			value, nLines := extractSummary(strings.Join(lines[ln:], "\n"))
			ln += nLines
			contentData[key] = value
		case "group":
			value := fieldInfo[2]
			contentData[key] = value
		default:
			return nil, fmt.Errorf("Error@line:%d\n->Invalid field found: '%s'", contentSection.From+ln, key)
		}
	}
	return contentData, nil
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
