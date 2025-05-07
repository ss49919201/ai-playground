package mcp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ss49919201/ai-playground/go/mcp-password-generator/internal/config"
)

func TestServerHandleToolsList(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		Generator: config.GeneratorConfig{
			DefaultLength: 12,
			MinLength:     8,
			MaxLength:     64,
			UseUppercase:  true,
			UseLowercase:  true,
			UseDigits:     true,
			UseSpecial:    true,
		},
		MCP: config.MCPConfig{
			RequireConsent: true,
		},
	}

	server := NewServer(cfg)

	server.RegisterTool(Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"test": map[string]interface{}{
					"type": "string",
				},
			},
		},
	})

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tools/list",
		"params":  map[string]interface{}{},
		"id":      1,
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSON-RPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.ID != float64(1) {
		t.Errorf("Expected ID %d, got %v", 1, resp.ID)
	}
	if resp.Error != nil {
		t.Errorf("Expected no error, got %v", resp.Error)
	}

	var result ToolsListResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Errorf("Failed to parse result: %v", err)
	}

	if len(result.Tools) != 1 {
		t.Errorf("Expected 1 tool, got %d", len(result.Tools))
	}
	if result.Tools[0].Name != "test_tool" {
		t.Errorf("Expected tool name 'test_tool', got '%s'", result.Tools[0].Name)
	}
}

func TestServerHandleToolsCall(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		Generator: config.GeneratorConfig{
			DefaultLength: 12,
			MinLength:     8,
			MaxLength:     64,
			UseUppercase:  true,
			UseLowercase:  true,
			UseDigits:     true,
			UseSpecial:    true,
		},
		MCP: config.MCPConfig{
			RequireConsent: true,
		},
	}

	server := NewServer(cfg)

	server.RegisterTool(Tool{
		Name:        "test_tool",
		Description: "A test tool",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"test": map[string]interface{}{
					"type": "string",
				},
			},
		},
	})

	server.executePasswordGenerator = func(input json.RawMessage) (interface{}, *Error) {
		return ToolResult{
			Content: []interface{}{
				TextContent{
					Type: "text",
					Text: "Test result",
				},
			},
		}, nil
	}

	consent := true
	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":    "generate_password",
			"input":   map[string]interface{}{},
			"consent": consent,
		},
		"id": 1,
	}
	reqBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSON-RPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.ID != float64(1) {
		t.Errorf("Expected ID %d, got %v", 1, resp.ID)
	}

	reqBody = map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":  "generate_password",
			"input": map[string]interface{}{},
		},
		"id": 2,
	}
	reqBytes, _ = json.Marshal(reqBody)
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSON-RPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.ID != float64(2) {
		t.Errorf("Expected ID %d, got %v", 2, resp.ID)
	}
	if resp.Error == nil {
		t.Errorf("Expected error, got nil")
	}
	if resp.Error.Code != int(ConsentRequired) {
		t.Errorf("Expected error code %d, got %d", ConsentRequired, resp.Error.Code)
	}
}

func TestServerInvalidRequest(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: 8080,
		},
		Generator: config.GeneratorConfig{
			DefaultLength: 12,
			MinLength:     8,
			MaxLength:     64,
			UseUppercase:  true,
			UseLowercase:  true,
			UseDigits:     true,
			UseSpecial:    true,
		},
		MCP: config.MCPConfig{
			RequireConsent: true,
		},
	}

	server := NewServer(cfg)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var resp Response
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSON-RPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.Error == nil {
		t.Errorf("Expected error, got nil")
	}
	if resp.Error.Code != int(ParseError) {
		t.Errorf("Expected error code %d, got %d", ParseError, resp.Error.Code)
	}

	reqBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "invalid_method",
		"params":  map[string]interface{}{},
		"id":      1,
	}
	reqBytes, _ := json.Marshal(reqBody)
	req = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSON-RPC version %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.Error == nil {
		t.Errorf("Expected error, got nil")
	}
	if resp.Error.Code != int(MethodNotFound) {
		t.Errorf("Expected error code %d, got %d", MethodNotFound, resp.Error.Code)
	}
}
