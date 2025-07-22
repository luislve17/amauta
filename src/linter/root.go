package linter

import (
	"fmt"
	"regexp"
	"strings"
)

func initRoot(rawBlocks []RawBlock) (*Node, error) {
	root := &Node{
		Info:  createRootNodeInfo(),
		Links: []*Node{},
	}

	rootSection, sectionFindErr := findSection(rawBlocks, rootSectionRegex, true, false)
	if sectionFindErr != nil {
		return root, nil
	}

	fields := strings.Split(rootSection[0].Content, "\n")
	rootFieldPattern := regexp.MustCompile(rootFieldRegex)

	for i, field := range fields[1:] {
		if field == "" {
			continue
		}
		match := rootFieldPattern.FindStringSubmatch(field)
		if len(match) == 0 {
			return nil, fmt.Errorf("Error@line:%d\n->Invalid root field format: %q", rootSection[0].StartLine+i+1, field)
		}
		value := match[2]
		rootData, ok := root.Info.(*Root)
		if !ok {
			return nil, fmt.Errorf("unexpected Info type for root node")
		}

		switch key := match[1]; key {
		case "LogoUrl":
			rootData.LogoUrl = value
		case "GithubUrl":
			rootData.GithubUrl = value
		}
	}

	return root, nil
}

func createRootNodeInfo() *Root {
	return &Root{
		Identifiable: Identifiable{
			Id: "Root",
		},
		BlockType: "Root",
	}
}
