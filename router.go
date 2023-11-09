package main

import (
	"fmt"
	"os"
	"regexp"
	"sync"
)

type route struct {
	fileName    string //name of file
	fileContent []byte //fileContent, err
}

//router will store all routes defined for our server, essentially doing all the routing

type router struct {
	routes []route
}

var lock sync.Mutex

func NewRouter() *router {
	return &router{routes: []route{}}
}

func (r *router) POST(fileName string) {
	path := "C:\\Users\\keiva\\Downloads\\" + fileName
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error POSTING file: ", err)
	}
	fmt.Println(fileContent)
	for _, route := range r.routes {
		if fileName == route.fileName{
			route.fileContent = fileContent
		}
	} 
	lock.Lock()
	route := route{fileName, fileContent}
	r.routes = append(r.routes, route)
	lock.Unlock()
}

func (r *router) GET(fileName string) []byte {
	for _, route := range r.routes {
		if fileName == route.fileName {
			return route.fileContent
		}
		//TODO
		//Implement 404 not found
	}
	fmt.Println("File not found")
	return []byte{}
}

func isFileNameValid(fileName string) bool {
	// Define a regular expression pattern for a valid file name
	pattern := `^[a-zA-Z0-9_.-]+$`

	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Use MatchString to check if the fileName matches the pattern
	return re.MatchString(fileName)
}
