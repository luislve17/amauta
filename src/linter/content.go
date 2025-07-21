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
