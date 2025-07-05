package main

import (
	"fmt"
	"os"

	"github.com/ss49919201/ai-kata/mini-container/internal/container"
)

func main() {
	args := os.Args[1:]
	cmd, err := container.ParseArgs(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	switch cmd.Action {
	case "help":
		fmt.Println("mini-container v0.1.0")
		fmt.Println("Usage: container <command> [args...]")
		fmt.Println("Commands:")
		fmt.Println("  help    Show this help message")
		fmt.Println("  run     Run a command in a container")
	case "run":
		fmt.Printf("Running: %s %v\n", cmd.Program, cmd.Args)
	}
}