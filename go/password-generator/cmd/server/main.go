package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ss49919201/ai-playground/go/password-generator/internal/config"
	"github.com/ss49919201/ai-playground/go/password-generator/pkg/generator"
)

type MCPRequest struct {
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type GeneratePasswordParams struct {
	Length    *int  `json:"length,omitempty"`
	Uppercase *bool `json:"uppercase,omitempty"`
	Lowercase *bool `json:"lowercase,omitempty"`
	Digits    *bool `json:"digits,omitempty"`
	Special   *bool `json:"special,omitempty"`
}

type PasswordResult struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
}

const (
	jsonRPCVersion        = "2.0"
	errCodeParseError     = -32700
	errCodeInvalidRequest = -32600
	errCodeMethodNotFound = -32601
	errCodeInvalidParams  = -32602
	errCodeInternalError  = -32603
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	gen, err := generator.NewGenerator(
		cfg.Generator.MinLength,
		cfg.Generator.MaxLength,
		generator.CharacterSet{
			Uppercase: cfg.Generator.UseUppercase,
			Lowercase: cfg.Generator.UseLowercase,
			Digits:    cfg.Generator.UseDigits,
			Special:   cfg.Generator.UseSpecial,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create password generator: %v", err)
	}

	// Set up signal handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Set up stdio
	scanner := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	// Main loop
	go func() {
		for scanner.Scan() {
			var req MCPRequest
			if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
				sendError(writer, req.ID, errCodeParseError, "Invalid request format")
				continue
			}

			if req.JSONRPC != jsonRPCVersion {
				sendError(writer, req.ID, errCodeInvalidRequest, "Invalid JSON-RPC version")
				continue
			}

			switch req.Method {
			case "generate_password":
				var params GeneratePasswordParams
				if err := json.Unmarshal(req.Params, &params); err != nil {
					sendError(writer, req.ID, errCodeInvalidParams, "Invalid parameters")
					continue
				}

				length := cfg.Generator.DefaultLength
				if params.Length != nil {
					length = *params.Length
				}

				charSet := generator.CharacterSet{
					Uppercase: getPointerBoolOrDefault(params.Uppercase, cfg.Generator.UseUppercase),
					Lowercase: getPointerBoolOrDefault(params.Lowercase, cfg.Generator.UseLowercase),
					Digits:    getPointerBoolOrDefault(params.Digits, cfg.Generator.UseDigits),
					Special:   getPointerBoolOrDefault(params.Special, cfg.Generator.UseSpecial),
				}

				password, err := gen.GenerateWithCharSet(length, charSet)
				if err != nil {
					sendError(writer, req.ID, errCodeInternalError, err.Error())
					continue
				}

				sendResult(writer, req.ID, PasswordResult{
					Password: password,
					Length:   length,
				})

			case "health_check":
				sendResult(writer, req.ID, map[string]string{"status": "ok"})

			default:
				sendError(writer, req.ID, errCodeMethodNotFound, fmt.Sprintf("Unknown method: %s", req.Method))
			}
		}
	}()

	<-quit
	log.Println("Shutting down server...")
}

func sendResult(writer *bufio.Writer, id any, result interface{}) {
	response := MCPResponse{
		JSONRPC: jsonRPCVersion,
		ID:      id,
		Result:  result,
	}
	data, _ := json.Marshal(response)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}

func sendError(writer *bufio.Writer, id any, code int, message string) {
	response := MCPResponse{
		JSONRPC: jsonRPCVersion,
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	data, _ := json.Marshal(response)
	writer.Write(data)
	writer.WriteByte('\n')
	writer.Flush()
}

func getPointerBoolOrDefault(value *bool, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}
