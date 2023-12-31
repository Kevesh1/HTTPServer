package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Parse the HTTP request
	request, err := http.ReadRequest(bufio.NewReader(conn))

	if err != nil {
		fmt.Println("Server: Error fetching request ", err)
		Status(conn, http.StatusBadRequest)
		// Handle parsing error and respond with a 400 Bad Request
	}

	if request.Method != "GET" && request.Method != "POST" {
		Status(conn, http.StatusNotImplemented)
		fmt.Println("Server: NOT IMPLEMENTED")
		return
	}

	// Handle different HTTP methods (GET, POST, etc.)

	if request.Method == "GET" {
		handleGetRequest(conn, request)
	} else if request.Method == "POST" {
		handlePostRequest(conn, request)
	}

	// Respond with appropriate status code, headers, and content
}

func handleGetRequest(conn net.Conn, request *http.Request) {
	requestedPath := request.URL.Path[1:]

	contentType := GetContentType(requestedPath)

	if contentType == "" {
		Status(conn, http.StatusBadRequest)
		return
	}

	fileContent := GET(requestedPath)
	if len(fileContent) == 0 {
		Status(conn, http.StatusNotFound)
		return
	}
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

func handlePostRequest(conn net.Conn, request *http.Request) {
	file, fileHeader, err := request.FormFile("file")

	if err != nil {
		fmt.Println("Server: Error submitting file: ", err)
	}
	defer file.Close()

	POST(fileHeader.Filename)

	response := "HTTP/1.1 201 Created\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"
	conn.Write([]byte(response))

}

func GetContentType(fileName string) string {
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

func Status(conn net.Conn, status int) {
	response := fmt.Sprintf("HTTP/1.1 %d %s\r\r\r\n", status, http.StatusText(status))

	// Write the response to the client's connection
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Server: Error writing response:", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	//port := ":8080"
	fmt.Println("Server: Enter what port to listen from: ")
	//var port string

	//fmt.Scanln(&port)
	port := os.Getenv("MAIN_PORT")

	go StartProxy()

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Server: Error ", err)
		return

	}
	defer listener.Close()

	max_processess := 10
	process := make(chan int, max_processess)

	// Create a worker pool for handling concurrent requests
	fmt.Println("Server: Running on port: ", port)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server: Error ", err)
			continue
		}
		process <- 1
		go func() {
			fmt.Println(len(process))
			handleRequest(clientConn)
			<-process // Release the worker when done
			fmt.Println("Server: Request finished", len(process))
			fmt.Println("_______________")
		}()
	}
}
