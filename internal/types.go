package internal

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	From           string    `json:"from_user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Body           string    `json:"body"`
	ConversationId string    `json:"conversation_id"`
}

type User struct {
	gorm.Model

	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Conversation struct {
	gorm.Model

	ID        string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}
