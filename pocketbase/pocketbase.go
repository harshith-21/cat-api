package pocketbase

const pocketbaseURI = "http://localhost:8090"

// UserCredentials represents the user's email and password
type UserCredentials struct {
	Identity string `json:"identity"`
	Password string `json:"password"`
}

// User represents the response structure from PocketBase
type User struct {
	Record struct {
		ID             string `json:"id"`
		CollectionID   string `json:"collectionId"`
		CollectionName string `json:"collectionName"`
		Username       string `json:"username"`
		Email          string `json:"email"`
		Verified       bool   `json:"verified"`
	} `json:"record"`
	Token string `json:"token"`
}

// NewCollection represents the structure of the collection to be created
type NewCollection struct {
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Schema []SchemaField `json:"schema"`
}

// SchemaField represents a field in the collection schema
// SchemaField represents a field in the collection schema
type SchemaField struct {
    Name    string                 `json:"name"`
    Type    string                 `json:"type"`
    Options map[string]interface{} `json:"options,omitempty"` // Add options for fields like "number"
}

// AdminAuthResponse represents the response structure for admin authentication
type AdminAuthResponse struct {
	Token string `json:"token"`
}

// NewUser represents the structure for creating a new user
type NewUser struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
	Verified        bool   `json:"verified"`
}
