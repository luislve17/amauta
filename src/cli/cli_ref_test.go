package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	// "github.com/luislve17/amauta/linter"
	"github.com/stretchr/testify/assert"
)

// Test uniqueness on refs ids
// Lint from root file
// If file doesnt have a explicit root, alert
// Alert contains assumed root path
// Root file

func TestFindsRefUsageWithinRootFolderPath(t *testing.T) {
	assert := assert.New(t)

	// Create a temp dir to hold all test files
	tempDir, err := os.MkdirTemp("", "amauta-test-*")
	assert.NoError(err)
	defer os.RemoveAll(tempDir) // cleans up everything inside

	// Root file
	rootPath := filepath.Join(tempDir, "root.amauta")
	assert.NoError(os.WriteFile(rootPath, []byte(rootFileContent), 0644))

	// Another file in the same directory as root
	subPath := filepath.Join(tempDir, "sub.amauta")
	assert.NoError(os.WriteFile(subPath, []byte(nonRootFileContent), 0644))

	// File inside a subfolder
	subDir := filepath.Join(tempDir, "nested")
	assert.NoError(os.Mkdir(subDir, 0755))

	nestedFilePath := filepath.Join(subDir, "nested.amauta")
	assert.NoError(os.WriteFile(nestedFilePath, []byte(nestedFileContent), 0644))

	// Run the method starting from root temp dir
	results := findRefUsage(tempDir)

	// Assertions
	assert.Len(results["refImport"], 2)
	assert.Len(results["refDeclaration"], 2)

	assert.Equal(regexLookupResult{
		Result:     "my-tags",
		LineNumber: 5,
		FilePath:   rootPath,
	}, results["refImport"][0])
	assert.Equal(regexLookupResult{
		Result:     "my-groups",
		LineNumber: 8,
		FilePath:   rootPath,
	}, results["refImport"][1])

	assert.Equal(regexLookupResult{
		Result:     "my-tags",
		LineNumber: 1,
		FilePath:   nestedFilePath,
	}, results["refDeclaration"][0])
	assert.Equal(regexLookupResult{
		Result:     "my-groups",
		LineNumber: 2,
		FilePath:   subPath,
	}, results["refDeclaration"][1])
}

func TestFindsDuplicatesInRefDeclarations(t *testing.T) {
	assert := assert.New(t)

	tempDir, err := os.MkdirTemp("", "amauta-test-*")
	assert.NoError(err)
	defer os.RemoveAll(tempDir) // cleans up everything inside

	// Files in root dir
	rootPath := filepath.Join(tempDir, "root.amauta")
	assert.NoError(os.WriteFile(rootPath, []byte(rootFileContent), 0644))

	otherPath1 := filepath.Join(tempDir, "sub1.amauta")
	assert.NoError(os.WriteFile(otherPath1, []byte(nonRootFileContent), 0644))

	otherPath2 := filepath.Join(tempDir, "sub2.amauta")
	assert.NoError(os.WriteFile(otherPath2, []byte(duplicatedRefDeclaration), 0644))

	// File inside a subfolder
	subDir := filepath.Join(tempDir, "nested")
	assert.NoError(os.Mkdir(subDir, 0755))
	nestedFilePath := filepath.Join(subDir, "nested.amauta")
	assert.NoError(os.WriteFile(nestedFilePath, []byte(nestedFileContent), 0644))

	results := findRefUsage(tempDir)
	checkUniqueErr := checkRefDeclarationUniqueness(results["refDeclaration"])
	assert.EqualError(checkUniqueErr, fmt.Sprintf("duplicate 'my-tags' ref found: first at %s:%d, again at %s:%d", nestedFilePath, 1, otherPath2, 2))
}
