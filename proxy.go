package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	//"net/http"
)

func StartProxy() {

	time.Sleep(1 * time.Second)

	fmt.Println("Enter what port to start proxy server from: ")
	var port string

	fmt.Scanln(&port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error ", err)
		return

	}
	defer listener.Close()

	// Create a worker pool for handling concurrent requests
	fmt.Println("Running on port: ", port)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error ", err)
			continue
		}
		fmt.Println("BEFORE HANDLE")
		go handleConn(clientConn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// Parse the HTTP request
	request, err := http.ReadRequest(bufio.NewReader(conn))

	if err != nil {
		fmt.Println("Error fetching request ", err)
		// Handle parsing error and respond with a 400 Bad Request
	}

	if request.Method != "GET" {
		Status(conn, http.StatusNotImplemented)
		return
	}

	mainServer, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to main server: ", err)
	}
	defer mainServer.Close()

	// Forward the client's GET request to the remote server
	err = request.Write(mainServer)
	if err != nil {
		fmt.Println("Error forwarding request ", err)
		// Handle forwarding error and respond with a 500 Internal Server Error
		Status(conn, http.StatusInternalServerError)
		return
	}

	// Read the remote server's response
	mainResponse, err := http.ReadResponse(bufio.NewReader(mainServer), nil)
	if err != nil {
		fmt.Println("Error reading response ", err)
		Status(conn, http.StatusNotFound)
		return

	}
	defer mainResponse.Body.Close()

	response := fmt.Sprintf("HTTP/1.1 %s\r\n", mainResponse.Status) +
		fmt.Sprintf("Content-Length: %d\r\n\r\n", mainResponse.ContentLength)

	// Write the constructed response headers to the client's connection
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Error writing response headers: %s\n", err)
		return
	}

	// Copy the main server's response body to the client's connection
	_, err = io.Copy(conn, mainResponse.Body)
	if err != nil {
		fmt.Printf("Error copying response body: %s\n", err)
		return
	}
}
