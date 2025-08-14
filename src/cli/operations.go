package cli

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/luislve17/amauta/linter"
)

func loadManifestFromFile(path string) linter.ManifestContent {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	return linter.ManifestContent(content)
}

func loadThemeFromFile(path string) (template.CSS, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return template.CSS(""), err
	}
	return template.CSS(content), nil
}

func (str styledString) style(args ...string) string {
	prefix := ""
	for _, arg := range args {
		prefix += arg
	}
	return fmt.Sprintf("%s%s%s", prefix, str, reset)
}

func resolveRefsInContent(rootPath string) linter.ManifestContent {
	rootDir := filepath.Dir(rootPath)
	refResolveErrorPrefix := "Unable to resolve ref declaration: %v"

	refLookupResult, fileContents, lookupErr := findRefUsage(rootDir)
	if lookupErr != nil {
		log.Fatalf(refResolveErrorPrefix, lookupErr)
	}

	refDeclarationUniquenessErr := checkRefDeclarationUniqueness(refLookupResult["refDeclaration"])
	if refDeclarationUniquenessErr != nil {
		log.Fatalf(refResolveErrorPrefix, refDeclarationUniquenessErr)
	}

	refUsageMissingDeclarationErr := checkMissingRefDeclarations(refLookupResult)
	if refUsageMissingDeclarationErr != nil {
		log.Fatalf(refResolveErrorPrefix, refUsageMissingDeclarationErr)
	}

	loadedRefFiles, loadingRefErr := replaceRefsInFileContents(refLookupResult, fileContents)
	if loadingRefErr != nil {
		log.Fatalf(refResolveErrorPrefix, loadingRefErr)
	}

	resolvedManifestContent := unifyFileContents(loadedRefFiles)
	return linter.ManifestContent(resolvedManifestContent)
}

func findRefUsage(rootPath string) (map[string][]regexLookupResult, map[string]string, error) {
	refResults := map[string][]regexLookupResult{
		"refImport":      {},
		"refDeclaration": {},
	}
	fileContents := make(map[string]string)

	reImport := regexp.MustCompile(rawRefImportRegex)
	reDeclaration := regexp.MustCompile(rawRefDeclarationRegex)

	rootDepth := strings.Count(filepath.Clean(rootPath), string(os.PathSeparator))

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		currentDepth := strings.Count(filepath.Clean(path), string(os.PathSeparator)) - rootDepth
		if currentDepth > MAX_REF_LOOKUP_RECURSIVE_DEPTH {
			log.Fatalf("Max recursive depth (%d) reached at: %s", MAX_REF_LOOKUP_RECURSIVE_DEPTH, path)
		}
		if d.IsDir() || !isAllowedExtension(d.Name()) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		fileContents[path] = string(data)
		FindRefsWithRegexes(string(data), refResults, reImport, reDeclaration, path)
		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return refResults, fileContents, nil
}

func isAllowedExtension(filename string) bool {
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	for _, allowed := range allowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func FindRefsWithRegexes(input string, results map[string][]regexLookupResult, reImport, reDeclaration *regexp.Regexp, filePath string) {
	lines := strings.Split(input, "\n")

	for i, line := range lines {
		if matches := reImport.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) > 1 {
					results["refImport"] = append(results["refImport"], regexLookupResult{
						Result:    match[1],
						FilePath:  filePath,
						LineRange: linter.LineRange{From: i + 1, To: i + 1},
					})
				}
			}
		}
	}

	blocks := linter.ExtractRawBlocks(linter.ManifestContent(input))
	for _, block := range blocks {
		lines := strings.Split(block.Content, "\n")
		if len(lines) == 0 {
			continue
		}
		declarationLineRange := linter.LineRange{From: block.From + 1, To: block.To - 1}

		if matches := reDeclaration.FindStringSubmatch(strings.TrimSpace(lines[0])); matches != nil && len(matches) > 1 {
			results["refDeclaration"] = append(results["refDeclaration"], regexLookupResult{
				Result:    matches[1],
				FilePath:  filePath,
				LineRange: declarationLineRange,
			})
		}
	}
}

func checkRefDeclarationUniqueness(refDeclarations []regexLookupResult) error {
	seen := make(map[string]regexLookupResult)

	for _, ref := range refDeclarations {
		if first, exists := seen[ref.Result]; exists {
			return fmt.Errorf(
				"duplicate '%s' ref found: first at %s:%d, again at %s:%d",
				ref.Result, first.FilePath, first.LineRange.From, ref.FilePath, ref.LineRange.From,
			)
		}
		seen[ref.Result] = ref
	}

	return nil
}

func checkMissingRefDeclarations(refLookupResults map[string][]regexLookupResult) error {
	for _, refImport := range refLookupResults["refImport"] {
		refId := refImport.Result
		refFound := false
		for _, refDeclaration := range refLookupResults["refDeclaration"] {
			if refDeclaration.Result == refId {
				refFound = true
				break
			}
		}
		if !refFound {
			return fmt.Errorf("Missing ref declaration for '%s' usage at %s:%d", refId, refImport.FilePath, refImport.LineRange.From)
		}
	}
	return nil
}

func replaceRefsInFileContents(
	refLookupResults map[string][]regexLookupResult,
	fileContents map[string]string,
) (map[string]string, error) {
	reImport := regexp.MustCompile(rawRefImportRegex)

	// Build declaration bodies map directly from refLookupResults + fileContents
	declBodies := make(map[string]string)
	for _, decl := range refLookupResults["refDeclaration"] {
		lines := strings.Split(fileContents[decl.FilePath], "\n")
		body := strings.Join(lines[decl.LineRange.From-1:decl.LineRange.To], "\n")
		declBodies[decl.Result] = body
	}

	for pass := 0; pass < MAX_REF_REPLACEMENT_ITERATIONS; pass++ {
		changed := false

		for path, content := range fileContents {
			newContent := reImport.ReplaceAllStringFunc(content, func(m string) string {
				match := reImport.FindStringSubmatch(m)
				if len(match) < 2 {
					return m
				}
				id := match[1]
				if body, ok := declBodies[id]; ok {
					changed = true
					return body
				}
				return m
			})
			fileContents[path] = newContent
		}

		if !changed {
			return fileContents, nil
		}
	}

	return nil, fmt.Errorf("too many iterations, possible cyclic refs")
}

func unifyFileContents(fileContents map[string]string) string {
	var builder strings.Builder
	for _, content := range fileContents {
		builder.WriteString(content)
		builder.WriteString("\n")
	}
	return builder.String()
}
