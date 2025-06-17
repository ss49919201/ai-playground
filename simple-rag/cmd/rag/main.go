package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"simple-rag/internal/config"
	"simple-rag/internal/document"
	"simple-rag/internal/llm"
	"simple-rag/internal/vector"
)

func main() {
	var (
		configPath = flag.String("config", "config.yaml", "Path to configuration file")
		command    = flag.String("cmd", "interactive", "Command to run: interactive, add, query, list")
		filePath   = flag.String("file", "", "File path for add command")
		query      = flag.String("query", "", "Query for query command")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Printf("Failed to load config, using defaults: %v", err)
		cfg = config.GetDefaultConfig()
	}

	// Initialize components
	db := vector.NewDatabase(cfg.VectorDB.StoragePath)
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	embeddingClient := vector.NewEmbeddingClient(cfg.Embedding.URL, cfg.Embedding.Model)
	llmClient := llm.NewClient(cfg.LLM.URL, cfg.LLM.Model, cfg.LLM.Temperature, cfg.LLM.MaxTokens)
	docReader := document.NewReader(cfg.Document.ChunkSize, cfg.Document.ChunkOverlap)

	// Create RAG system
	ragSystem := &RAGSystem{
		db:              db,
		embeddingClient: embeddingClient,
		llmClient:       llmClient,
		docReader:       docReader,
		config:          cfg,
	}

	// Execute command
	switch *command {
	case "add":
		if *filePath == "" {
			log.Fatal("File path is required for add command")
		}
		if err := ragSystem.AddDocument(*filePath); err != nil {
			log.Fatalf("Failed to add document: %v", err)
		}
		fmt.Printf("Document added successfully: %s\n", *filePath)

	case "query":
		if *query == "" {
			log.Fatal("Query is required for query command")
		}
		response, err := ragSystem.Query(*query)
		if err != nil {
			log.Fatalf("Failed to query: %v", err)
		}
		printResponse(response)

	case "list":
		documents := ragSystem.ListDocuments()
		printDocuments(documents)

	case "interactive":
		runInteractive(ragSystem)

	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func runInteractive(ragSystem *RAGSystem) {
	fmt.Println("Simple RAG System - Interactive Mode")
	fmt.Println("Commands:")
	fmt.Println("  add <file_path>  - Add a document")
	fmt.Println("  query <question> - Ask a question")
	fmt.Println("  list             - List all documents")
	fmt.Println("  help             - Show this help")
	fmt.Println("  exit             - Exit the program")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		command := parts[0]

		switch command {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <file_path>")
				continue
			}
			filePath := parts[1]
			if err := ragSystem.AddDocument(filePath); err != nil {
				fmt.Printf("Error adding document: %v\n", err)
			} else {
				fmt.Printf("Document added successfully: %s\n", filePath)
			}

		case "query":
			if len(parts) < 2 {
				fmt.Println("Usage: query <question>")
				continue
			}
			question := parts[1]
			response, err := ragSystem.Query(question)
			if err != nil {
				fmt.Printf("Error querying: %v\n", err)
			} else {
				printResponse(response)
			}

		case "list":
			documents := ragSystem.ListDocuments()
			printDocuments(documents)

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  add <file_path>  - Add a document")
			fmt.Println("  query <question> - Ask a question")
			fmt.Println("  list             - List all documents")
			fmt.Println("  help             - Show this help")
			fmt.Println("  exit             - Exit the program")

		case "exit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
		}
	}
}