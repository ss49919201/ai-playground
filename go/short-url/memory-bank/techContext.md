# Tech Context: Go Short URL Server

## 1. Language & Runtime

- **Language:** Go (latest stable version recommended)
- **Runtime:** Go runtime environment

## 2. Key Dependencies

- **Standard Library:**
  - `net/http`: For building the HTTP server and handling requests.
  - `fmt`: For formatting strings (used in ID generation and output).
  - `log`: For basic logging (primarily for fatal errors).
  - `sync`: For `RWMutex` to protect concurrent access to the in-memory store.
  - `io`: For reading request bodies.
  - `net/url`: For basic URL validation.
  - `strings`: For string manipulation (trimming path prefix).
- **External Dependencies:** None currently.

## 3. Development Setup

- Go compiler and tools installed.
- A text editor or IDE (like VS Code with the Go extension).
- Terminal for running commands.

## 4. Build & Run

- **Build:** `go build` in the `go/short-url` directory.
- **Run:** `./short-url` (or `go run main.go`) in the `go/short-url` directory.
- The server listens on port `8080` by default.

## 5. Technical Constraints

- Relies on in-memory storage, so data is lost on server restart.
- Simple sequential ID generation might not be suitable for high-load or public scenarios.
- Minimal error handling and validation.
