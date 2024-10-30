package pocketbase

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const pocketbaseURI = "http://localhost:8090"

// UserCredentials represents the user's email and password
type UserCredentials struct {
    Identity string `json:"identity"`
    Password string `json:"password"`
}

// User represents the response structure from PocketBase
type User struct {
    Record struct {
        ID              string `json:"id"`
        CollectionID    string `json:"collectionId"`
        CollectionName  string `json:"collectionName"`
        Username        string `json:"username"`
        Email           string `json:"email"`
        Verified        bool   `json:"verified"`
        // Other fields as needed...
    } `json:"record"`
    Token string `json:"token"`
}

// AdminAuthResponse represents the response structure for admin authentication
type AdminAuthResponse struct {
    Token string `json:"token"`
}

// NewUser represents the structure for creating a new user
type NewUser struct {
    Username      string `json:"username"`
    Email         string `json:"email"`
    Password      string `json:"password"`
    PasswordConfirm string `json:"passwordConfirm"`
    Verified        bool   `json:"verified"`
}

// CreateUserResponseFromString creates a User struct from a JSON string
func CreateUserResponseFromString(data string) (User, error) {
    var response User
    err := json.Unmarshal([]byte(data), &response)
    if err != nil {
        return User{}, err
    }
    return response, nil
}

// Authenticate authenticates the user and returns the token
func Authenticate(credentials UserCredentials) (string, error) {
    url := fmt.Sprintf("%s/api/collections/users/auth-with-password", pocketbaseURI)
    
    jsonData, err := json.Marshal(credentials)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }
    
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("authentication failed: %s", string(body))
    }

    user, err := CreateUserResponseFromString(string(body))
    if err != nil {
        return "", err
    }
    
    return user.Token, nil
}

// GetAdminToken retrieves the admin token from PocketBase
func GetAdminToken(username, password string) (string, error) {
    url := fmt.Sprintf("%s/api/admins/auth-with-password", pocketbaseURI)

    data := UserCredentials{
        Identity: username,
        Password: password,
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("authentication failed: %s", string(body))
    }

    var authResponse AdminAuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
        return "", err
    }

    return authResponse.Token, nil
}

// UserExistsByEmail checks if a user exists with the given email
func UserExistsByEmail(email string, adminToken string) (bool, error) {
    // URL to list all users
    url := fmt.Sprintf("%s/api/collections/users/records", pocketbaseURI)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return false, err
    }

    // Set the Authorization header with the admin token
    req.Header.Set("Authorization", "Bearer "+adminToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()

    // Log the raw response for debugging
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, err
    }

    fmt.Printf("Raw response body: %s\n", string(bodyBytes)) // Log the raw response

    // Check if the response status is OK
    if resp.StatusCode != http.StatusOK {
        return false, fmt.Errorf("error listing users: %s", resp.Status)
    }

    // Decode the response into a suitable struct
    var response struct {
        Total     int `json:"totalItems"`
        Items     []struct {
            Email string `json:"email"`
            Username string `json:"username"` // Add username if needed
        } `json:"items"`
    }

    // Decode the logged response
    err = json.Unmarshal(bodyBytes, &response)
    if err != nil {
        return false, err
    }

    // Iterate over the users to find the one with the specified email
    for _, user := range response.Items {
        if user.Email == email {
            fmt.Printf("User with email %s exists.\n", email)
            return true, nil // User exists
        }
    }

    fmt.Printf("User with email %s does not exist.\n", email)
    return false, nil // User does not exist
}

// CreateUser creates a new user in PocketBase
func CreateUser(newUser NewUser, token string) (string, error) {
    url := fmt.Sprintf("%s/api/collections/users/records", pocketbaseURI)

    jsonData, err := json.Marshal(newUser)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+token)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to create user: %s", string(body))
    }

    user, err := CreateUserResponseFromString(string(body))
    if err != nil {
        return "", err
    }

    return user.Record.ID, nil
}