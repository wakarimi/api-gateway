package types

// Error represents a standard application error.
type Error struct {
	// The error message.
	// Example: User already exists
	Error string `json:"error" binding:"required"`
}
