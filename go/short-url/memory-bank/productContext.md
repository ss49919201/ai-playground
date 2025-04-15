# Product Context: Go Short URL Server

## 1. Problem Solved

Provides a simple way to shorten long, unwieldy URLs into shorter, more manageable links, suitable for sharing or use in character-limited environments.

## 2. How it Works

- **Shortening (API):** Users send a POST request to `/shorten` with the long URL in the request body. The server generates a unique short ID, stores the mapping (short ID -> long URL), and returns the short URL (e.g., `http://localhost:8080/{shortID}`).
- **Shortening (Web UI):** Users navigate to the root path (`/`) in their browser, enter a long URL into the form, and submit it. The JavaScript on the page sends a request to the `/shorten` API and displays the resulting short URL.
- **Redirection:** Users access the short URL (e.g., `http://localhost:8080/{shortID}`). The server looks up the corresponding long URL based on the `shortID` and redirects the user's browser to the original long URL using an HTTP 302 Found status code.

## 3. User Experience Goals

- **Simplicity:** Easy to use via simple HTTP requests.
- **Speed:** Fast redirection from short URL to long URL.
- **Reliability:** Consistently redirects to the correct original URL.
- **Accessibility:** Provides a simple web form for non-technical users to shorten URLs.

## 4. Target Users

- Developers or systems needing a basic internal URL shortening capability via API.
- End-users needing a simple web tool to shorten URLs.
