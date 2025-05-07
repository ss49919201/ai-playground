package mcp

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ss49919201/ai-playground/go/mcp-password-generator/internal/config"
)

func TestExecutePasswordGenerator(t *testing.T) {
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

	input := json.RawMessage(`{}`)
	result, err := executePasswordGenerator(input, cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	toolResult, ok := result.(ToolResult)
	if !ok {
		t.Errorf("Expected ToolResult, got %T", result)
	}

	if len(toolResult.Content) < 2 {
		t.Errorf("Expected at least 2 content items, got %d", len(toolResult.Content))
	}

	outputJSON, err := json.Marshal(toolResult.Content[1])
	if err != nil {
		t.Errorf("Failed to marshal output: %v", err)
	}

	var output PasswordGeneratorOutput
	if err := json.Unmarshal(outputJSON, &output); err != nil {
		t.Errorf("Failed to unmarshal output: %v", err)
	}

	if output.Length != cfg.Generator.DefaultLength {
		t.Errorf("Expected length %d, got %d", cfg.Generator.DefaultLength, output.Length)
	}
	if len(output.Password) != cfg.Generator.DefaultLength {
		t.Errorf("Expected password length %d, got %d", cfg.Generator.DefaultLength, len(output.Password))
	}

	length := 16
	uppercase := false
	lowercase := true
	digits := true
	special := false
	input = json.RawMessage(`{
		"length": 16,
		"uppercase": false,
		"lowercase": true,
		"digits": true,
		"special": false
	}`)
	result, err = executePasswordGenerator(input, cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	toolResult, ok = result.(ToolResult)
	if !ok {
		t.Errorf("Expected ToolResult, got %T", result)
	}

	if len(toolResult.Content) < 2 {
		t.Errorf("Expected at least 2 content items, got %d", len(toolResult.Content))
	}

	outputJSON, err = json.Marshal(toolResult.Content[1])
	if err != nil {
		t.Errorf("Failed to marshal output: %v", err)
	}

	if err := json.Unmarshal(outputJSON, &output); err != nil {
		t.Errorf("Failed to unmarshal output: %v", err)
	}

	if output.Length != length {
		t.Errorf("Expected length %d, got %d", length, output.Length)
	}
	if len(output.Password) != length {
		t.Errorf("Expected password length %d, got %d", length, len(output.Password))
	}

	hasUppercase := false
	hasLowercase := false
	hasDigit := false
	hasSpecial := false

	for _, char := range output.Password {
		c := string(char)
		if strings.Contains(uppercaseChars, c) {
			hasUppercase = true
		} else if strings.Contains(lowercaseChars, c) {
			hasLowercase = true
		} else if strings.Contains(digitChars, c) {
			hasDigit = true
		} else if strings.Contains(specialChars, c) {
			hasSpecial = true
		}
	}

	if hasUppercase != uppercase {
		t.Errorf("Expected uppercase %v, got %v", uppercase, hasUppercase)
	}
	if hasLowercase != lowercase {
		t.Errorf("Expected lowercase %v, got %v", lowercase, hasLowercase)
	}
	if hasDigit != digits {
		t.Errorf("Expected digits %v, got %v", digits, hasDigit)
	}
	if hasSpecial != special {
		t.Errorf("Expected special %v, got %v", special, hasSpecial)
	}

	input = json.RawMessage(`{
		"length": 4
	}`)
	result, err = executePasswordGenerator(input, cfg)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if err.Code != int(ToolExecutionFail) {
		t.Errorf("Expected error code %d, got %d", ToolExecutionFail, err.Code)
	}
}
