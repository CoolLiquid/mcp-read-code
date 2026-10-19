package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/nanami7777777/gitlab-code-reader-mcp/internal/gitlab"
)

func clientFromContext(ctx context.Context, pool *gitlab.ClientPool) (*gitlab.Client, *mcp.CallToolResult, error) {
	client, err := pool.ClientFromContext(ctx)
	if err != nil {
		return nil, mcp.NewToolResultError(fmt.Sprintf("❌ %v", err)), nil
	}
	return client, nil, nil
}
