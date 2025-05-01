package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/ss49919201/ai-playground/go/mcp-password-generator/internal/config"
	"github.com/ss49919201/ai-playground/go/mcp-password-generator/pkg/generator"
)

type PasswordGeneratorInput struct {
	Length    *int  `json:"length,omitempty"`
	Uppercase *bool `json:"uppercase,omitempty"`
	Lowercase *bool `json:"lowercase,omitempty"`
	Digits    *bool `json:"digits,omitempty"`
	Special   *bool `json:"special,omitempty"`
}

type PasswordGeneratorOutput struct {
	Password string `json:"password"`
	Length   int    `json:"length"`
}

func RegisterPasswordGeneratorTool(server *Server, cfg *config.Config) {
	inputSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"length": map[string]interface{}{
				"type":        "integer",
				"description": "Length of the password to generate",
				"minimum":     cfg.Generator.MinLength,
				"maximum":     cfg.Generator.MaxLength,
				"default":     cfg.Generator.DefaultLength,
			},
			"uppercase": map[string]interface{}{
				"type":        "boolean",
				"description": "Include uppercase letters",
				"default":     cfg.Generator.UseUppercase,
			},
			"lowercase": map[string]interface{}{
				"type":        "boolean",
				"description": "Include lowercase letters",
				"default":     cfg.Generator.UseLowercase,
			},
			"digits": map[string]interface{}{
				"type":        "boolean",
				"description": "Include digits",
				"default":     cfg.Generator.UseDigits,
			},
			"special": map[string]interface{}{
				"type":        "boolean",
				"description": "Include special characters",
				"default":     cfg.Generator.UseSpecial,
			},
		},
	}

	server.RegisterTool(Tool{
		Name:        "generate_password",
		Description: "Generates a secure random password with customizable options",
		InputSchema: inputSchema,
		Annotations: map[string]interface{}{
			"security": map[string]interface{}{
				"requiresConsent": cfg.MCP.RequireConsent,
			},
		},
	})

	server.passwordGeneratorHandler = func(input json.RawMessage) (interface{}, *Error) {
		return executePasswordGenerator(input, cfg)
	}
}

func executePasswordGenerator(input json.RawMessage, cfg *config.Config) (interface{}, *Error) {
	var params PasswordGeneratorInput
	if err := json.Unmarshal(input, &params); err != nil {
		return nil, NewError(InvalidParams, "Invalid parameters", err.Error())
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

	gen, err := generator.NewGenerator(
		cfg.Generator.MinLength,
		cfg.Generator.MaxLength,
		charSet,
	)
	if err != nil {
		return nil, NewError(InternalError, "Failed to create generator", err.Error())
	}

	password, err := gen.Generate(length)
	if err != nil {
		return nil, NewError(ToolExecutionFail, "Failed to generate password", err.Error())
	}

	result := ToolResult{
		Content: []interface{}{
			TextContent{
				Type: "text",
				Text: fmt.Sprintf("Generated password: %s", password),
			},
			PasswordGeneratorOutput{
				Password: password,
				Length:   length,
			},
		},
	}

	return result, nil
}

func getPointerBoolOrDefault(value *bool, defaultValue bool) bool {
	if value == nil {
		return defaultValue
	}
	return *value
}
