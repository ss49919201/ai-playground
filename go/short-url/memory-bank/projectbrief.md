# Project Brief: Go Short URL Server

## 1. Overview

This project aims to create a simple URL shortening service using Go. The server will accept long URLs, generate short unique identifiers for them, store the mapping, and redirect users from the short URL to the original long URL.

## 2. Core Requirements

- **URL Shortening:** Accept a long URL via an API endpoint (e.g., POST /shorten) and return a short URL.
- **Redirection:** Accept a request to a short URL path (e.g., GET /{shortID}) and redirect the user to the corresponding original long URL using an HTTP 301 or 302 redirect.
- **Persistence:** Store the mapping between short IDs and long URLs. Initially, an in-memory store will be used for simplicity.
- **Unique ID Generation:** Implement a mechanism to generate unique and reasonably short identifiers for the URLs.

## 3. Goals

- Implement a functional MVP (Minimum Viable Product) of the URL shortener.
- Write clean, maintainable Go code.
- Provide clear instructions on how to build and run the server.

## 4. Scope (Initial)

- Basic HTTP server implementation.
- In-memory data storage (no database integration initially).
- Simple short ID generation (e.g., random string or counter-based).
- No user authentication or management features.
- Minimal error handling.
- No UI/frontend.

## 5. Target Directory

- `/go/short-url`
