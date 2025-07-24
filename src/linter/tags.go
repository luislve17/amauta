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

	nodes, tagFetchErr := getTagsInSection(*tagSection[0], tagPattern)
	if tagFetchErr != nil {
		return nil, tagFetchErr
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

func getTagsInSection(tagSection RawBlock, tagPattern *regexp.Regexp) ([]*Node, error) {
	foundTagNodes := []*Node{}
	sectionLines := strings.Split(tagSection.Content, "\n")
	for i, line := range sectionLines[1:] {
		if line == "" {
			continue
		}
		match := tagPattern.FindStringSubmatch(line)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", tagSection.StartLine+i+1, line)
		}
		node := &Node{
			Info:  createTagNodeInfo(match),
			Links: []*Node{},
		}
		foundTagNodes = append(foundTagNodes, node)
	}
	return foundTagNodes, nil
}
