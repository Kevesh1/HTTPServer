package main

import (
	"fmt"
	"os"
	"regexp"
)

type route struct {
	fileName    string //name of file
	fileContent []byte //fileContent, err
}

//router will store all routes defined for our server, essentially doing all the routing

type router struct {
	routes []route
}

func NewRouter() *router {
	return &router{routes: []route{}}
}

func (r *router) POST(fileName string) {
	path := "C:\\Users\\Keivan\\Downloads\\" + fileName
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error POSTING file: ", err)
	}
	fmt.Println(fileContent)
	route := route{fileName, fileContent}
	r.routes = append(r.routes, route)
}

func (r *router) GET(fileName string) []byte {
	for _, route := range r.routes {
		if fileName == route.fileName {
			return route.fileContent
		}
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
