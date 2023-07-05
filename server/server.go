package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// type File struct {
// 	Name string `json:"name"`
// }
var dummyFiles = [10]string{
	"file1.txt",
	"file2.txt",
	"file3.txt",
	"file4.txt",
	"file5.txt",
	"file6.txt",
	"file7.txt",
	"file8.txt",
	"file9.txt",
	"file10.txt",
}

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
	// Define the list of filenames
	// Create an HTTP handler function
	http.HandleFunc("/files", serveFiles)

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
