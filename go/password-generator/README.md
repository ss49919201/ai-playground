# Password Generator MCP Server

A secure password generation server implemented in Go.

## Features

- Generate random passwords with customizable options
- RESTful API with JSON support
- Configurable password length and character sets
- Secure random number generation using crypto/rand
- Comprehensive test suite
- Simple HTML documentation

## API Endpoints

### GET /api/v1/password

Generate a password using query parameters.

**Query Parameters:**
- `length`: Length of the password (default: configured default length)
- `uppercase`: Include uppercase letters (true/false)
- `lowercase`: Include lowercase letters (true/false)
- `digits`: Include digits (true/false)
- `special`: Include special characters (true/false)

**Example:**
```
GET /api/v1/password?length=12&uppercase=true&lowercase=true&digits=true&special=false
```

### POST /api/v1/password

Generate a password using JSON request body.

**Request Body:**
```json
{
  "length": 12,
  "uppercase": true,
  "lowercase": true,
  "digits": true,
  "special": false
}
```

**Response:**
```json
{
  "password": "generated-password",
  "length": 12
}
```

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "ok"
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
```

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

1. Clone the repository
```bash
git clone https://github.com/ss49919201/ai-kata.git
cd ai-kata/go/password-generator
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
password-generator/
├── cmd/
│   └── server/
│       └── main.go         # Entry point for the server
├── internal/
│   ├── config/
│   │   └── config.go       # Configuration handling
│   └── utils/              # Internal utilities
├── pkg/
│   ├── api/
│   │   └── server.go       # HTTP server and API handlers
│   └── generator/
│       ├── generator.go    # Password generation logic
│       └── generator_test.go # Tests for password generator
├── test/                   # Integration tests
├── config.yaml             # Default configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # Project documentation
```
