package auth

// Simple user model for the hardcoded demo
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // plaintext for demo only — do NOT do this in production
}

// Single hardcoded user: username=admin, password=password123
var HardcodedUser = User{
	ID:       1,
	Username: "admin",
	Password: "password123",
}

