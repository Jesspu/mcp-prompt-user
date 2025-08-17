package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type promptResult struct {
	answer string
	err    error
}

func main() {
	mcpServer := server.NewMCPServer("mcp-prompt-user", "1.0.0", server.WithToolCapabilities(true))

	// Tools
	promptUserSchema := mcp.NewTool(
		"promptUser",
		mcp.WithDescription("Prompts the user with a question asking for more context."),
		mcp.WithString("prompt", mcp.Description("The question to prompt the user with."), mcp.Required()),
	)

	mcpServer.AddTool(promptUserSchema, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		log.Println("promptUser tool called")
		prompt, ok := req.GetArguments()["prompt"].(string)
		if !ok {
			log.Println("Error: missing prompt argument")
			return mcp.NewToolResultError("missing prompt argument"), nil
		}

		progressTokenValue := req.Params.Meta.ProgressToken
		if progressTokenValue == nil {
			// No progress token, so we'll just block and wait for the answer.
			log.Printf("Prompting user with: %s", prompt)
			answer, err := RunPrompt(prompt)
			if err != nil {
				log.Printf("Error from RunPrompt: %v", err)
				return mcp.NewToolResultError(err.Error()), nil
			}
			log.Printf("Received answer from user: %s", answer)
			return mcp.NewToolResultText(answer), nil
		}

		resultChan := make(chan promptResult)
		go func() {
			log.Printf("Prompting user with: %s", prompt)
			answer, err := RunPrompt(prompt)
			resultChan <- promptResult{answer, err}
		}()

		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		session := server.ClientSessionFromContext(ctx)
		if session == nil {
			return mcp.NewToolResultError("could not get client session from context"), nil
		}

		curProgress := 1
		for {
			select {
			case result := <-resultChan:
				if result.err != nil {
					log.Printf("Error from RunPrompt: %v", result.err)
					return mcp.NewToolResultError(result.err.Error()), nil
				}
				log.Printf("Received answer from user: %s", result.answer)
				return mcp.NewToolResultText(result.answer), nil
			case <-ticker.C:
				curProgress = curProgress + 1
				log.Println("Sending progress notification")
				mcpServer.SendNotificationToClient(ctx, "notifications/progress", map[string]any{
					"progressToken": progressTokenValue,
					"progress":      curProgress,
					"total":         100,
					"message":       "Waiting for user input...",
				})
			}
		}
	})

	// Resources
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

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
