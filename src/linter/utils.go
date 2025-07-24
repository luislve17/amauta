package linter

import (
	"html/template"
	"strings"
)

func getHTMLContent(raw string) template.HTML {
	return renderMarkdown(raw)
}

func extractSummary(raw string) string {
	var insideSummary bool
	var mdLines []string

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "summary: <md>" {
			insideSummary = true
			continue
		}

		if trimmed == "</md>" {
			insideSummary = false
			break
		}

		if strings.HasPrefix(trimmed, "summary:") && !strings.Contains(trimmed, "<md>") {
			summary := strings.TrimPrefix(trimmed, "summary:")
			return strings.TrimSpace(summary)
		}

		if insideSummary {
			mdLines = append(mdLines, line)
		}
	}
	return strings.Join(mdLines, "\n")
}
