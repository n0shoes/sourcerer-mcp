package main

import (
	"context"
	"fmt"
	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
)

func main() {
	fmt.Println("Creating analyzer...")
	a, err := analyzer.New(context.Background(), "/Users/shaune/temp")
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Analyzer created successfully!")
	fmt.Println("Now triggering a search to force embedding creation...")
	
	// This will actually call the embedding function
	results, err := a.SemanticSearch(context.Background(), "test query", []string{"memory"})
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return
	}
	
	fmt.Printf("Search returned %d results\n", len(results))
}
