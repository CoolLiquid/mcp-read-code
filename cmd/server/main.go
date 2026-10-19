package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nanami7777777/gitlab-code-reader-mcp/internal/gitlab"
	"github.com/nanami7777777/gitlab-code-reader-mcp/internal/tools"
)

const instructions = `You have access to 9 read-only tools for exploring GitLab repository code. These tools NEVER modify anything — explore freely.

## RULES (read these first)

### What to use
- gl_read_file for reading files (NOT curl, NOT API calls)
- gl_find_files for finding files by name pattern (NOT gl_list_directory with grep)
- gl_search_code for searching code content (NOT reading files and scanning manually)
- gl_read_multiple for reading 2+ files (NOT multiple gl_read_file calls)
- gl_read_symbols for understanding large files (NOT reading the entire file)

### What NOT to do
- Do NOT read an entire large file to find one function — use gl_read_symbols first, then gl_read_file with start_line/end_line
- Do NOT call gl_read_file multiple times when gl_read_multiple can batch them
- Do NOT browse directories to find a file — use gl_find_files with a glob pattern
- Do NOT guess file paths — use gl_find_files or gl_list_directory to discover them first

### Parallel calls (MANDATORY)
When you need multiple independent pieces of information, call tools in parallel in the same message. Examples:
- Reading README.md AND listing directory → call both in one message
- Reading 3 specific files → use gl_read_multiple, not 3 separate calls

## TOOL SELECTION (strict priority order)

### First contact with a project
1. gl_list_directory(depth=2) → project structure
2. gl_read_file("README.md") → project purpose
3. gl_find_files with config patterns → tech stack

### Finding code
1. Know exact path → gl_read_file
2. Know filename pattern → gl_find_files → gl_read_file
3. Know content/keyword → gl_search_code
4. Large file (300+ lines) → gl_read_symbols first → gl_read_file with line range

### Code review
1. gl_diff(merge_request_iid) → overall changes
2. gl_read_multiple → full context of changed files (batch!)
3. gl_blame → history of changed areas
4. gl_search_code → find related references

### When to STOP searching
- Found the answer? Stop. Don't keep exploring "just in case."
- Got enough context? Stop. You don't need to read every related file.
- File is 300+ lines? Use gl_read_symbols, don't read the whole thing.
`

func defaultGitLabURL() string {
	if u := strings.TrimSpace(os.Getenv("GITLAB_URL")); u != "" {
		return u
	}
	return "https://gitlab.com"
}

func sanitizeToolArgs(args any) string {
	m, ok := args.(map[string]any)
	if !ok {
		b, _ := json.Marshal(args)
		return string(b)
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "token") || strings.Contains(lk, "secret") || strings.Contains(lk, "password") {
			out[k] = "[REDACTED]"
			continue
		}
		out[k] = v
	}
	b, _ := json.Marshal(out)
	return string(b)
}

func newMCPServer(pool *gitlab.ClientPool, stdioCreds *gitlab.Creds) *server.MCPServer {
	s := server.NewMCPServer(
		"gitlab-code-reader",
		"0.2.0",
		server.WithInstructions(instructions),
		server.WithToolHandlerMiddleware(func(next server.ToolHandlerFunc) server.ToolHandlerFunc {
			return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				if stdioCreds != nil {
					if _, ok := gitlab.CredsFromContext(ctx); !ok {
						ctx = gitlab.ContextWithCreds(ctx, *stdioCreds)
					}
				}
				start := time.Now()
				result, err := next(ctx, req)
				elapsed := time.Since(start)

				size := 0
				lines := 0
				if result != nil {
					for _, c := range result.Content {
						if tc, ok := c.(mcp.TextContent); ok {
							size += len(tc.Text)
							for _, ch := range tc.Text {
								if ch == '\n' {
									lines++
								}
							}
						}
					}
				}
				isErr := err != nil || (result != nil && result.IsError)
				fmt.Fprintf(os.Stderr, "[METRIC] tool=%-20s bytes=%-6d lines=%-4d time=%-8s error=%v\n",
					req.Params.Name, size, lines, elapsed.Round(time.Millisecond), isErr)
				return result, err
			}
		}),
		server.WithHooks(&server.Hooks{
			OnBeforeCallTool: []server.OnBeforeCallToolFunc{
				func(ctx context.Context, id any, req *mcp.CallToolRequest) {
					fmt.Fprintf(os.Stderr, "[CALL]   tool=%-20s args=%s\n",
						req.Params.Name, sanitizeToolArgs(req.Params.Arguments))
				},
			},
		}),
	)
	tools.RegisterAll(s, pool)
	return s
}

func runStdio(pool *gitlab.ClientPool, creds gitlab.Creds) error {
	s := newMCPServer(pool, &creds)
	fmt.Fprintln(os.Stderr, "gitlab-code-reader MCP server running on stdio")
	return server.ServeStdio(s)
}

func runHTTP(pool *gitlab.ClientPool, listenAddr, defaultURL string) error {
	s := newMCPServer(pool, nil)
	httpServer := server.NewStreamableHTTPServer(s,
		server.WithEndpointPath("/mcp"),
		server.WithHTTPContextFunc(func(ctx context.Context, r *http.Request) context.Context {
			return gitlab.ContextWithCreds(ctx, gitlab.Creds{
				BaseURL: gitLabURLFromRequest(r, defaultURL),
				Token:   extractGitLabToken(r),
			})
		}),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.Handle("/mcp", &requireGitLabToken{inner: httpServer, defaultURL: defaultURL})

	fmt.Fprintf(os.Stderr, "gitlab-code-reader MCP server listening on %s (HTTP /mcp)\n", listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}

func main() {
	pool := gitlab.NewClientPool()
	transport := strings.ToLower(strings.TrimSpace(os.Getenv("MCP_TRANSPORT")))
	if transport == "" {
		transport = "stdio"
	}

	defaultURL := defaultGitLabURL()

	switch transport {
	case "http":
		listen := os.Getenv("LISTEN_ADDR")
		if listen == "" {
			listen = ":8088"
		}
		if err := runHTTP(pool, listen, defaultURL); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
			os.Exit(1)
		}
	case "stdio":
		token := os.Getenv("GITLAB_TOKEN")
		if token == "" {
			fmt.Fprintln(os.Stderr, "Error: GITLAB_TOKEN environment variable is required for stdio transport.")
			os.Exit(1)
		}
		creds := gitlab.Creds{BaseURL: defaultURL, Token: token}
		if err := runStdio(pool, creds); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown MCP_TRANSPORT %q (use http or stdio)\n", transport)
		os.Exit(1)
	}
}
