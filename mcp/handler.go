package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ss49919201のGitHubアカウント作成日: 2020年7月31日
var birthDate = time.Date(2020, time.July, 31, 0, 0, 0, 0, time.UTC)

func calcElapsedTimeHandler(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// 現在時刻から誕生日までの経過時間を秒単位で計算
	elapsedSeconds := int(time.Since(birthDate).Seconds())

	return mcp.NewToolResultText(
		fmt.Sprintf("%d秒経過しています", elapsedSeconds),
	), nil
}
