package linter

import (
	"fmt"
	"regexp"
	"strings"
)

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
			Info:  createTagNodeInfo(match),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func createTagNodeInfo(tagMatch []string) Tag {
	return Tag{
		BlockType: "Tag",
		Identifiable: Identifiable{
			Id: tagMatch[1],
		},
		color:       tagMatch[2],
		Description: tagMatch[3],
	}
}
