package main

import (
	"log"

	"github.com/Jesspu/mcp-prompt-user/resources"
	"github.com/Jesspu/mcp-prompt-user/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := server.NewMCPServer("mcp-prompt-user", "1.0.0", server.WithToolCapabilities(true))

	// Tools
	tools.RegisterPromptUserTool(mcpServer)

	// Resources
	resources.RegisterPromptUserDocResource(mcpServer, "promptUser.md")

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
