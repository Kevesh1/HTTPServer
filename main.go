package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
)

func handleRequest(conn net.Conn, router *router) {
	defer conn.Close()

	// Parse the HTTP request
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error fetching request ", err)
		// Handle parsing error and respond with a 400 Bad Request
	}

	if request.Method != "GET" && request.Method != "POST" {
		errorStatus(conn, http.StatusNotImplemented)
		return
	}

	// Handle different HTTP methods (GET, POST, etc.)

	if request.Method == "GET" {
		handleGetRequest(conn, request, router)
	} else if request.Method == "POST" {
		handlePostRequest(conn, request, router)
	}

	// Respond with appropriate status code, headers, and content
}

func handleGetRequest(conn net.Conn, request *http.Request, router *router) {
	requestedPath := request.URL.Path[1:]

	contentType := getContentType(requestedPath)

	if contentType == "" {
		errorStatus(conn, http.StatusBadRequest)
		return
	}

	fileContent := router.GET(requestedPath)
	fmt.Println(fileContent)

	// Create an HTTP response with a 200 OK status and the appropriate headers
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + fmt.Sprint(len(fileContent)) + "\r\n\r\n" +
		string(fileContent)

	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing response: %s\n", err)
	}
}

func handlePostRequest(conn net.Conn, request *http.Request, router *router) {
	file, fileHeader, err := request.FormFile("file")

	if err != nil {
		fmt.Println("Error submitting file: ", err)
	}
	defer file.Close()

	router.POST(fileHeader.Filename)

	response := "HTTP/1.1 201 Created\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"
	conn.Write([]byte(response))

}

func getContentType(fileName string) string {
	// Determine the content type based on the file's extension
	// You can map file extensions to content types (e.g., .html -> "text/html")
	// or use the standard library's mime package for a more comprehensive mapping.
	// For simplicity, you can create a basic mapping for the required file types.
	// Example:
	switch filepath.Ext(fileName) {
	case ".html":
		return "text/html"
	case ".txt":
		return "text/plain"
	case ".gif":
		return "image/gif"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".css":
		return "text/css"
	default:
		return ""
	}
}

func errorStatus(conn net.Conn, status int) {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\r\r\n", status, http.StatusText(status))

	// Write the response to the client's connection
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing response:", err)
	}
}

func main() {
	port := ":8080"

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error ", err)
		return

	}
	defer listener.Close()

	router := NewRouter()

	// Create a worker pool for handling concurrent requests
	fmt.Println("Running on port: ", port)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error ", err)
			continue
		}
		go handleRequest(clientConn, router)
	}
}
