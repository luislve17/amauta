package linter

import (
	"fmt"
	"regexp"
	"strings"
)

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
