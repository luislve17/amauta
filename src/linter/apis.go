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

	nodesRegistry := NewNodeRegistry()

	var nodes []*Node

	for i, apiSection := range apiSections {
		rawHeader := strings.Split(apiSection.Content, "\n")[0]
		headerMatch := apiHeaderPattern.FindStringSubmatch(rawHeader)
		if len(headerMatch) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid tag format: %q", apiSection.From+i+1, rawHeader)
		}
		apiData, apiDataErr := getAPIData(apiSection, nodesRegistry)
		if apiDataErr != nil {
			return nil, apiDataErr
		}
		APINode := &Node{
			Info: createAPINodeInfo(headerMatch, apiData),
		}
		linkNodeOneToMany(APINode, apiData["innerNodes"].([]*Node))
		nodes = append(nodes, APINode)
	}

	return nodes, nil
}

func getAPIData(apiSection *RawBlock, nodesRegistry *NodeRegistry) (map[string]interface{}, error) {
	apiData := make(map[string]interface{})
	lines := strings.Split(apiSection.Content, "\n")
	innerBlocks := ExtractInnerBlocks(*apiSection)
	apiData["innerNodes"] = loadNodesFromInnerBlocks(innerBlocks, nodesRegistry)

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
			// TODO: Warning msg?
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

func loadInnerBlockDataNode(innerBlock InnerBlock, nodesRegistry *NodeRegistry) *Node {
	blockContent := strings.Split(innerBlock.Content, "\n")
	blockHeader := blockContent[0]
	apiRequestHeaderPattern := regexp.MustCompile(apiInnerSectionHeaderRegex)
	if apiRequestHeaderPattern.MatchString(blockHeader) {
		matches := apiRequestHeaderPattern.FindStringSubmatch(blockHeader)
		path := matches[3]
		httpVerb := matches[2]
		payloadType := matches[1]

		endpointNode := nodesRegistry.GetOrCreate(path, "Endpoint", func() *Node {
			return &Node{
				Info: Endpoint{
					Identifiable: Identifiable{Id: path},
					BlockType:    "Endpoint",
				},
			}
		})

		httpNode := nodesRegistry.GetOrCreate(httpVerb, "Http", func() *Node {
			return &Node{
				Info: HTTPVerb{
					Identifiable: Identifiable{Id: httpVerb},
					BlockType:    "Http",
				},
			}
		})

		payloadNode := &Node{
			Info: loadPayloadInfo(blockContent[1:], payloadType),
		}

		linkNodeOneToOne(endpointNode, httpNode)
		linkNodeOneToOne(httpNode, payloadNode)
		return endpointNode
	}
	return nil
}

func loadNodesFromInnerBlocks(innerBlocks []InnerBlock, nodesRegistry *NodeRegistry) []*Node {
	var result []*Node
	seen := make(map[*Node]bool)

	for _, block := range innerBlocks {
		innerBlockNode := loadInnerBlockDataNode(block, nodesRegistry)
		if innerBlockNode != nil && !seen[innerBlockNode] {
			result = append(result, innerBlockNode)
			seen[innerBlockNode] = true
		}
	}

	return result
}

func loadPayloadInfo(blockContent []string, payloadType string) NodeInfo {
	data := map[string]string{}
	for ln, line := range blockContent {
		if line == "" {
			continue
		}

		rawFieldRegex := `^([-_\w]+):\s*(.*)`
		fieldRegex := regexp.MustCompile(rawFieldRegex)
		fieldInfo := fieldRegex.FindStringSubmatch(line)

		if fieldInfo == nil {
			continue
		}
		key := fieldInfo[1]
		switch key {
		case "summary":
			value, nLines := extractSummary(strings.Join(blockContent[ln:], "\n"))
			ln += nLines
			data[key] = value
		default:
			data[key] = fieldInfo[2]
		}
	}
	if payloadType == "request" {
		payloadInfo := RequestPayload{
			Summary:   getHTMLContent(data["summary"]),
			BlockType: "RequestPayload",
		}
		return payloadInfo
	}

	if payloadType == "response" {
		payloadInfo := RequestPayload{
			Summary:   getHTMLContent(data["summary"]),
			BlockType: "ResponsePayload",
		}
		return payloadInfo
	}
	return nil
}
