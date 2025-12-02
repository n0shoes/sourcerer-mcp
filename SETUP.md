## Setup ollama 
Install ollama on macos
> ~/Applications/ollama

Install nomic-embed-text in ollama 
> using a terminal run:
> ollama pull nomic-embed-text

Check the endpoint works
> curl http://localhost:11434/api/embed -d '{ "model": "nomic-embed-text", "input": "test input" }'
> It will probably contain no data, but you should get something like: {"model":"nomic-embed-text","embeddings": ..

## Setup sourcerer and it'd main dependency chromem-go
Now clone the repo for sourcerer and for chromem-go  (these two are customised to work together)
> mkdir -p ~/github/n0shoes
> cd ~/github/n0shoes 
> git clone https://github.com/n0shoes/sourcerer-mcp.git
> git clone https://github.com/n0shoes/chromem-go.git

Setup the module path in sourcerer to use the local chromem-go
> cd ~/github/n0shoes/sourcerer-mcp
> Edit go.mod and change the last line to point to the chromem-go you just cloned
> The last line should read: where claudeuser is the account you're using. (i actually have an account for claude called claudeuser)
> replace github.com/philippgille/chromem-go => /Users/claudeuser/github/n0shoes/chromem-go

Build sourcerer using the pre-built shell script (take a look at it first so you know what it does - it's very short and simple!)
> From inside the dir ~/github/n0shoes/sourcerer-mcp
> ./build-sourcerer.sh
> That should create a file called sourcerer in the directory. Something like:
> -rwxrwxr-x   1 claudeuser  staff  16121426  2 Dec 16:13 sourcerer

## Setup a test project to try the MCP server.
Create a test project (or use the one provided) to test it out.
> cd test_project
> Create a file called .mcp.json in the test project with the following content (edit the path as needed)
> Note: Create one of these files in the root of each project you want to use sourcerer.
{
	"mcpServers": {
		"sourcerer": {
			"command": "/Users/claudeuser/github/n0shoes/sourcerer-mcp/sourcerer",
			"env": {
				"SOURCERER_WORKSPACE_ROOT": "/Users/claudeuser/github/n0shoes/sourcerer-mcp/test_project"
			}
		}
    }
}

## Run claude code to try the MCP server.
Change to the test_project directory and run claude code as usual
> use the claude code mcp command: /mcp to ensure sourcerer is connected.
> now ask claude something like: index the project with sourcerer
> then ask: check the status of the index
> then ask: why did we decide to use python for this project?
> then ask: when was this project used?
> Claude should use the sourcerer tooling to find the answer to the queries rather than only reading the files directly.
> Two files are key files for sourcerer at the moment: **MEMORY.md and decisions.md**.
> You should be able to put specifics in those files and have them indexed into "persistent" memory.
> Be careful not to let claude run rampant with those files and only add key information.
> Other md files will be indexed anyway, but these two files should take preference when searching using "memory". 
> The idea is that these memories will persist over compaction and help retain project knowledge over time. (hopefully)
> This will hopefully save tokens for Claude grepping around trying to find information after compaction.
> When md files or source files change it should automatically re-index the files.

It's important to understand this is very much alpha software. 
I will try to continue to test and maintain it, so keep an eye out for updates if you decide to use it.
The original sourcerer (this is a fork) was oriented toward searching source code files to reduce token use. 
This implementation has a different primary use case - persistent memory for claude code.

## Improvements
> The way you use MEMORY and decisions files will have an effect on how useful this tooling is.
> I'm thiking about and will add information for, creating good MEMORY and decision files as well as ways to add 
> temporal information, this might be instructions in CLAUDE.md and/or hooks or commands - unsure how this will work ATM.
> Experimenting with md files and headings containing dates etc might help formulate ways to ensure memory does not 
> get "stale" over time, that is, provide a way for new memories to be taken into consideration without too much manual effort.
>

## Troubleshooting
There are a couple of test scripts in the sourcer dir (one to test search and one to test the index)
to use these you need to ensure each script has the correct path to your test project
that is, you need to edit the script and look for the path - it should be obvious.
If things are failing to work as expected and you're trying to index your test project you should 
delete the .sourcer directory in the test project (to start with a fresh index) and try the test scripts 
again. Hopefully this project "just works" for you and you don't need to do this.



