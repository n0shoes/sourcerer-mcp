package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
)

func main() {
	// Use a test workspace path
	testWorkspace := "/Users/claudeuser/claude_home/test_project"

	fmt.Printf("Testing DB path creation...\n")
	fmt.Printf("Workspace: %s\n\n", testWorkspace)

	// Check workspace exists
	if _, err := os.Stat(testWorkspace); os.IsNotExist(err) {
		fmt.Printf("❌ Test workspace doesn't exist: %s\n", testWorkspace)
		fmt.Printf("Creating test workspace...\n")
		os.MkdirAll(testWorkspace, 0755)

		// Create a simple test file
		testFile := filepath.Join(testWorkspace, "MEMORY.md")
		content := "# Project Memory\n\n## Decision: Use Go\nWe chose Go for performance."
		os.WriteFile(testFile, []byte(content), 0644)
		fmt.Printf("Created test file: %s\n\n", testFile)
	}

	// Create index - this should create .sourcerer/db in the workspace
	fmt.Println("Creating index (this will create .sourcerer/db)...")
	idx, err := index.New(context.Background(), testWorkspace)
	if err != nil {
		fmt.Printf("❌ Failed to create index: %v\n", err)
		os.Exit(1)
	}
	_ = idx

	// Check where the DB was created
	expectedDBPath := filepath.Join(testWorkspace, ".sourcerer/db")
	wrongDBPath := ".sourcerer/db" // relative to current dir

	if _, err := os.Stat(expectedDBPath); err == nil {
		fmt.Printf("✅ SUCCESS: DB created in workspace at:\n   %s\n", expectedDBPath)
	} else {
		fmt.Printf("❌ FAILED: DB NOT found in workspace at:\n   %s\n", expectedDBPath)
	}

	if _, err := os.Stat(wrongDBPath); err == nil {
		fmt.Printf("⚠️  WARNING: DB also created in current directory at:\n   %s\n", wrongDBPath)
		fmt.Printf("   (This should NOT happen with the fix)\n")
	}

	fmt.Println("\n=== Test Complete ===")
}
