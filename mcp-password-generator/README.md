# Password Generator MCP Server

A secure password generation server implemented in Go following the Model Context Protocol (MCP) specification.

## Features

- Implements the Model Context Protocol (MCP) for standardized communication
- JSON-RPC 2.0 based server with MCP-compliant endpoints
- Secure random password generation with customizable options
- Human-in-the-loop consent mechanism for tool execution
- Comprehensive test suite for both MCP protocol and password generation
- Configurable via YAML configuration file

## MCP Protocol Implementation

This server implements the [Model Context Protocol (MCP)](https://spec.modelcontextprotocol.io/specification/2025-03-26/) which provides a standardized way for language models to interact with external tools. The implementation includes:

- `tools/list` endpoint to discover available tools
- `tools/call` endpoint to execute tools
- Proper error handling according to MCP specification
- Human-in-the-loop consent mechanism for security

## Password Generation Tool

The server exposes a password generation tool with the following features:

- Configurable password length
- Options to include/exclude uppercase letters, lowercase letters, digits, and special characters
- Secure random number generation using crypto/rand
- Validation to ensure generated passwords meet requirements

## API Usage

### Listing Available Tools

```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "params": {},
  "id": 1
}
```

Response:

```json
{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "generate_password",
        "description": "Generates a secure random password with customizable options",
        "inputSchema": {
          "type": "object",
          "properties": {
            "length": {
              "type": "integer",
              "description": "Length of the password to generate",
              "minimum": 8,
              "maximum": 64,
              "default": 12
            },
            "uppercase": {
              "type": "boolean",
              "description": "Include uppercase letters",
              "default": true
            },
            "lowercase": {
              "type": "boolean",
              "description": "Include lowercase letters",
              "default": true
            },
            "digits": {
              "type": "boolean",
              "description": "Include digits",
              "default": true
            },
            "special": {
              "type": "boolean",
              "description": "Include special characters",
              "default": true
            }
          }
        },
        "annotations": {
          "security": {
            "requiresConsent": true
          }
        }
      }
    ]
  },
  "id": 1
}
```

### Calling the Password Generator Tool

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "generate_password",
    "input": {
      "length": 16,
      "uppercase": true,
      "lowercase": true,
      "digits": true,
      "special": false
    },
    "consent": true
  },
  "id": 2
}
```

Response:

```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Generated password: Abc123DefGhi456J"
      },
      {
        "password": "Abc123DefGhi456J",
        "length": 16
      }
    ]
  },
  "id": 2
}
```

## Configuration

The server can be configured using a YAML configuration file:

```yaml
server:
  port: 8080

generator:
  defaultLength: 12
  minLength: 8
  maxLength: 64
  useUppercase: true
  useLowercase: true
  useDigits: true
  useSpecial: true

mcp:
  requireConsent: true
```

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

1. Clone the repository

```bash
git clone https://github.com/ss49919201/ai-playground.git
cd ai-playground/go/mcp-password-generator
```

2. Install dependencies

```bash
go mod tidy
```

### Running the Server

```bash
go run cmd/server/main.go
```

Or with a custom configuration file:

```bash
go run cmd/server/main.go -config=/path/to/config.yaml
```

### Running Tests

```bash
go test ./...
```

## Project Structure

```
mcp-password-generator/
├── cmd/
│   └── server/
│       └── main.go         # Entry point for the server
├── internal/
│   └── config/
│       └── config.go       # Configuration handling
├── pkg/
│   ├── generator/
│   │   └── generator.go    # Password generation logic
│   └── mcp/
│       ├── types.go        # MCP protocol types
│       ├── server.go       # MCP server implementation
│       └── password_tool.go # Password generator tool implementation
├── config.yaml             # Default configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
```

## Security Considerations

- The server uses cryptographically secure random number generation
- MCP protocol implementation includes consent mechanisms for tool execution
- Input validation is performed to prevent security issues
- Error handling follows MCP specification to avoid information leakage
