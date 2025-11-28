package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
)

type Server struct {
	workspaceRoot string
	mcp           *server.MCPServer
	analyzer      *analyzer.Analyzer
}

func NewServer(workspaceRoot, version string) (*Server, error) {
	a, err := analyzer.New(context.Background(), workspaceRoot)
	if err != nil {
		return nil, err
	}

	s := &Server{
		workspaceRoot: workspaceRoot,
		analyzer:      a,
	}

	s.mcp = server.NewMCPServer(
		"Sourcerer",
		version,
		server.WithInstructions(`
You have access to Sourcerer MCP tools for efficient codebase navigation AND
project memory access. Sourcerer provides surgical precision - you can jump
directly to specific functions, classes, and code chunks without reading entire
files, AND recall past decisions, design rationale, and documented context.
This reduces token waste & cognitive load.

WHAT SOURCERER INDEXES:
Sourcerer maintains a semantic index of your entire project:
- Source code: Functions, classes, methods, types, implementations
- Documentation: General markdown files (README, docs/, etc.)
- Project memory: MEMORY.md, decisions.md - architectural decisions, design
  rationale, constraints, lessons learned, and project context
- Tests: Test suites and test implementations

WHEN TO USE SOURCERER:
1. Finding code implementations and understanding how things work
2. Recalling past decisions: "What did we decide about X?"
3. Understanding design rationale: "Why did we choose this approach?"
4. Accessing documented constraints, patterns, and architectural context
5. Finding related code and documentation together

MEMORY & CONTEXT QUERIES:
Use semantic search to find project decisions, design rationale, and documented
context from MEMORY.md and decisions.md files. Sourcerer will return relevant
sections that answer questions about past choices and project evolution.

Example memory queries:
✓ "What did we decide about authentication?"
✓ "Why did we choose this database?"
✓ "What are the known limitations of the API?"
✓ "What's our approach to error handling?"
✓ "What were the tradeoffs we considered for X?"
✓ "What dependencies should we avoid and why?"

SEARCH STRATEGY:
Sourcerer's semantic search understands concepts and relationships:
- Describe the purpose/behavior you're seeking
- Use natural language to explain the concept
- Include context about what the code should accomplish
- Mention related functionality or typical patterns

The line numbers shown in search results (e.g., "lines 45-67") reference the
exact location in the original file and can be used with standard file tools
if you need to read or edit those specific sections.

Use the file_types param to filter search results (defaults to ['src', 'docs']):
- src: Source code implementations
- docs: General documentation (README, guides, API docs)
- memory: Project memory (MEMORY.md, decisions.md) - architectural decisions
  and design rationale
- tests: Test code

FILE TYPE FILTERING EXAMPLES:
- For code implementation only:
  file_types: ['src']
  Query: "rate limiting implementation"

- For decisions and context only:
  file_types: ['memory']
  Query: "What did we decide about rate limiting?"

- For both code AND decisions (comprehensive):
  file_types: ['src', 'memory']
  Query: "rate limiting" → finds both implementation AND decision rationale

- Search everything (default behavior):
  file_types: ['src', 'docs', 'memory']
  or omit file_types parameter

COMBINED CODE + MEMORY QUERIES:
When you want to understand both the "what" (implementation) and "why"
(decision), search across both src and memory:

✓ "authentication" with file_types: ['src', 'memory']
  → finds auth code AND the decision to use JWT tokens

✓ "database" with file_types: ['src', 'memory']
  → finds DB code AND the rationale for choosing PostgreSQL

✓ "error handling" with file_types: ['src', 'memory']
  → finds error handling code AND documented patterns/conventions

AVOID SEMANTIC SEARCH FOR EXACT MATCHES:
If you need to find specific names or exact text, use pattern-based tools
like grep & glob instead:

Good: "authentication logic and session management"
Avoid: "AuthService class definition" (use grep instead)

CHUNK IDs:
Use chunk IDs to retrieve source code:
- Type definition: path/to/file.ext::Type
- Specific method in Type: path/to/file.ext::Type::method
- Variable: path/to/file.ext::Var
- Content-based chunks: file.ext::695fffd41945e08d (imports, markdown sections)

Chunk IDs are stable across minor edits but update when code structure
changes (renames, moves, deletions). Use get_chunk_code with these precise
ids to get exactly the code you need.

If you already know the specific function/class/method/struct/etc and file
location from previous context, construct the chunk ID yourself and use
get_chunk_code directly rather than semantic searching again.

MARKDOWN CHUNKS:
Markdown files are chunked by section (## headers). Each section becomes a
searchable chunk. For example:
- MEMORY.md::Authentication Decision
- decisions.md::Database Selection Rationale

BATCHING:
Batch operations instead of making separate requests which waste tokens and
time (round-trips).

DO NOT try pulling all chunks within a specific file (with an id like file.ext).
That defeats the purpose of surgical precision. If you need the entire file,
just read it directly with your standard tools.

WORKFLOW EXAMPLES:

Starting work on authentication:
1. Search memory: "authentication approach" with file_types: ['memory']
   → Understand past decisions and constraints
2. Search code: "authentication" with file_types: ['src']
   → Find current implementation
3. Use get_chunk_code to examine specific functions

Understanding an unfamiliar module:
1. Search broadly: "payment processing" with file_types: ['src', 'memory', 'docs']
   → Get code, decisions, and documentation
2. Use chunk IDs to dive deeper into relevant pieces

Recalling project context:
1. Search memory: "why did we choose Redis?" with file_types: ['memory']
   → Find decision rationale and tradeoffs considered
`),
	)

	s.mcp.AddTool(
		mcp.NewTool("semantic_search",
			mcp.WithDescription("Find relevant code using semantic search"),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Your search"),
			),
			mcp.WithArray("file_types",
				mcp.WithStringItems(),
				mcp.Description("Filter by file type(s)"),
			),
		),
		s.semanticSearch,
	)

	s.mcp.AddTool(
		mcp.NewTool("find_similar_chunks",
			mcp.WithDescription("Find code chunks semantically similar to a given chunk"),
			mcp.WithString("id",
				mcp.Required(),
				mcp.Description("The chunk ID to find similar code for"),
			),
		),
		s.findSimilarChunks,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_chunk_code",
			mcp.WithDescription("Get the actual code you need to examine"),
			mcp.WithArray("ids",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description("Chunks to get code for"),
			),
		),
		s.getChunkCode,
	)

	s.mcp.AddTool(
		mcp.NewTool("index_workspace",
			mcp.WithDescription("Index all pending files in the workspace"),
		),
		s.indexWorkspace,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_index_status",
			mcp.WithDescription("Get the codebase's indexing status"),
		),
		s.getIndexStatus,
	)

	//new memory specific tool. (SE)
	s.mcp.AddTool(
		mcp.NewTool("search_memory",
			mcp.WithDescription("Search project memory for past decisions, design rationale, and documented context. Searches MEMORY.md and decisions.md files."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("What decision, rationale, or context to find (e.g., 'authentication approach', 'why we chose PostgreSQL')"),
			),
		),
		s.searchMemory,
	)

	return s, nil
}

func (s *Server) Serve() error {
	return server.ServeStdio(s.mcp)
}

func (s *Server) semanticSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")
	fileTypes := request.GetStringSlice("file_types", []string{"src", "docs"})

	results, err := s.analyzer.SemanticSearch(ctx, query, fileTypes)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No matching chunks found."), nil
	}

	content := strings.Join(results, "\n")
	return mcp.NewToolResultText(content), nil
}

func (s *Server) findSimilarChunks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	chunkID := request.GetString("id", "")

	results, err := s.analyzer.FindSimilarChunks(ctx, chunkID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No similar chunks found."), nil
	}

	content := strings.Join(results, "\n")
	return mcp.NewToolResultText(content), nil
}

func (s *Server) getChunkCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ids := request.GetStringSlice("ids", []string{})

	chunks := s.analyzer.GetChunkCode(ctx, ids)

	return mcp.NewToolResultText(chunks), nil
}

func (s *Server) indexWorkspace(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	go s.analyzer.IndexWorkspace(ctx)

	return mcp.NewToolResultText("Indexing in progress..."), nil
}

func (s *Server) getIndexStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pendingFiles, lastIndexedAt := s.analyzer.GetIndexStatus()

	status := fmt.Sprintf("Number of pending files: %d, last indexed: ", pendingFiles)
	if lastIndexedAt.IsZero() {
		status += "in progress"
	} else {
		status += humanize.Time(lastIndexedAt)
	}

	return mcp.NewToolResultText(status), nil
}

//new memory search function (SE)
func (s *Server) searchMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")

	// Search only memory file type
	results, err := s.analyzer.SemanticSearch(ctx, query, []string{"memory"})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Memory search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No matching decisions or context found in project memory."), nil
	}

	content := strings.Join(results, "\n")
	return mcp.NewToolResultText(content), nil
}

func (s *Server) Close() error {
	if s.analyzer != nil {
		s.analyzer.Close()
	}

	return nil
}
