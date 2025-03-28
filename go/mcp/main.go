package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"ss49919201 profile",
		"1.0.0",
	)

	s.AddTool(
		mcp.NewTool("elapsed_time_since_ss49919201_was_born",
			mcp.WithDescription("Return the elapsed time since ss49919201 was born in seconds"),
		),
		calcElapsedTimeHandler,
	)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// ss49919201のGitHubアカウント作成日: 2020年7月31日
var birthDate = time.Date(2020, time.July, 31, 0, 0, 0, 0, time.UTC)

func calcElapsedTimeHandler(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 現在時刻から誕生日までの経過時間を秒単位で計算
	elapsedSeconds := int(time.Since(birthDate).Seconds())

	return mcp.NewToolResultText(
		fmt.Sprintf("%d秒経過しています", elapsedSeconds),
	), nil
}
