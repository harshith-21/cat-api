package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/harshith-21/cat-api/pocketbase"
)

const (
    adminEmail    = "admin@admin.com"
    adminPassword = "adminadmin"
)

var pocketbaseURI = "http://localhost:8090"

// Home Page
func homePage(w http.ResponseWriter, r *http.Request) {
    // Retrieve existing APIs (collections)
    adminToken, err := pocketbase.GetAdminToken(adminEmail, adminPassword)
    if err != nil {
        http.Error(w, "Unable to authenticate", http.StatusInternalServerError)
        return
    }

    // Retrieve the list of existing APIs
    collections, err := getCollections(adminToken)
    if err != nil {
        http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
        return
    }

    // Prepare data for the template
    data := struct {
        Collections []string
    }{
        Collections: collections,
    }

    // Parse and execute the home template
    tmpl, err := template.ParseFiles("frontend/templates/home.html")
    if err != nil {
        http.Error(w, "Unable to load template", http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, data); err != nil {
        http.Error(w, "Unable to execute template", http.StatusInternalServerError)
    }
}

// Create API Page
func createApiPage(w http.ResponseWriter, r *http.Request) {
    // Render the create API form
    tmpl, err := template.ParseFiles("frontend/templates/create_api.html")
    if err != nil {
        http.Error(w, "Unable to load template", http.StatusInternalServerError)
        return
    }

    if err := tmpl.Execute(w, nil); err != nil {
        http.Error(w, "Unable to execute template", http.StatusInternalServerError)
    }
}

// Create a new API collection
func createApiCollection(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    collectionName := r.FormValue("collectionName")

    // Validate collection name
    if collectionName == "" {
        http.Error(w, "Collection name cannot be empty", http.StatusBadRequest)
        return
    }

    adminToken, err := pocketbase.GetAdminToken(adminEmail, adminPassword)
    if err != nil {
        http.Error(w, "Unable to authenticate", http.StatusInternalServerError)
        return
    }

    // Call the function to create a new collection
    if err := createCollection(adminToken, collectionName); err != nil {
        http.Error(w, "Failed to create collection", http.StatusInternalServerError)
        return
    }

    // Redirect to home after successful creation
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Main function
func main() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/createapi", createApiPage)
    http.HandleFunc("/createapi/submit", createApiCollection)

    // Start the server
    fmt.Println("Starting server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

// Helper function to get collections
func getCollections(adminToken string) ([]string, error) {
    url := fmt.Sprintf("%s/api/collections", pocketbaseURI)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+adminToken)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to list collections: %s", resp.Status)
    }

    var response struct {
        Items []struct {
            Name string `json:"name"`
        } `json:"items"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    collections := make([]string, len(response.Items))
    for i, item := range response.Items {
        collections[i] = item.Name
    }

    return collections, nil
}

// Create a collection in PocketBase
func createCollection(token, name string) error {
    url := fmt.Sprintf("%s/api/collections", pocketbaseURI)
    newCollection := struct {
        Name string `json:"name"`
    }{
        Name: name,
    }

    jsonData, err := json.Marshal(newCollection)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("failed to create collection: %s", string(body))
    }

    return nil
}