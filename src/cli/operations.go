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

func resolveRefsInContent(rootPath string) string {
	refResolveErrorPrefix := "Unable to resolve ref declaration: %v"

	refLookupResult, lookupErr := findRefUsage(rootPath)
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
	return ""
}

func findRefUsage(rootPath string) (map[string][]regexLookupResult, error) {
	results := map[string][]regexLookupResult{
		"refImport":      {},
		"refDeclaration": {},
	}

	reImport := regexp.MustCompile(rawRefImportRegex)
	reDeclaration := regexp.MustCompile(rawRefDeclarationRegex)

	rootDepth := strings.Count(filepath.Clean(rootPath), string(os.PathSeparator))

	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Limit recursion depth
		currentDepth := strings.Count(filepath.Clean(path), string(os.PathSeparator)) - rootDepth
		if currentDepth > MAX_REF_LOOKUP_RECURSIVE_DEPTH {
			log.Fatalf("Max recursive depth (%d) reached at: %s", MAX_REF_LOOKUP_RECURSIVE_DEPTH, path)
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Skip files without allowed extension
		if !isAllowedExtension(d.Name()) {
			return nil
		}

		// Read and process file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		FindRefsWithRegexes(string(data), results, reImport, reDeclaration, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
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
		lineNumber := i + 1

		// ref imports
		if matches := reImport.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) > 1 {
					results["refImport"] = append(results["refImport"], regexLookupResult{
						Result:     match[1],
						LineNumber: lineNumber,
						FilePath:   filePath,
					})
				}
			}
		}

		// ref declarations
		if matches := reDeclaration.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				if len(match) > 1 {
					results["refDeclaration"] = append(results["refDeclaration"], regexLookupResult{
						Result:     match[1],
						LineNumber: lineNumber,
						FilePath:   filePath,
					})
				}
			}
		}
	}
}

func checkRefDeclarationUniqueness(refDeclarations []regexLookupResult) error {
	seen := make(map[string]regexLookupResult)

	for _, ref := range refDeclarations {
		if first, exists := seen[ref.Result]; exists {
			return fmt.Errorf(
				"duplicate '%s' ref found: first at %s:%d, again at %s:%d",
				ref.Result, first.FilePath, first.LineNumber, ref.FilePath, ref.LineNumber,
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
		for _, refDeclaration := range refLookupResults["refDecalaration"] {
			if refDeclaration.Result == refId {
				refFound = true
				break
			}
		}
		if !refFound {
			return fmt.Errorf("Missing ref declaration for '%s' usage at %s:%d", refId, refImport.FilePath, refImport.LineNumber)
		}
	}
	return nil
}
