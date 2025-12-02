package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run debug_file_types.go <workspace_path>")
		fmt.Println("Example: go run debug_file_types.go /Users/claudeuser/claude_home/test_project")
		os.Exit(1)
	}

	testDir := os.Args[1]
	fmt.Printf("=== File Type Classification Debug ===\n")
	fmt.Printf("Workspace: %s\n\n", testDir)

	// Create markdown parser
	mdParser, err := parser.NewMarkdownParser(testDir)
	if err != nil {
		fmt.Printf("Failed to create parser: %v\n", err)
		os.Exit(1)
	}

	// Find all markdown files
	var mdFiles []string
	err = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			// Get relative path from workspace root
			relPath, _ := filepath.Rel(testDir, path)
			mdFiles = append(mdFiles, relPath)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Failed to walk directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d markdown files:\n\n", len(mdFiles))

	// Test classification for each file
	for _, relPath := range mdFiles {
		// We need to access the private classifyFileType method
		// So instead we'll check the FileTypeRules directly
		baseName := filepath.Base(relPath)

		fmt.Printf("File: %s\n", relPath)
		fmt.Printf("  Basename: %s\n", baseName)

		// Manually check which rule would match
		fileType := classifyFile(mdParser, relPath)
		fmt.Printf("  Classified as: %s\n\n", fileType)
	}

	// Now check what's actually indexed
	fmt.Println("\n=== Checking Index Contents ===")
	a, err := analyzer.New(context.Background(), testDir)
	if err != nil {
		fmt.Printf("Failed to create analyzer: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Note: Index may be empty if not yet indexed.")
	fmt.Println("Run test_memory_system to index first if needed.")
}

// Simulate the classification logic
func classifyFile(p *parser.Parser, filePath string) parser.FileType {
	// This is a hacky way since we can't call private method
	// We'll try to parse and check the chunk types
	file, err := p.Chunk(filePath)
	if err != nil {
		return parser.FileTypeIgnore
	}

	if len(file.Chunks) > 0 {
		return parser.FileType(file.Chunks[0].Type)
	}

	return parser.FileTypeSrc
}
