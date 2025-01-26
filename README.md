# Web Analyzer Service (Go Backend)

This project provides a web page analyzer service that extracts key details from web pages. The service is capable of detecting the HTML version, counting headings, identifying links (both internal and external), checking for inaccessible links, and detecting login forms.

## Features

- **HTML Version Detection**: Detects the HTML version used in a webpage (e.g., HTML5, HTML 4.01).
- **Title Extraction**: Extracts the title of the webpage.
- **Heading Count**: Counts the number of headings by level (H1, H2, etc.).
- **Link Count**: Counts the number of internal and external links, and identifies broken/inaccessible links.
- **Login Form Detection**: Detects the presence of login forms based on input types.
- **Error Handling**: The service gracefully handles errors such as invalid URLs, unreachable pages, and unexpected content types. Appropriate error messages are provided to the user.

## Project Structure

- `internal/`: Contains the core application logic (services, validators, and utils).
- `handlers/`: Contains HTTP handlers for interacting with the service.
- `services/`: Contains the business logic for analyzing the web pages.
- `validators/`: Contains URL validation logic.
- `utils/`: Contains utility functions for extracting data from HTML.

## Prerequisites

- Go 1.23.5 or later
- `git` installed
- Docker installed (If only using Docker to run the application)
- Node.js and npm installed (for the React frontend)

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

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
FRONTEND_URL=http://localhost:3000
SERVER_PORT=8081
```

### Run Locally

1. Start the server:
   ```bash
   go run cmd/main.go
   ```

2. The server will run on `http://localhost:8081` by default.

### Run with Docker

1. Build the Docker image:
   ```bash
   docker build -t web-analyzer-service .
   ```

2. Run the container:
   ```bash
   docker run -p 8081:8081 web-analyzer-service
   ```

3. The server will be available at `http://localhost:8081`.

### Run Backend and Frontend Together

No need to create a `.env` file manually. The required environment variables are set within the script.

1. Ensure the frontend code is in a separate directory (e.g., `../web-analyzer-frontend`).
2. Make the script executable:
   ```bash
   chmod +x run.sh
   ```
3. Run the following script from the `web-analyzer-service` directory:
   ```bash
   ./run.sh
   ```
4. This script will:
   - Install backend and frontend dependencies
   - Start the backend server on port `8081`
   - Start the React frontend on port `3000` (this may take some time)

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
#### Example Response:

```json
{
  "title": "Test Page",
  "htmlVersion": "HTML5",
  "headings": {
    "h1": 1
  },
  "internalLinks": 1,
  "externalLinks": 1,
  "inaccessibleLinks": 1,
  "hasLoginForm": false
}
```
#### Example UI:

![Screenshot from 2025-01-26 21-53-33](https://github.com/user-attachments/assets/7a00b5fb-1e37-4bbd-b029-8c956d04acc4)


### Error Handling

The service provides robust error handling and will return clear error messages in cases like:

- **Invalid URL**: If the URL format is incorrect or unsupported.
- **Page Unreachable**: If the page is not accessible (e.g., network issues, 404 or 500 errors).
- **Invalid Content**: If the page content cannot be parsed correctly (e.g., missing HTML structure).
  
#### Example Error Response:

```json
{
   "status":400,
   "error":"invalid URL format, please provide a valid URL"
}
```
#### Example Error Response on UI:

![Screenshot from 2025-01-26 21-52-46](https://github.com/user-attachments/assets/5a1f02f5-b8c5-43af-a1df-07a2d920f579)

![Screenshot from 2025-01-26 21-54-07](https://github.com/user-attachments/assets/7f89ed2d-969b-4896-9a3a-d2cc4380e05d)

## Testing

The project uses unit tests to validate the functionality of the core services. The following tools are used for testing:

- **Testify**: For assertions and mocking.
- **Gock**: For HTTP request mocking.
- **HTTPMock**: For simulating HTTP responses.
- **Gin**: For HTTP handler testing.

### Checking Test Coverage

To check the test coverage for the project, use the following command:

```bash
go test -cover ./...
```
This command will run the tests and display the coverage percentage, giving you an insight into the test coverage of the project.

### Running Tests

To run the tests, use the following command:

```bash
go test ./...
```

This will execute all the unit tests in the project and provide output on the success or failure of each test.

### Testing Details

- **URL Validator**: Tests for validating the format and reachability of URLs.
- **HTML Analyzer**: Tests for extracting the title, detecting HTML version, counting headings, and checking links.
- **Login Form Detector**: Tests for detecting login forms in HTML pages.


## Additional Notes

- The backend uses the Gin framework.
- Assuming the user provides a reachable URL, return a 400 error when the URL is not reachable due to network or gateway issues.
- In the basic implementation, it took more time to analyze a simple webpage (1.5 minutes), but by using channels and goroutines, the code was optimized to reduce the response time to around 5 seconds.
- Used `zerolog` to add logs and configured it to store the logs in the app.log file.
- CORS is configured to allow requests from `http://localhost:3000`.
- Using Docker ensures a consistent environment across different systems.

## Possible Improvements

While the project provides the basic functionality of web page analysis, there are several areas that could be improved or extended in the future:

1. **Performance Optimization**:
   - Although goroutines and channels have been used to improve performance, further optimization can be achieved by adjusting the concurrency levels based on the size and complexity of the web page being analyzed. For example, URLs like `https://go.dev/dl/` containing over 6000 links take more than a minute to analyze due to network latency. Performance can likely be improved by fine-tuning the concurrency level.

2. **Enhanced Error Handling**:
   - Error handling can be more detailed and cover specific error types such as network issues, malformed HTML, or invalid URLs. This would improve the user experience and make the system more robust.

3. **Caching Mechanism**:
   - Implement a caching mechanism to store results of frequently analyzed pages. This would reduce processing time when the same page is analyzed multiple times.

4. **Integration with CI/CD**:
   - Set up continuous integration and continuous deployment (CI/CD) pipelines for automatic testing and deployment of the application.

5. **User Authentication**:
   - Implement user authentication to allow users to save their past page analyses and view them later. This would add personalization to the app.

6. **Unit Test Coverage**:
    - Increase the unit test coverage to ensure that all components of the project are thoroughly tested, including edge cases.

These improvements can help enhance the functionality, performance, and user experience of the web analysis tool.
