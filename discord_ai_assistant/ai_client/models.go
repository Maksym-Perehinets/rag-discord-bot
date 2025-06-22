package ai_client

type ChatMessage struct {
	Role    string `json:"role"`    // "user", "assistant", or "system"
	Content string `json:"content"` // The content of the message
}

type ChatRequest struct {
	UserID    string        `json:"user_id"`    // Unique identifier for the user
	MessageID string        `json:"message_id"` // Unique identifier for the message
	Query     []ChatMessage `json:"query"`      // List of messages in the chat
}

type ChatResponse struct {
	UserID    string      `json:"user_id"`        // Unique identifier for the user
	MessageID string      `json:"message_id"`     // Unique identifier for the message
	Answer    ChatMessage `json:"query_response"` // List of messages in the response
}
