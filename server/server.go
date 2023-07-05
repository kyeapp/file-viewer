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

// type File struct {
// 	Name string `json:"name"`
// }
var dummyFiles []string

func serveFiles(w http.ResponseWriter, r *http.Request) {
	defer measureTime()()
	log.Println("/files")

	// Encode the list of files as JSON and write it to the response
	err := json.NewEncoder(w).Encode(dummyFiles)
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

func main() {
	datadir := "./data"
	err := filepath.Walk(datadir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filename := filepath.Base(path)
			dummyFiles = append(dummyFiles, filename)
		}

		return nil
	})
	if err != nil {
		log.Fatal("Error:", err)
	}

	// natural sort
	natsort.Sort(dummyFiles)

	// Create an HTTP handler function
	http.HandleFunc("/files", serveFiles)

	// Start the HTTP server on port 8080
	log.Println("Server listening on port 8080")
	err = http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux))
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
