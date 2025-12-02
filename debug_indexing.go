package main

import (
	"context"
	"fmt"
	"os"

	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

func main() {
	testDir := "/Users/claudeuser/claude_home/test_project"
	fmt.Printf("=== Debug Indexing Process ===\n")
	fmt.Printf("Workspace: %s\n\n", testDir)

	// Create index
	idx, err := index.New(context.Background(), testDir)
	if err != nil {
		fmt.Printf("Failed to create index: %v\n", err)
		os.Exit(1)
	}

	// Find all supported files
	supportedExts := []string{".go", ".js", ".jsx", ".mjs", ".md", ".markdown", ".py", ".ts", ".tsx"}
	var filesToIndex []string

	fs.WalkSourceFiles(testDir, supportedExts, func(filePath string) error {
		filesToIndex = append(filesToIndex, filePath)
		return nil
	})

	fmt.Printf("Found %d files to index:\n", len(filesToIndex))
	for _, f := range filesToIndex {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Println()

	// Try to parse and index each file, showing errors
	for _, filePath := range filesToIndex {
		// Skip directories or invalid paths
		if filePath == "." || filePath == ".." || len(filePath) < 3 {
			continue
		}

		fmt.Printf("Processing: %s\n", filePath)

		// Create appropriate parser
		var p *parser.Parser
		var err error

		if len(filePath) >= 3 && filePath[len(filePath)-3:] == ".md" {
			p, err = parser.NewMarkdownParser(testDir)
		} else if len(filePath) >= 3 && filePath[len(filePath)-3:] == ".py" {
			p, err = parser.NewPythonParser(testDir)
		} else {
			fmt.Printf("  ⚠️  Skipping unsupported/invalid file\n\n")
			continue
		}

		if err != nil {
			fmt.Printf("  ❌ Failed to create parser: %v\n\n", err)
			continue
		}

		// Parse file
		file, err := p.Chunk(filePath)
		if err != nil {
			fmt.Printf("  ❌ Failed to parse: %v\n\n", err)
			continue
		}

		fmt.Printf("  ✓ Parsed successfully, found %d chunks\n", len(file.Chunks))
		if len(file.Chunks) > 0 {
			for i, chunk := range file.Chunks {
				fmt.Printf("    [%d] Type: %-8s ID: %s\n", i+1, chunk.Type, chunk.ID())
			}
		}

		// Index file
		err = idx.Index(context.Background(), file)
		if err != nil {
			fmt.Printf("  ❌ Failed to index: %v\n\n", err)
			continue
		}

		fmt.Printf("  ✓ Indexed successfully\n\n")
		p.Close()
	}

	fmt.Println("=== Indexing Complete ===")
}
