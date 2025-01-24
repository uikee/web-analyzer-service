# Web Analyzer Service (Go Backend)

This is the backend service for analyzing web pages. It provides an API endpoint to analyze a webpage and extract useful information.

## Prerequisites

- Go 1.18 or later
- `git` installed
- `.env` file created in the project root

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
FRONTEND_URL=http://localhost:3000
SERVER_PORT=8081
```

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd web-analyzer-service
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the App

1. Start the server:
   ```bash
   go run main.go
   ```

2. The server will run on `http://localhost:8081` by default.

## Endpoints

### Analyze Webpage
- **URL:** `/analyze`
- **Method:** `GET`
- **Query Parameters:**
  - `url` (required): The URL of the webpage to analyze.

Example:
```bash
curl "http://localhost:8081/analyze?url=https://example.com"
```

---

## Additional Notes

- The backend uses the Gin framework.
- CORS is configured to allow requests from the `FRONTEND_URL` defined in `.env`.
- For production, update the `.env` file with your frontend's URL.

