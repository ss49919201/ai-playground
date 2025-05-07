package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"ss49919201 profile",
		"1.0.0",
	)

	fmt.Println("Server created with name: ss49919201 profile, version: 1.0.0")

	s.AddTool(
		mcp.NewTool("elapsed_time_since_ss49919201_was_born",
			mcp.WithDescription("Return the elapsed time since ss49919201 was born in seconds"),
		),
		calcElapsedTimeHandler,
	)

	fmt.Println("Tool added: elapsed_time_since_ss49919201_was_born")

	// Start the stdio server
	fmt.Println("Start Server ...")
	fmt.Println("Starting stdio server...")

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}

	fmt.Println("Server stopped")
}
