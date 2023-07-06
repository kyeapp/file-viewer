package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/facette/natsort"
)

type FsEntry struct {
	Name     string `json:"name"`
	IsFile   bool   `json:"isFile"`
	IsFolder bool   `json:"isFolder"`
}

type SearchRes struct {
	Name string
	Line []string
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

const indexDir = "bleve-indexes"

func main() {
	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(indexDir)
	if err != nil {
		log.Fatalf("error reading data dir: %v", err)
	}

	for _, dirInfo := range dirEntries {
		indexPath := indexDir + string(os.PathSeparator) + dirInfo.Name()

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("not registering %s, skipping", indexPath)
			continue
		}

		i, err := bleve.Open(indexPath)
		if err != nil {
			log.Printf("error opening index %s: %v", indexPath, err)
			panic("no index")
		}
		log.Printf("registered index: %s", dirInfo.Name())
		// set correct name in stats
		i.SetName(dirInfo.Name())
		i.Close()
	}

	// Create an HTTP handler function
	http.HandleFunc("/directory", serveFoldersAndFiles)
	http.HandleFunc("/search", searchHandler)

	// Start the HTTP server on port 8080
	log.Println("Server listening on port 8080")
	err = http.ListenAndServe(":8080", addCorsHeaders(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	defer measureTime()()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// example query: http://localhost:8095/search?i=hpotter.bleve&q=nimbus
	indexPath := r.URL.Query().Get("i")
	searchTerm := r.URL.Query().Get("q")
	searchResults, err := performSearch(indexPath, searchTerm)
	if err != nil {
		return
	}
	// printStruct(searchResults)
	// fmt.Printf("%v", searchResults)

	hitResp := make([]SearchRes, len(searchResults.Hits))

	for i, hit := range searchResults.Hits {
		hitResp[i].Name = hit.ID
		hitResp[i].Line = hit.Fragments["Line"]
	}

	res := struct {
		SearchStat string
		Hits       []SearchRes
	}{
		SearchStat: fmt.Sprintf("%d results (%s)", searchResults.Total, searchResults.Took),
		Hits:       hitResp,
	}

	jsonResponse, err := json.Marshal(res)
	if err != nil {
		log.Printf("JSON marshaling error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func performSearch(indexPath string, searchTerm string) (*bleve.SearchResult, error) {
	log.Printf(`Searching through index "%s" for "%s"`, indexPath, searchTerm)

	path := indexDir + string(os.PathSeparator) + indexPath
	index, err := bleve.Open(path)
	if err != nil {
		log.Printf("error opening index %s: %v", indexPath, err)
		return nil, err
	}

	indexQuery := bleve.NewMatchQuery(searchTerm)
	searchReq := bleve.NewSearchRequest(indexQuery)
	searchReq.Size = math.MaxInt64
	searchReq.Highlight = bleve.NewHighlight()
	searchResults, err := index.Search(searchReq)
	if err != nil {
		log.Printf("index search error: %v", err)
		return nil, err
	}

	err = index.Close()
	if err != nil {
		log.Printf("close index error: %v", err)
		return nil, err
	}

	return searchResults, nil
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
