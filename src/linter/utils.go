package linter

import (
	"fmt"
	"html/template"
	"strings"
)

func getHTMLContent(raw string) template.HTML {
	return renderMarkdown(raw)
}

func extractSummary(raw string) (string, int) {
	var insideSummary bool
	var mdLines []string

	lines := strings.Split(raw, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "summary: <md>" {
			insideSummary = true
			continue
		}

		if trimmed == "</md>" {
			insideSummary = false
			break
		}

		if insideSummary {
			mdLines = append(mdLines, line)
			continue
		}

		if strings.HasPrefix(trimmed, "summary:") {
			summary := strings.TrimPrefix(trimmed, "summary:")
			return strings.TrimSpace(summary), i
		}
	}

	return strings.Join(mdLines, "\n"), len(lines)
}

// TEST
func lineInRanges(line int, ranges []LineRange) bool {
	for _, r := range ranges {
		if line >= r.From && line < r.To {
			return true
		}
	}
	return false
}

// TEST
func getLinesToSkipInSections(rawBlocks []RawBlock) ([]LineRange, error) {
	var skipRanges []LineRange
	var globalLine int

	for _, block := range rawBlocks {
		lines := strings.Split(block.Content, "\n")
		var insideSummary bool
		var startLine int

		for ln := 0; ln < len(lines); ln++ {
			trimmed := strings.TrimSpace(lines[ln])

			if trimmed == "summary: <md>" {
				if insideSummary {
					return nil, fmt.Errorf("nested <md> found at global line %d", globalLine)
				}
				insideSummary = true
				startLine = globalLine
				globalLine++
				continue
			}

			if insideSummary {
				if trimmed == "</md>" {
					skipRanges = append(skipRanges, LineRange{From: startLine, To: globalLine + 1})
					insideSummary = false
				}
				globalLine++
				continue
			}

			if strings.HasPrefix(trimmed, "summary:") && !strings.Contains(trimmed, "<md>") {
				skipRanges = append(skipRanges, LineRange{From: globalLine, To: globalLine + 1})
			}

			globalLine++
		}

		if insideSummary {
			return nil, fmt.Errorf("unterminated <md> block starting at global line %d", startLine+1)
		}
	}

	return skipRanges, nil
}
