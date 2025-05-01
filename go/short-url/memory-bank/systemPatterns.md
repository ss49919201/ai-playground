# System Patterns: Go Short URL Server

## 1. Architecture

- **Monolithic Service:** A single Go application handles URL shortening API, redirection, and serves the static Web UI.
- **Standard HTTP Server:** Uses Go's built-in `net/http` package for handling API requests, serving static files (`http.FileServer`, `http.StripPrefix`), and serving the main HTML page (`http.ServeFile`).
- **In-Memory Storage:** Utilizes a simple map (`map[string]string`) protected by a `sync.RWMutex` for storing URL mappings. This is suitable for MVP but not persistent or scalable.
- **Frontend:** Simple vanilla JavaScript interacts with the backend API via `fetch`.

## 2. Key Technical Decisions

- **Language:** Go (chosen for performance and simplicity in building web services).
- **Storage:** In-memory map (chosen for simplicity in the initial version).
- **ID Generation:** Simple integer counter (`fmt.Sprintf("%d", count)`). This is predictable but not ideal for public-facing services (potential for guessing URLs, limited character set).
- **Routing:** Basic routing using `http.HandleFunc` and `http.Handle`. The root path `/` serves the HTML, `/static/` serves static assets, `/shorten` is the API endpoint, and other paths are treated as potential short IDs for redirection.

## 3. Component Relationships

```mermaid
graph TD
    Browser --> Server[/ GET]
    Browser --> Server[/static/* GET]
    Browser --> Server[/shorten POST]
    Browser --> Server[/{shortID} GET]
    Server --> App[App Struct]
    App --> Store[URLStore Interface]
    Store --> InMemoryStore[InMemoryURLStore]

    subgraph Server [HTTP Server (main.go)]
        direction LR
        RootHandler["/" Handler (Serve index.html or Redirect)]
        StaticHandler["/static/" Handler (FileServer)]
        ShortenHandler["/shorten" Handler (API)]
    end

    subgraph Frontend [Web UI (static/index.html)]
        direction LR
        HTMLForm[HTML Form]
        JSLogic[JavaScript (fetch API)]
    end

     subgraph Storage [Data Storage (store.go)]
        direction LR
        InMemoryStore
    end

    RootHandler --> App
    ShortenHandler --> App
    HTMLForm --> JSLogic
    JSLogic --> ShortenHandler
```

## 4. Future Considerations

- Replace in-memory store with a persistent database (e.g., Redis, PostgreSQL).
- Implement a more robust short ID generation algorithm (e.g., base62 encoding of a counter, random strings with collision checks).
- Use a dedicated HTTP router (e.g., `gorilla/mux`, `chi`) for more complex routing needs.
- Add configuration management.
- Implement logging and monitoring.
