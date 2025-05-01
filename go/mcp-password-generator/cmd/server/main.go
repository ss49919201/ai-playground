package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ss49919201/ai-kata/go/mcp-password-generator/internal/config"
	"github.com/ss49919201/ai-kata/go/mcp-password-generator/pkg/mcp"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	server := mcp.NewServer(cfg)

	mcp.RegisterPasswordGeneratorTool(server, cfg)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Printf("Starting MCP server on %s", addr)
		if err := http.ListenAndServe(addr, server); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
