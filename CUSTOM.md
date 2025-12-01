
This fork is for using with local ollama (no openai) and focuses on creating "memory" for claude code.
Note: Claude code is a lazy fucking cunt and will still probably ignore this and do fuck all. 


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
