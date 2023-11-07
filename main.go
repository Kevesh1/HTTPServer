package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
)

var (
	maxConcurrency = 10
	currentWorkers = 0
	workerLock     sync.Mutex
)
var dataStore = make(map[string][]byte)

func main() {
	port := ":8080" // Change this to the desired port, don't forget docker

	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}

		workerLock.Lock()
		if currentWorkers >= maxConcurrency {
			workerLock.Unlock()
			conn.Close()
			fmt.Println("Max concurrency reached. Waiting for a free worker.")
		} else {
			currentWorkers++
			workerLock.Unlock()
			go handleRequest(conn)
		}
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading request: %s\n", err)
		return
	}

	reader := bufio.NewReader(bytes.NewReader(buffer[:n]))
	request, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Printf("Error parsing request: %s\n", err)
		return
	}

	if request.Method == "GET" {
		handleGetRequest(conn, request)
	} else if request.Method == "POST" {
		// Implement POST request handling here. You should read the request body and store it.
		handlePostRequest(conn, request)
		return
	} else if request.Method == "PUT" {
		// Implement PUT request handling here. You should read the request body and store it.
		handlePutRequest(conn, *request)
	} else if request.Method == "DELETE" {
		// Implement DELETE request handling here. You should delete the file if it exists.
		sendErrorResponse(conn, http.StatusNotImplemented)
		return
	} else {
		sendErrorResponse(conn, http.StatusNotImplemented)
		return
	}
}

func getContentType(filePath string) string {
	switch {
	case strings.HasSuffix(filePath, ".html"):
		return "text/html"
	case strings.HasSuffix(filePath, ".json"):
		return "application/json"
	case strings.HasSuffix(filePath, ".txt"):
		return "text/plain"
	case strings.HasSuffix(filePath, ".gif"):
		return "image/gif"
	case strings.HasSuffix(filePath, ".jpeg"), strings.HasSuffix(filePath, ".jpg"):
		return "image/jpeg"
	case strings.HasSuffix(filePath, ".css"):
		return "text/css"
	default:
		return ""
	}
}

func readFile(filePath string) ([]byte, error) {
	// Implement file reading logic here. You should check if the file exists and return its content.
	// If the file does not exist, return an error.
	// Handle file I/O, errors, and storage as per your requirements.
	return []byte("Hello, World!"), nil
}

/*
 * @Author Oscar Cronvall
 * @Description Handles a GET request
 */
func handleGetRequest(conn net.Conn, request *http.Request) {
	filePath := request.URL.Path
	if filePath == "/" {
		filePath = "/index.html"
	}

	contentType := getContentType(filePath)
	if contentType == "" {
		sendErrorResponse(conn, http.StatusBadRequest)
		return
	}

	data, err := readFile(filePath)
	if err != nil {
		sendErrorResponse(conn, http.StatusNotFound)
		return
	}

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + fmt.Sprint(len(data)) + "\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"\r\n" +
		string(data)

	conn.Write([]byte(response))
	workerLock.Lock()
	currentWorkers--
	workerLock.Unlock()
}

/*
 * @Author Oscar Cronvall
 * @Description Handles a POST request
 */
func handlePostRequest(conn net.Conn, request *http.Request) {

	buffer := make([]byte, 4096)
	n, err := request.Body.Read(buffer)
	filepath := request.URL.Path
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading POST data: %s\n", err)
		sendErrorResponse(conn, http.StatusInternalServerError)
		return
	}

	data := buffer[:n]

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: " + fmt.Sprint(len(data)) + "\r\n" +
		"Content-Type: " + getContentType(filepath) + "\r\n" +
		"\r\n" +
		string(data)

	conn.Write([]byte(response))

	workerLock.Lock()
	currentWorkers--
	workerLock.Unlock()
}

/*
 * @Author Oscar Cronvall
 * @Description Handles a PUT request
 */
func handlePutRequest(conn net.Conn, request http.Request) {
	defer conn.Close()

	filePath := request.URL.Path
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		sendErrorResponse(conn, http.StatusInternalServerError)
		return
	}
	dataStore[filePath] = body

	response := "HTTP/1.1 200 OK\r\n"
	conn.Write([]byte(response))
}

/*
 * @Author Oscar Cronvall
 * @Description Sends a HTTP error response
 */
func sendErrorResponse(conn net.Conn, statusCode int) {
	statusText := http.StatusText(statusCode)
	response := "HTTP/1.1 " + statusText + "\r\n\r\n"
	conn.Write([]byte(response))
}
