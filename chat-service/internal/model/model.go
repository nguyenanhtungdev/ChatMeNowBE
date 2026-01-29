package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID             uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID uuid.UUID              `json:"conversationId" gorm:"type:uuid;not null;index"`
	SenderID       uuid.UUID              `json:"senderId" gorm:"type:uuid;not null;index"`
	Content        string                 `json:"content" gorm:"type:text;not null"`
	Type           string                 `json:"type" gorm:"type:varchar(20);not null;default:'text'"` // text, image, file, video
	Metadata       map[string]interface{} `json:"metadata,omitempty" gorm:"type:jsonb"`
	CreatedAt      time.Time              `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt         `json:"-" gorm:"index"`
}

func (Message) TableName() string {
	return "messages"
}

type Conversation struct {
	ID        uuid.UUID            `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string               `json:"name" gorm:"type:varchar(255)"`
	Type      string               `json:"type" gorm:"type:varchar(20);not null"` // direct, group
	AvatarURL string               `json:"avatarUrl,omitempty" gorm:"type:varchar(500)"`
	CreatedBy uuid.UUID            `json:"createdBy" gorm:"type:uuid;not null"`
	CreatedAt time.Time            `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time            `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt       `json:"-" gorm:"index"`
	Members   []ConversationMember `json:"members,omitempty" gorm:"foreignKey:ConversationID"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type ConversationMember struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID uuid.UUID      `json:"conversationId" gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID      `json:"userId" gorm:"type:uuid;not null;index"`
	Role           string         `json:"role" gorm:"type:varchar(20);not null;default:'member'"` // admin, member
	JoinedAt       time.Time      `json:"joinedAt" gorm:"autoCreateTime"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ConversationMember) TableName() string {
	return "conversation_members"
}

// Request DTOs
type CreateConversationRequest struct {
	Name      string      `json:"name" binding:"required"`
	Type      string      `json:"type" binding:"required,oneof=direct group"`
	MemberIDs []uuid.UUID `json:"memberIds" binding:"required,min=1"`
}

type SendMessageRequest struct {
	ConversationID uuid.UUID              `json:"conversationId" binding:"required"`
	Content        string                 `json:"content" binding:"required"`
	Type           string                 `json:"type" binding:"required,oneof=text image file video"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

type GetMessagesRequest struct {
	ConversationID uuid.UUID `json:"conversationId" binding:"required"`
	Limit          int       `json:"limit" binding:"min=1,max=100"`
	Offset         int       `json:"offset" binding:"min=0"`
}

// WebSocket Message Types
type WSMessage struct {
	Type    string      `json:"type"` // join, leave, message, typing, read
	Payload interface{} `json:"payload"`
}

type JoinRoomPayload struct {
	ConversationID uuid.UUID `json:"conversationId"`
	UserID         uuid.UUID `json:"userId"`
}

type TypingPayload struct {
	ConversationID uuid.UUID `json:"conversationId"`
	UserID         uuid.UUID `json:"userId"`
	IsTyping       bool      `json:"isTyping"`
}
