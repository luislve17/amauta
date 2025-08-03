package linter

import (
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
