package internal

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model

	UserID         string    `json:"from"`
	CreatedAt      time.Time `json:"created_at"`
	Body           string    `json:"body"`
	ConversationId string    `json:"conversation_id"`
	User           User      `json:"user" gorm:"foreignKey:UserID"`
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
