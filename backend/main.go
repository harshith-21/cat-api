package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "text/template"
)

// Load the template files
var tmpl = template.Must(template.ParseFiles("frontend/templates/index.html"))

// Handler for the index route
func indexHandler(w http.ResponseWriter, r *http.Request) {
    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// API handler (example endpoint)
func apiHandler(w http.ResponseWriter, r *http.Request) {
    pbURL := os.Getenv("POCKETBASE_URL") 
	// pbURL := "localhost:8090"
    response := fmt.Sprintf("This is data from PocketBase API running at: %s", pbURL)
    w.Write([]byte(response))
	fmt.Println(response)
}

func main() {
    // Serve static files (CSS, etc.)
    fs := http.FileServer(http.Dir("./frontend/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Handle routes
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/api", apiHandler)

    // Start the server
    fmt.Println("Server is running on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}