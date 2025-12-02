package main

import (
	"context"
	"fmt"
	"os"

	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
)

func main() {
	testDir := "/Users/claudeuser/claude_home/test_project"

	fmt.Println("=== Checking Index Contents ===\n")

	a, err := analyzer.New(context.Background(), testDir)
	if err != nil {
		fmt.Printf("Failed to create analyzer: %v\n", err)
		os.Exit(1)
	}

	// Search for each file type to see what's indexed
	types := []string{"memory", "docs", "src", "tests"}

	for _, fileType := range types {
		results, err := a.SemanticSearch(context.Background(), "*", []string{fileType})
		if err != nil {
			fmt.Printf("Error searching type %s: %v\n", fileType, err)
			continue
		}

		fmt.Printf("Type: %-10s → %d chunks\n", fileType, len(results))
		if len(results) > 0 {
			for i, result := range results {
				if i >= 5 {
					fmt.Printf("  ... and %d more\n", len(results)-5)
					break
				}
				fmt.Printf("  %s\n", result)
			}
		}
		fmt.Println()
	}
}
