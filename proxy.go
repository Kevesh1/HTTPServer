package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	//"net/http"
)

func StartProxy() {

	time.Sleep(1 * time.Second)
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	fmt.Println("Proxy: Enter what port to start proxy server from: ")
	//var port string

	//fmt.Scanln(&port)

	port := os.Getenv("PROXY_PORT")

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Proxy: Error ", err)
		return

	}
	defer listener.Close()

	// Create a worker pool for handling concurrent requests
	fmt.Println("Proxy: Running on port: ", port)
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			fmt.Println("Proxy: Error ", err)
			continue
		}
		go handleConn(clientConn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	// Parse the HTTP request
	request, err := http.ReadRequest(bufio.NewReader(conn))

	if err != nil {
		fmt.Println("Proxy: Error fetching request ", err)
		// Handle parsing error and respond with a 400 Bad Request
	}

	if request.Method != "GET" {
		Status(conn, http.StatusNotImplemented)
		conn.Write([]byte("501 Req Method Not Implemented"))
		return
	}
	mainPort := os.Getenv("MAIN_PORT")
	ipAdress := strings.Split(request.Host, ":")[0]
	//fmt.Println(ipAdress)

	mainServer, err := net.Dial("tcp", ipAdress+":"+mainPort)
	if err != nil {
		fmt.Println("Proxy: Error connecting to main server: ", err)
	}
	defer mainServer.Close()

	// Forward the client's GET request to the remote server
	err = request.Write(mainServer)
	if err != nil {
		fmt.Println("Proxy: Error forwarding request ", err)
		// Handle forwarding error and respond with a 500 Internal Server Error
		Status(conn, http.StatusInternalServerError)
		return
	}

	// Read the remote server's response
	mainResponse, err := http.ReadResponse(bufio.NewReader(mainServer), nil)
	if err != nil {
		fmt.Println("Proxy: Error reading response ", err)
		Status(conn, http.StatusNotFound)
		return

	}
	defer mainResponse.Body.Close()

	response := fmt.Sprintf("HTTP/1.1 %s\r\n", mainResponse.Status) +
		fmt.Sprintf("Content-Length: %d\r\n\r\n", mainResponse.ContentLength)

	// Write the constructed response headers to the client's connection
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Printf("Proxy: Error writing response headers: %s\n", err)
		return
	}

	// Copy the main server's response body to the client's connection
	_, err = io.Copy(conn, mainResponse.Body)
	if err != nil {
		fmt.Printf("Proxy: Error copying response body: %s\n", err)
		return
	}
}
