package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "./uploads"

func main() {

	// Create upload directory
	err := os.Mkdir(uploadDir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Build router
	http.HandleFunc("/upload", uploadHandler)

	// Start server
	fmt.Println("Server started at :3000")
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Check request method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart request
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Retrieve file from the request data
	file, handler, err := r.FormFile("file")
	if (err != nil) || (handler.Filename == "") {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			http.Error(w, "Error closing uploaded file", http.StatusInternalServerError)
		}
	}(file)

	// Create file
	dst, err := os.Create(filepath.Join(uploadDir, handler.Filename))
	if err != nil {
		http.Error(w, "Error creating file", http.StatusBadRequest)
		return
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			http.Error(w, "Error closing created file", http.StatusInternalServerError)
		}
	}(dst)

	// Copy file content
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing to new file", http.StatusInternalServerError)
		return
	}

}
