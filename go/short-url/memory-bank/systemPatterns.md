# System Patterns: Go Short URL Server

## 1. Architecture

- **Monolithic Service:** A single Go application handles both URL shortening and redirection.
- **Standard HTTP Server:** Uses Go's built-in `net/http` package for handling requests.
- **In-Memory Storage:** Utilizes a simple map (`map[string]string`) protected by a `sync.RWMutex` for storing URL mappings. This is suitable for MVP but not persistent or scalable.

## 2. Key Technical Decisions

- **Language:** Go (chosen for performance and simplicity in building web services).
- **Storage:** In-memory map (chosen for simplicity in the initial version).
- **ID Generation:** Simple integer counter (`fmt.Sprintf("%d", count)`). This is predictable but not ideal for public-facing services (potential for guessing URLs, limited character set).
- **Routing:** Basic routing using `http.HandleFunc` based on URL path prefixes.

## 3. Component Relationships

```mermaid
graph TD
    Client --> Server[/shorten POST]
    Client --> Server[/{shortID} GET]
    Server --> App[App Struct]
    App --> Store[URLStore Interface]
    Store --> InMemoryStore[InMemoryURLStore]

    subgraph Server [HTTP Server (main.go)]
        direction LR
        HandleShorten[handleShorten]
        HandleRedirect[handleRedirect]
    end

    subgraph Storage [Data Storage (store.go)]
        direction LR
        InMemoryStore
    end

    HandleShorten --> App
    HandleRedirect --> App
```

## 4. Future Considerations

- Replace in-memory store with a persistent database (e.g., Redis, PostgreSQL).
- Implement a more robust short ID generation algorithm (e.g., base62 encoding of a counter, random strings with collision checks).
- Use a dedicated HTTP router (e.g., `gorilla/mux`, `chi`) for more complex routing needs.
- Add configuration management.
- Implement logging and monitoring.
