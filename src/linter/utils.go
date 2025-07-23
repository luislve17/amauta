package linter

import (
	"html/template"
	"strings"
)

func getHTMLContent(raw string) template.HTML {
	lines := strings.Split(raw, "\n")
	var inMarkdown bool
	var mdLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "summary: <md>" {
			inMarkdown = true
			continue
		}

		if trimmed == "</md>" {
			inMarkdown = false
			break
		}

		if strings.HasPrefix(trimmed, "summary:") && !strings.Contains(trimmed, "<md>") {
			summary := strings.TrimPrefix(trimmed, "summary:")
			return renderMarkdown(strings.TrimSpace(summary))
		}

		if inMarkdown {
			mdLines = append(mdLines, line)
		}
	}

	return renderMarkdown(strings.Join(mdLines, "\n"))
}
