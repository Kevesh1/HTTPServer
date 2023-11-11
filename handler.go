package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var lock sync.Mutex

func POST(fileName string) {

	lock.Lock()
	savePath := filepath.Join("./files/" + fileName)
	writePath := filepath.Join("./test/" + fileName)
	file, err := os.Create(savePath)
	if err != nil {
		fmt.Println("Error POSTING: ", err)
	}
	defer file.Close()

	fileContent, err := os.ReadFile(writePath)
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}
	file.Write(fileContent)
	lock.Unlock()
}

func GET(fileName string) []byte {
	path := "./files/" + fileName
	fileContent, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file: ", err)
	}

	savedFiles, err := os.Open("./files")
	if err != nil {
		fmt.Println("error opening directory: ", err)
	}
	defer savedFiles.Close()

	files, err := savedFiles.Readdir(-1)
	if err != nil {
		fmt.Println("error reading directory:", err)
	}

	for _, files := range files {
		if files.Name() == fileName {
			return fileContent
		}
	}

	fmt.Println("File not found")
	return []byte{}
}
