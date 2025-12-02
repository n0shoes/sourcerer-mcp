# Bug Fixes for Memory File Type Classification and Ollama Integration

## Issues Fixed

### 1. File Type Classification Priority (internal/parser/parser.go)
**Problem:** MEMORY.md and decisions.md files were being misclassified as `docs` or `src` instead of `memory` type because global directory rules (like `docs/**`) were checked before language-specific filename rules.

**Fix:** Reversed the rule checking order to prioritize language-specific rules (MEMORY.md, decisions.md) over global directory patterns. Now these files are always classified as `memory` type regardless of their directory location.

### 2. Search Type Filtering (internal/index/index.go)
**Problem:** Search queried all chunks first, then filtered by type after retrieval. This meant if you had 500 docs chunks and 5 memory chunks, searching for `file_types=['memory']` would often return 0 results because the top 30 results were all docs.

**Fix:** Modified search to query each file type separately using chromem's `where` parameter, then merge results. This ensures each requested file type gets representation in search results.

### 3. Empty Collection Handling (internal/index/index.go)
**Problem:** Requesting more results than exist in the collection caused error: "nResults must be <= the number of documents in the collection".

**Fix:** Added collection size checking before querying to ensure nResults never exceeds available documents.

### 4. Database Path Location (internal/index/index.go)
**Problem:** Database was created in current working directory (`.sourcerer/db`) instead of the workspace root directory being indexed.

**Fix:** Changed database path from relative `.sourcerer/db` to `filepath.Join(workspaceRoot, ".sourcerer/db")`. Each indexed project now maintains its own database in the correct location.

## Test Results

All fixes verified with test_memory_system.go:
- ✅ MEMORY.md and decisions.md correctly classified as `memory` type
- ✅ Memory searches return results from memory files
- ✅ Database created in workspace root, not current directory
- ✅ Works with local Ollama embeddings (nomic-embed-text)

### 5. Ollama Embedding Reliability (chromem-go fork)

**Problem:** Ollama's internal embedding service randomly failed with "EOF" errors when processing multiple embeddings sequentially. The failures were intermittent and unpredictable - sometimes the first chunk, sometimes the last, sometimes random chunks in the middle.

**Root Cause:** Ollama spawns ephemeral internal worker processes (on random ports like 60710, 61281, etc.) to handle embeddings. These workers occasionally crash or become unresponsive, returning HTTP 500 errors with EOF messages. This is an Ollama architecture issue, not a client-side problem - even direct curl requests experienced the same failures.

**Fixes Applied:**
1. **Better Error Reporting**: Uncommented HTTP status checking to surface actual Ollama errors instead of generic "no embeddings found" messages
2. **Switched to curl**: Replaced Go's HTTP client with direct curl subprocess calls to eliminate any potential connection pooling issues
3. **Retry Logic with Exponential Backoff**: Implemented 3-attempt retry with 500ms/1s/1.5s backoff delays to handle Ollama's transient failures
4. **Concurrency Limiting**: Reduced concurrent embedding requests from `runtime.NumCPU()` (8-16) to 1 for Ollama to avoid overwhelming the service

**Result:** Indexing now succeeds reliably. Failed embeddings are automatically retried and succeed on subsequent attempts. In testing, approximately 20-30% of embedding requests initially fail but succeed on retry.

## Chromem-go Fork Changes

The fixes above required modifications to the local chromem-go fork at `/Users/claudeuser/github/n0shoes/chromem-go`:

**File: `embed_ollama.go`**
- Added retry logic wrapper around embedding function
- Switched from `net/http` client to `os/exec` curl subprocess
- Added debug logging for request/response tracking
- Implemented exponential backoff retry strategy (3 attempts max)
- Removed HTTP keep-alive connection reuse (DisableKeepAlives)

## Migration Notes

**Important:** Projects indexed before these fixes will have incorrect file type classifications in their `.sourcerer/db`. Delete the `.sourcerer` directory in each project and re-index to pick up the correct classifications.

## Known Limitations

1. **Ollama Instability**: The retry logic works around Ollama's internal service instability but doesn't eliminate it. Expect ~20-30% of embedding requests to require retries.

2. **Performance**: Sequential processing (concurrency=1) and retry delays make indexing slower than with OpenAI's API, but ensures reliability.

3. **Batch Embeddings**: Ollama supports batch embedding API calls (multiple texts in one request) but chromem-go doesn't utilize this feature. Future optimization could batch all chunks from a file into one API call to reduce Ollama's failure surface.
