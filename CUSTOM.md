
This fork is for using with local ollama (no openai) and focuses on creating "memory" for claude code.


MCP server instructions have been updated from the original project to focus on the memory use-case.

You can ask claude code to use memory to find xyz about your project after you have indexed the project using the mcp server function. I'll look at the md files and it'll look for specific files as well MEMORY.md and decisions.md.

Put the mcp config below in your project's directory and it creates a sourcerer db in SOURCERER_WORKSPACE_ROOT.
.mcp.json
{
   "mcpServers": {
      "sourcerer": {
         "command": "/Users/xxxx/github/xxxx/sourcerer-mcp/sourcerer",
         "env": {
            "SOURCERER_WORKSPACE_ROOT": "/Users/xxxx/project_dir"
         }
      }
   }
}


Make sure you update go.mod to point to your own chromem-go project. (has a fix for ollama embeddings)
This is at the bottom of the file, change this path to point to your local chromem-go project.
Grab that from the n0shoes repository on github.

//Checkout the version and use the local version - change this to your own local path.
replace github.com/philippgille/chromem-go => /Users/xxxx/github/xxxx/chromem-go


Changes

Just for reference, the key files which have changed in the latest build are:

  sourcerer-mcp:
  - internal/index/index.go - Search fixes, concurrency, DB path
  - internal/parser/parser.go - File type classification priority
  - internal/parser/markdown.go - MEMORY.md/decisions.md rules
  - FIXES.md - Documentation of all fixes
  - go.mod - Points to local chromem-go fork

  chromem-go:
  - embed_ollama.go - Retry logic, curl implementation, error handling
  The test files (debug_*.go, test_*.go, etc.) were just for debugging - you can keep or delete those as needed.


I think the only debugging files needed are the two that are left in the project dir. 
The files in the debugging directory are temp files created by CC. But i think the 
debug_indexing.go and test_search_only.go will be enough to show if it's working or not with 
your project documents. You will need to ensure there is a project path that matches the one 
defined in that file though if you try to run them.
Use that and the build-sourcerer.sh file. 

If you want to re-index your project make sure you remove the relevant .sourcerer directory before
re-indexing. This should live in the project directory eg. test_project/.sourcerer

My test_project has the following files:
- CLAUDE.md 
- decisions.md (mentions some information about the project like why python was chosen and some date information)
- http_request.py (just a python program with some functions that can be indexed)

after you have setup your mcp.json file (described above for ollama) you should be able to ask it to:
"Use memory and tell me why we chose python for the project"
and information should be pulled from ollama rather than searched for on disk.


