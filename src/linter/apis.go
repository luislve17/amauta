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
	innerBlocks := ExtractInnerBlocks(*apiSection)
	for ln, line := range lines {

		if line == "" || lineOverlapsLineRanges(ln, innerBlocks) {
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

func ExtractInnerBlocks(rawBlock RawBlock) []InnerBlock {
	lines := strings.Split(rawBlock.Content, "\n")
	var blocks []InnerBlock
	var current []string
	var inBlock bool
	var inComment bool
	var inCodeBlock bool
	var inMD bool

	startLine := rawBlock.LineRange.From

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCodeBlock = !inCodeBlock
		}

		if !inMD && !inCodeBlock && strings.HasPrefix(trimmed, "summary:") && strings.Contains(trimmed, "<md>") {
			inMD = true
		}

		if inMD && !inCodeBlock && strings.TrimSpace(trimmed) == "</md>" {
			inMD = false
		}

		// Skip comments (only outside markdown)
		if !inMD {
			if strings.HasPrefix(trimmed, "-->") {
				inComment = false
				continue
			}
			if inComment || strings.HasPrefix(trimmed, "--") {
				continue
			}
			if strings.HasPrefix(trimmed, "<--") {
				inComment = true
				continue
			}
		}

		if strings.HasPrefix(trimmed, "[") && !strings.HasPrefix(trimmed, "[[") && !inMD && !inCodeBlock {
			if inBlock && len(current) > 0 {
				blocks = append(blocks, InnerBlock{
					Content: strings.Join(current, "\n"),
					LineRange: LineRange{
						From: startLine,
						To:   i,
					},
				})
			}
			current = []string{line}
			startLine = i + 1
			inBlock = true
		} else if inBlock {
			current = append(current, line)
		}
	}

	// Append last block
	if inBlock && len(current) > 0 {
		blocks = append(blocks, InnerBlock{
			Content: strings.Join(current, "\n"),
			LineRange: LineRange{
				From: startLine,
				To:   len(lines),
			},
		})
	}

	return blocks
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
