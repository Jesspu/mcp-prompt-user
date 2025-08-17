package tools

import (
	"context"
	"time"

	"github.com/Jesspu/mcp-prompt-user/tui"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type promptResult struct {
	answer string
	err    error
}

func RegisterPromptUserTool(mcpServer *server.MCPServer) {
	promptUserSchema := mcp.NewTool(
		"promptUser",
		mcp.WithDescription("Prompts the user with a question asking for more context."),
		mcp.WithString("prompt", mcp.Description("The question to prompt the user with."), mcp.Required()),
		mcp.WithDestructiveHintAnnotation(false),
		mcp.WithReadOnlyHintAnnotation(true),
	)

	mcpServer.AddTool(promptUserSchema, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		prompt, ok := req.GetArguments()["prompt"].(string)
		if !ok {
			return mcp.NewToolResultError("missing prompt argument"), nil
		}

		progressTokenValue := req.Params.Meta.ProgressToken
		if progressTokenValue == nil {
			// No progress token, so we'll just block and wait for the answer.
			answer, err := tui.RunPrompt(prompt)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText(answer), nil
		}

		resultChan := make(chan promptResult)
		go func() {
			answer, err := tui.RunPrompt(prompt)
			resultChan <- promptResult{answer, err}
		}()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		session := server.ClientSessionFromContext(ctx)
		if session == nil {
			return mcp.NewToolResultError("could not get client session from context"), nil
		}

		curProgress := 1
		totalProgress := 100
		for {
			select {
			case result := <-resultChan:
				if result.err != nil {
					return mcp.NewToolResultError(result.err.Error()), nil
				}
				return mcp.NewToolResultText(result.answer), nil
			case <-ticker.C:
				curProgress = curProgress + 1
				mcpServer.SendNotificationToClient(ctx, "notifications/progress", map[string]any{
					"progressToken": progressTokenValue,
					"progress":      curProgress,
					"total":         totalProgress,
					"message":       "Waiting for user input...",
				})
			}
		}
	})
}
