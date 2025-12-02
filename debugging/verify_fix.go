package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

func main() {
	fmt.Println("=== Verifying File Type Classification Fix ===\n")

	// Test the classification logic
	testCases := []struct {
		filePath     string
		expectedType parser.FileType
	}{
		{"MEMORY.md", parser.FileTypeMemory},
		{"decisions.md", parser.FileTypeMemory},
		{"docs/MEMORY.md", parser.FileTypeMemory},     // Should be memory, not docs
		{"project/decisions.md", parser.FileTypeMemory}, // Should be memory, not src
		{"README.md", parser.FileTypeDocs},
		{"docs/api.md", parser.FileTypeDocs},
		{"test/example_test.go", parser.FileTypeTests},
		{"src/main.go", parser.FileTypeSrc},
	}

	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "sourcerer-test-*")
	if err != nil {
		fmt.Printf("Failed to create temp dir: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	for _, tc := range testCases {
		fullPath := filepath.Join(tmpDir, tc.filePath)
		os.MkdirAll(filepath.Dir(fullPath), 0755)

		content := fmt.Sprintf("# Test File: %s\n\nSome content here.", tc.filePath)
		os.WriteFile(fullPath, []byte(content), 0644)
	}

	// Test with markdown parser
	mdParser, err := parser.NewMarkdownParser(tmpDir)
	if err != nil {
		fmt.Printf("Failed to create parser: %v\n", err)
		os.Exit(1)
	}
	defer mdParser.Close()

	allPassed := true
	for _, tc := range testCases {
		// Only test markdown files
		if filepath.Ext(tc.filePath) != ".md" {
			continue
		}

		file, err := mdParser.Chunk(tc.filePath)
		if err != nil {
			fmt.Printf("❌ FAILED: %s - couldn't parse: %v\n", tc.filePath, err)
			allPassed = false
			continue
		}

		if len(file.Chunks) == 0 {
			fmt.Printf("⚠️  WARNING: %s - no chunks extracted\n", tc.filePath)
			continue
		}

		actualType := parser.FileType(file.Chunks[0].Type)
		if actualType == tc.expectedType {
			fmt.Printf("✓ PASS: %-30s → %s\n", tc.filePath, actualType)
		} else {
			fmt.Printf("❌ FAIL: %-30s → got %s, expected %s\n", tc.filePath, actualType, tc.expectedType)
			allPassed = false
		}
	}

	fmt.Println()
	if allPassed {
		fmt.Println("=== ✅ All classification tests PASSED ===")
		fmt.Println("\nNext steps:")
		fmt.Println("1. Delete your existing .sourcerer/db directory to clear old misclassified data")
		fmt.Println("2. Rebuild: go build -o sourcerer cmd/sourcerer/main.go")
		fmt.Println("3. Re-index your workspace to pick up the fixes")
	} else {
		fmt.Println("=== ❌ Some tests FAILED ===")
		os.Exit(1)
	}
}
