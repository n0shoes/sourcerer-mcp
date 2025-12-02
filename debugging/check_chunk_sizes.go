package main

import (
	"fmt"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

func main() {
	testDir := "/Users/claudeuser/claude_home/test_project"
	files := []string{"CLAUDE.md", "decisions.md", "http_request.py"}

	fmt.Println("=== Chunk Size Analysis ===\n")

	for _, fileName := range files {
		var p *parser.Parser
		var err error

		if fileName[len(fileName)-3:] == ".md" {
			p, err = parser.NewMarkdownParser(testDir)
		} else if fileName[len(fileName)-3:] == ".py" {
			p, err = parser.NewPythonParser(testDir)
		}

		if err != nil {
			fmt.Printf("Failed to create parser for %s: %v\n", fileName, err)
			continue
		}

		file, err := p.Chunk(fileName)
		if err != nil {
			fmt.Printf("Failed to parse %s: %v\n", fileName, err)
			p.Close()
			continue
		}

		fmt.Printf("File: %s\n", fileName)
		fmt.Printf("  Total chunks: %d\n", len(file.Chunks))

		for i, chunk := range file.Chunks {
			chunkSize := len(chunk.Source)
			fmt.Printf("  [%d] %-8s %6d bytes - %s\n", i+1, chunk.Type, chunkSize, chunk.ID())

			if chunkSize > 2000 {
				fmt.Printf("       ⚠️  Large chunk! First 100 chars: %s...\n", chunk.Source[:100])
			}
		}
		fmt.Println()

		p.Close()
	}
}
