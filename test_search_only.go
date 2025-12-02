package main

import (
	"context"
	"fmt"
	"os"

	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
)

func main() {
	testDir := "/Users/claudeuser/claude_home/test_project"
	fmt.Println("=== Testing Search on Existing Index ===\n")

	// Create analyzer (connects to existing database)
	a, err := analyzer.New(context.Background(), testDir)
	if err != nil {
		fmt.Printf("Failed to create analyzer: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("1. Searching for 'python' in memory files...")
	memResults, err := a.SemanticSearch(context.Background(), "python", []string{"memory"})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d results:\n", len(memResults))
		for i, r := range memResults {
			fmt.Printf("     [%d] %s\n", i+1, r)
		}
	}

	fmt.Println("\n2. Searching for 'fuzzer' in all types...")
	allResults, err := a.SemanticSearch(context.Background(), "fuzzer", []string{"src", "docs", "memory"})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d results:\n", len(allResults))
		for i, r := range allResults {
			fmt.Printf("     [%d] %s\n", i+1, r)
		}
	}

	fmt.Println("\n3. Searching for 'function' in src files...")
	srcResults, err := a.SemanticSearch(context.Background(), "function", []string{"src"})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   Found %d results:\n", len(srcResults))
		for i, r := range srcResults {
			fmt.Printf("     [%d] %s\n", i+1, r)
		}
	}

	fmt.Println("\n=== Search Test Complete ===")
}
