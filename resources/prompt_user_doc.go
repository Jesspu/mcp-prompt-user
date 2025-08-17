package resources

import (
	"context"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterPromptUserDocResource(mcpServer *server.MCPServer) {
	promptUserResource := mcp.NewResource("promptUser.md", "promptUser tool documentation", mcp.WithResourceDescription("Tool definition and documentation for the promptUser tool."), mcp.WithMIMEType("text/markdown"))

	mcpServer.AddResource(promptUserResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("promptUser.md")
		if err != nil {
			return nil, err
		}
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://readme",
				MIMEType: "text/markdown",
				Text:     string(content),
			},
		}, nil
	})
}
