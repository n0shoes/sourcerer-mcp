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
persistent project memory. Sourcerer acts as your EXTENDED MEMORY that survives
conversation compaction - preserving institutional knowledge, past decisions,
and project context across sessions. It also provides surgical precision for
code navigation without reading entire files.

CRITICAL CONCEPT - MEMORY PERSISTENCE:
Sourcerer's memory search is your safety net against context loss. When
conversations get compacted or reset, project memory (MEMORY.md, decisions.md)
remains searchable. Think of it as institutional knowledge that persists beyond
your working memory. ALWAYS check memory before making architectural decisions
or significant changes.

WHAT SOURCERER INDEXES:
Sourcerer maintains a semantic index of your entire project:
- Project memory: MEMORY.md, decisions.md - architectural decisions, design
  rationale, constraints, lessons learned, and project context (PERSISTS ACROSS
  COMPACTION)
- Source code: Functions, classes, methods, types, implementations
- Documentation: General markdown files (README, docs/, etc.)
- Tests: Test suites and test implementations

WHEN TO USE SOURCERER (MEMORY-FIRST APPROACH):

PRIMARY USE - Project Memory & Context:
1. BEFORE implementing new features: "What decisions exist about X?"
2. BEFORE refactoring: "Why was this built this way?"
3. BEFORE choosing libraries/tools: "What did we already decide to use/avoid?"
4. Understanding constraints: "What limitations should I know about?"
5. Learning from history: "What did we try before? What failed? Why?"
6. Accessing project patterns: "What's our approach to error handling?"
7. Timeline questions: "When was this used? How long is this expected to last?"

SECONDARY USE - Code Navigation:
8. Finding code implementations when you need surgical precision
9. Understanding how things work without reading entire files
10. Locating specific functions, classes, or modules across the codebase

MEMORY-FIRST WORKFLOW:
When starting ANY implementation task, follow this pattern:

1. CHECK MEMORY FIRST - Search for relevant past decisions
   - Use search_memory or semantic_search with file_types: ['memory']
   - Ask: "What do we already know about this?"
   - Ask: "What did we decide before?"

2. UNDERSTAND THE CODE - Then explore current implementation
   - Use semantic_search with file_types: ['src']
   - Use chunk IDs to examine specific functions

3. IMPLEMENT - Armed with institutional knowledge + current state
   - You now know WHY things are the way they are
   - You can avoid repeating past mistakes
   - You respect established patterns and constraints

Example memory-first workflow:
User: "Add caching to the API"
✓ FIRST: Search memory: "caching decisions performance" with file_types: ['memory']
  → Discover we already evaluated Redis vs in-memory, chose Redis, documented why
✓ THEN: Search code: "API endpoints" with file_types: ['src']
  → Find current API implementation
✓ FINALLY: Implement Redis caching respecting past decisions

MEMORY & CONTEXT QUERIES:
Use search_memory (dedicated memory tool) or semantic_search with file_types:
['memory'] to find project decisions, design rationale, and documented context
from MEMORY.md and decisions.md files. This returns relevant sections about
past choices and project evolution.

Example memory queries:
✓ "What did we decide about authentication?"
✓ "Why did we choose this database?"
✓ "What are the known limitations of the API?"
✓ "What's our approach to error handling?"
✓ "What were the tradeoffs we considered for X?"
✓ "What dependencies should we avoid and why?"
✓ "When was this project used?"
✓ "What patterns should I follow?"
✓ "What did we learn from the previous implementation?"
✓ "What constraints exist for this feature?"

TRIGGERS FOR MEMORY SEARCH:
Automatically check memory when you encounter:
- "Add [feature]" → Check if we decided on an approach
- "Refactor [component]" → Check why it was built that way
- "Fix [bug]" → Check if this was addressed before
- "Improve [performance]" → Check what we already tried
- "Use [library/tool]" → Check if we evaluated it already
- "Change [architecture]" → Check past architectural decisions
- User mentions previous work or decisions
- You're about to make a significant technical choice

SEARCH STRATEGY:
Sourcerer's semantic search understands concepts and relationships:
- Describe the purpose/behavior you're seeking
- Use natural language to explain the concept
- Include context about what the code should accomplish
- Mention related functionality or typical patterns

The line numbers shown in search results (e.g., "lines 45-67") reference the
exact location in the original file and can be used with standard file tools
if you need to read or edit those specific sections.

FILE TYPE FILTERING:
Use the file_types param to filter search results (defaults to ['src', 'docs']):
- memory: Project memory (MEMORY.md, decisions.md) - USE THIS FIRST
- src: Source code implementations
- docs: General documentation (README, guides, API docs)
- tests: Test code

FILE TYPE FILTERING EXAMPLES:

Memory-first pattern (RECOMMENDED):
1. file_types: ['memory']
   Query: "What did we decide about rate limiting?"
   → Get institutional knowledge first

2. file_types: ['src']
   Query: "rate limiting implementation"
   → Then find current code

Combined comprehensive search:
- file_types: ['src', 'memory']
  Query: "rate limiting"
  → Finds both implementation AND decision rationale in one search

- file_types: ['src', 'memory', 'docs']
  Query: "authentication"
  → Complete picture: code + decisions + documentation

Code-only (when you already checked memory):
- file_types: ['src']
  Query: "rate limiting implementation"

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
Use chunk IDs to retrieve source code with surgical precision:
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
- MEMORY.md::Lessons Learned
- decisions.md::Why We Chose PostgreSQL

BATCHING:
Batch operations instead of making separate requests which waste tokens and
time (round-trips).

DO NOT try pulling all chunks within a specific file (with an id like file.ext).
That defeats the purpose of surgical precision. If you need the entire file,
just read it directly with your standard tools.

WORKFLOW EXAMPLES:

Example 1: Starting work on a new feature (MEMORY-FIRST):
User: "Add user authentication to the app"
1. Search memory FIRST: "authentication approach decisions" with file_types: ['memory']
   → Discover we already decided on JWT, documented security requirements
2. Search code: "authentication" with file_types: ['src']
   → Find any existing auth-related code
3. Search docs: "authentication" with file_types: ['docs']
   → Check if there's API documentation
4. Implement with full context of past decisions and current state

Example 2: Understanding an unfamiliar module:
1. Search memory: "payment processing decisions" with file_types: ['memory']
   → Understand why we built it this way, what we considered
2. Search broadly: "payment processing" with file_types: ['src', 'docs']
   → Get code and documentation
3. Use chunk IDs to dive deeper into relevant pieces

Example 3: Refactoring existing code:
User: "Refactor the database layer"
1. Search memory FIRST: "database architecture decisions" with file_types: ['memory']
   → Learn why it was designed this way, what constraints exist
2. Search memory: "database lessons learned problems" with file_types: ['memory']
   → Avoid repeating past mistakes
3. Search code: "database layer" with file_types: ['src']
   → Examine current implementation
4. Refactor respecting documented constraints and lessons

Example 4: Bug fix with historical context:
User: "The API is timing out"
1. Search memory: "API performance timeout" with file_types: ['memory']
   → Check if this happened before, what fixed it
2. Search memory: "API limitations constraints" with file_types: ['memory']
   → Understand known issues
3. Search code: "API timeout handling" with file_types: ['src']
   → Find current implementation
4. Fix with knowledge of past solutions

Example 5: After conversation compaction:
User: "Continue working on the authentication feature"
1. Search memory: "authentication" with file_types: ['memory']
   → Recover context about what was decided, even if conversation was compacted
2. Search code: "authentication" with file_types: ['src']
   → Find current implementation state
3. Continue work with full context restored

WHY MEMORY-FIRST MATTERS:
- Prevents re-litigating decisions already made
- Avoids repeating past mistakes documented in lessons learned
- Respects established patterns and constraints
- Provides continuity across conversation resets
- Saves time by leveraging institutional knowledge
- Ensures consistency with project direction
- Helps you understand the "why" not just the "what"

MEMORY AS YOUR SAFETY NET:
When you're unsure, check memory. When conversations reset, check memory. When
starting new work, check memory. The project memory is specifically designed to
survive context loss and provide you with the institutional knowledge needed to
make informed decisions even after compaction events.

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
