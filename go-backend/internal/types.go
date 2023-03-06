package internal

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             string    `json:"uuid" gorm:"primaryKey"`
	From           string    `json:"from_user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Body           string    `json:"body"`
	ConversationId string    `json:"conversation_id"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	m.ID = uuid.New().String()
	return
}

type User struct {
	gorm.Model

	ID        string    `json:"uuid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
}

type Conversation struct {
	gorm.Model

	ID        string    `json:"uuid"`
	CreatedAt time.Time `json:"created_at"`
	Messages  []Message `json:"messages"`
}

func (c *Conversation) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}