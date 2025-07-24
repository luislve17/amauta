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

	nodes, groupFetchErr := getGroupsInSection(groupSection[0], groupPattern)
	if groupFetchErr != nil {
		return nil, groupFetchErr
	}
	return nodes, nil
}

func getGroupsInSection(groupSection *RawBlock, groupPattern *regexp.Regexp) ([]*Node, error) {
	var nodes []*Node
	lines := strings.Split(groupSection.Content, "\n")

	for i, group := range lines[1:] {
		if group == "" {
			continue
		}
		match := groupPattern.FindStringSubmatch(group)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid group format: %q", groupSection.StartLine+i+1, group)
		}

		node := &Node{
			Info:  createGroupNodeInfo(match),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func createGroupNodeInfo(groupMatch []string) Group {
	return Group{
		BlockType: "Group",
		Identifiable: Identifiable{
			Id: groupMatch[1],
		},
		Description: groupMatch[3],
		LinkFields: LinkFields{
			_tagIds: strings.Split(groupMatch[2], ","),
		},
	}
}
