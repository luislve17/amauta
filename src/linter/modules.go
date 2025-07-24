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
		moduleData := getModuleData(moduleSection)
		node := &Node{
			Info:  createModuleNodeInfo(headerMatch, moduleData),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getModuleData(moduleSection *RawBlock) map[string]interface{} {
	moduleData := make(map[string]interface{})
	lines := strings.Split(moduleSection.Content, "\n")
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
			moduleData[key] = value
		case "summary":
			value := extractSummary(strings.Join(lines[ln:], "\n"))
			moduleData[key] = value
		}
	}
	return moduleData
}

func createModuleNodeInfo(headerMatch []string, moduleData map[string]interface{}) Module {
	return Module{
		Identifiable: Identifiable{
			Id: headerMatch[1],
		},
		BlockType: "Module",
		Summary:   getHTMLContent(moduleData["summary"].(string)),
		LinkFields: LinkFields{
			_tagIds:   strings.Split(headerMatch[2], ","),
			_groupIds: []string{moduleData["group"].(string)},
		},
	}
}
