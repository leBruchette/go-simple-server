
# Go Simple Server Example

This project is a simple web server written in Go. It demonstrates structured logging with [Uber Zap](https://github.com/uber-go/zap), request handling for GET and POST endpoints, and a health check endpoint.  This was build simply to log GET and POST request information in a publicly hosted environment

## Features
- **Endpoints:**
    - `GET  /get`
    - `POST /post`
    - `GET  /health`
    - `GET  /` (default)


- **Structured logging** of all requests using Zap
-  Graceful error handling for unsupported methods

## Requirements

- Go 1.18 or newer

## Getting Started

1. **Clone the repository:**
   ```sh
   git clone <your-repo-url>
   cd <project-directory>
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Run the server:**
   ```sh
   go run main.go
   ```

## Running with Docker

You can run this project using Docker or Docker Compose.  Both `Dockerfile` and `docker-compose.yml` are provided.

### Using Docker

1. Build the Docker image and run the container:
   ```sh
   docker build -t go-simple-server . && \
   docker run -p 8080:8080 go-simple-server
   ```

### Using Docker Compose 
1. Start the service:
   ```sh
   docker-compose up
   ```



## Example Requests

- **GET request:**
  ```sh
  curl http://localhost:8080/get
  ```

- **POST request:**
  ```sh
  curl -X POST -H "Content-Type: application/json" -d '{"foo":"bar"}' http://localhost:8080/post
  ```

- **Health check:**
  ```sh
  curl http://localhost:8080/health
  ```

## Logging

All requests are logged in structured JSON format using Zap.

## License

MIT