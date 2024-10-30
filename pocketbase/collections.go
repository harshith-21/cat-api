package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateCollection creates a new collection in PocketBase
func CreateCollection(newCollection NewCollection, token string) error {
	url := fmt.Sprintf("%s/api/collections", pocketbaseURI)

	jsonData, err := json.Marshal(newCollection)
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create collection: %s", string(bodyBytes))
	}

	return nil
}

// GetCollection retrieves a collection by its name
func GetCollection(collectionName string, token string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/collections/%s", pocketbaseURI, collectionName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to retrieve collection: %s", string(bodyBytes))
	}

	var collection map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return collection, nil
}

// DeleteCollection deletes a collection by name in PocketBase
func DeleteCollection(collectionName, token string) error {
	url := fmt.Sprintf("%s/api/collections/%s", pocketbaseURI, collectionName)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete collection: %s", string(bodyBytes))
	}

	return nil
}

// Helper function to get collections
func GetCollections(adminToken string) ([]string, error) {
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
