package response

// Error is the format of the response to the request in case an error occurred
type Error struct {
	// Human-readable error message
	Message string `json:"message"`
	// Internal error description
	Reason string `json:"reason"`
}
