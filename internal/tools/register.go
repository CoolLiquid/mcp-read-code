package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/nanami7777777/gitlab-code-reader-mcp/internal/gitlab"
)

// RegisterAll adds the nine GitLab code reader tools to the MCP server.
func RegisterAll(s *server.MCPServer, pool *gitlab.ClientPool) {
	t, h := ReadFile(pool)
	s.AddTool(t, h)
	t, h = ReadMultiple(pool)
	s.AddTool(t, h)
	t, h = FindFiles(pool)
	s.AddTool(t, h)
	t, h = SearchCode(pool)
	s.AddTool(t, h)
	t, h = ListDirectory(pool)
	s.AddTool(t, h)
	t, h = ReadSymbols(pool)
	s.AddTool(t, h)
	t, h = Diff(pool)
	s.AddTool(t, h)
	t, h = Blame(pool)
	s.AddTool(t, h)
	t, h = CommitHistory(pool)
	s.AddTool(t, h)
}
