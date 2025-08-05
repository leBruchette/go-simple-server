package main

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"time"
)

var logger *zap.Logger

// RequestInfo represents the structure for logging request information
type RequestInfo struct {
	ID          string            `json:"id"`
	Timestamp   string            `json:"timestamp"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	IP          string            `json:"ip"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	Body        interface{}       `json:"body,omitempty"`
}

// logRequest logs the request details as structured JSON using zap
func logRequest(r *http.Request, body interface{}) {
	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0] // Take first value if multiple
		}
	}

	// Convert query parameters to map
	queryParams := make(map[string]string)
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0] // Take first value if multiple
		}
	}

	currentTime := time.Now()
	requestInfo := RequestInfo{
		ID:          fmt.Sprintf("%v", currentTime.UnixNano()),
		Timestamp:   currentTime.Format(time.RFC3339),
		Method:      r.Method,
		Path:        r.URL.Path,
		IP:          r.RemoteAddr,
		Headers:     headers,
		QueryParams: queryParams,
		Body:        body,
	}

	logger.Info("request received",
		zap.String("id", requestInfo.ID),
		zap.String("timestamp", requestInfo.Timestamp),
		zap.String("method", requestInfo.Method),
		zap.String("path", requestInfo.Path),
		zap.String("ip", requestInfo.IP),
		zap.Any("headers", requestInfo.Headers),
		zap.Any("query_params", requestInfo.QueryParams),
		zap.Any("body", requestInfo.Body),
	)
}

// handleGet handles GET requests
func handleGet(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		buildErrorResponse(w)
		return
	}

	// Log the GET request (no body)
	logRequest(r, nil)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":     "GET request received successfully",
		"path":        r.URL.Path,
		"query":       r.URL.Query(),
		"status_code": http.StatusOK,
	}

	json.NewEncoder(w).Encode(response)
}

// handlePost handles POST requests
func handlePost(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		buildErrorResponse(w)
		return
	}
	// Read the request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Try to parse as JSON, fallback to string if not valid JSON
	var bodyData interface{}
	if len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
			// If not valid JSON, store as string
			bodyData = string(bodyBytes)
		}
	}

	// Log the POST request with body
	logRequest(r, bodyData)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":      "POST request received successfully",
		"path":         r.URL.Path,
		"body_length":  len(bodyBytes),
		"content_type": r.Header.Get("Content-Type"),
		"status_code":  http.StatusOK,
	}

	json.NewEncoder(w).Encode(response)
}

func buildErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	response := map[string]interface{}{
		"error":       "Method Not Allowed",
		"status_code": http.StatusMethodNotAllowed,
	}
	json.NewEncoder(w).Encode(response)
}

// healthCheck handles health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	var err error
	logger, err = zap.NewProduction(zap.WithCaller(false))
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create a new HTTP server mux
	mux := http.NewServeMux()

	// Register handlers
	mux.HandleFunc("/get", handleGet)
	mux.HandleFunc("/post", handlePost)
	mux.HandleFunc("/health", healthCheck)

	// Default handler for undefined routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r, nil)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Welcome to the Go Web Server",
			"hint":    "Try /get, /post, or /health endpoints",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Server configuration
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server
	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /get")
	fmt.Println("  POST /post")
	fmt.Println("  GET  /health")
	fmt.Println("  GET  / (default)")
	fmt.Println("\nServer logs will appear below\n")

	logger.Info("server started", zap.String("address", server.Addr))
	log.Fatal(server.ListenAndServe())
}
