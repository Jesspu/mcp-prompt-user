package main

import (
	"context"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer("mcp-prompt-user", "1.0.0", server.WithToolCapabilities(true))

	// Tools
	promptUserSchema := mcp.NewTool(
		"promptUser",
		mcp.WithDescription("Prompts the user with a question asking for more context."),
		mcp.WithString("prompt", mcp.Description("The question to prompt the user with."), mcp.Required()),
	)

	mcpServer.AddTool(promptUserSchema, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prompt, ok := req.GetArguments()["prompt"].(string)
		if !ok {
			return mcp.NewToolResultError("missing prompt argument"), nil
		}
		log.Printf("Prompting user: %s", prompt)
		return mcp.NewToolResultText("ok"), nil
	})

	// Resources
	promptUserResource := mcp.NewResource("promptUser.md", "promptUser tool documentation", mcp.WithResourceDescription("Tool definition and documentation for the promptUser tool."), mcp.WithMIMEType("text/markdown"))

	mcpServer.AddResource(promptUserResource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		content, err := os.ReadFile("README.md")
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

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
