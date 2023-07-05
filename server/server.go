package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/facette/natsort"
)

type FsEntry struct {
	Name     string `json:"name"`
	IsFile   bool   `json:"isFile"`
	IsFolder bool   `json:"isFolder"`
}

func serveFoldersAndFiles(w http.ResponseWriter, r *http.Request) {
	defer measureTime()()
	// Extract the path query parameter from the URL
	queryParams := r.URL.Query()
	path := queryParams.Get("path")

	fsEntries := getFiles(path)
	log.Println("serveDirectory:", path)

	// Encode the list of files as JSON and write it to the response
	err := json.NewEncoder(w).Encode(fsEntries)
	if err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

func measureTime() func() {
	start := time.Now()

	return func() {
		duration := time.Since(start)
		fmt.Println("Execution time:", duration)
	}
}

func getFiles(dir string) []FsEntry {
	var folders []string
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Exclude the root directory itself
		if path == dir {
			return nil
		}

		if info.IsDir() {
			folders = append(folders, filepath.Base(path))
			return filepath.SkipDir
		} else {
			files = append(files, filepath.Base(path))
		}

		return nil
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	// natural sort
	natsort.Sort(folders)
	natsort.Sort(files)

	var entries []FsEntry
	for _, v := range folders {
		entry := FsEntry{
			Name:     v,
			IsFolder: true,
		}
		entries = append(entries, entry)
	}
	for _, v := range files {
		entry := FsEntry{
			Name:   v,
			IsFile: true,
		}
		entries = append(entries, entry)
	}

	return entries
}

func main() {
	// Create an HTTP handler function
	http.HandleFunc("/directory", serveFoldersAndFiles)

	// Start the HTTP server on port 8080
	log.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

// So the server can be hit during dev
func addCorsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			// Handle preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
