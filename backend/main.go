package main

import (
	"fmt"
	"html/template"
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
	// Get admin token
	adminToken, err := pocketbase.GetAdminToken(adminEmail, adminPassword)
	if err != nil {
		http.Error(w, "Unable to authenticate", http.StatusInternalServerError)
		return
	}

	// Retrieve the list of existing APIs
	collections, err := pocketbase.GetCollections(adminToken)
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

    // Get the admin token
    adminToken, err := pocketbase.GetAdminToken(adminEmail, adminPassword)
    if err != nil {
        http.Error(w, "Unable to authenticate", http.StatusInternalServerError)
        return
    }

    // Define the new collection schema
    newCollection := pocketbase.NewCollection{
		Name: collectionName,
		Type: "base", 
		Schema: []pocketbase.SchemaField{
			{
				Name: "field1",
				Type: "text",
			},
			{
				Name: "jsonField",
				Type: "json",
				Options: map[string]interface{}{
					"minSize": 0,    // Optional: set minimum value
					"maxSize": 1000, // Optional: set maximum value
				},
			},
		},
    }

    // Call the function to create the collection
    err = pocketbase.CreateCollection(newCollection, adminToken)
    if err != nil {
        fmt.Printf("Error creating collection: %v\n", err)
        http.Error(w, "Error creating collection", http.StatusInternalServerError)
        return
    }

    fmt.Printf("Collection %s created successfully!", collectionName)

    // Redirect to home after successful creation
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Main function
func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/createapi", createApiPage)
	http.HandleFunc("/createapi/submit", createApiCollection)

	// TODO : get env variables for pocketbase url, admin user, admin password. use default if not provided

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
