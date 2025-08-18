package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func getAPIs(rawBlocks []RawBlock) ([]*Node, error) {
	apiHeaderPattern := regexp.MustCompile(apiSectionHeaderRegex)
	apiSections, sectionFindErr := findSection(rawBlocks, apiSectionHeaderRegex, false, true)
	if sectionFindErr != nil {
		return nil, sectionFindErr
	}

	if len(apiSections) == 0 {
		return nil, nil
	}

	var nodes []*Node

	for i, apiSection := range apiSections {
		rawHeader := strings.Split(apiSection.Content, "\n")[0]
		headerMatch := apiHeaderPattern.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", apiSection.From+i+1, rawHeader)
		}
		apiData, apiDataErr := getAPIData(apiSection)
		if apiDataErr != nil {
			return nil, apiDataErr
		}
		node := &Node{
			Info:  createAPINodeInfo(headerMatch, apiData),
			Links: []*Node{},
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func getAPIData(apiSection *RawBlock) (map[string]interface{}, error) {
	apiData := make(map[string]interface{})
	lines := strings.Split(apiSection.Content, "\n")
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
		case "summary":
			value, nLines := extractSummary(strings.Join(lines[ln:], "\n"))
			ln += nLines
			apiData[key] = value
		case "group":
			value := fieldInfo[2]
			apiData[key] = value
		default:
			continue
		}
	}
	return apiData, nil
}

func createAPINodeInfo(headerMatch []string, apiData map[string]interface{}) API {
	return API{
		Identifiable: Identifiable{
			Id: headerMatch[1],
		},
		BlockType: "API",
		Summary:   getHTMLContent(apiData["summary"].(string)),
		LinkFields: LinkFields{
			_tagIds:   strings.Split(headerMatch[2], ","),
			_groupIds: []string{apiData["group"].(string)},
		},
	}
}
