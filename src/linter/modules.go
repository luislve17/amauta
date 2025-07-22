package linter

import (
	"fmt"
	"regexp"
	"strings"
)

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
			Info:  createModuleNodeInfo(headerMatch, moduleSection),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func createModuleNodeInfo(headerMatch []string, moduleSection *RawBlock) Module {
	return Module{
		Identifiable: Identifiable{
			Id: headerMatch[1],
		},
		BlockType: "Module",
		Summary:   getHTMLContent(moduleSection.Content),
		LinkFields: LinkFields{
			_tagIds:   strings.Split(headerMatch[2], ","),
			_groupIds: strings.Split(getGroupsInSection(moduleSection.Content), ","),
		},
	}
}
