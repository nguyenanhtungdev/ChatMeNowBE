package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/chatmenow/chat-service/internal/model"
	"github.com/chatmenow/chat-service/internal/service"
	"github.com/google/uuid"
)

type Hub struct {
	clients         map[string]*Client          // userID -> Client
	conversations   map[string]map[*Client]bool // conversationID -> set of clients
	broadcast       chan *BroadcastMessage
	register        chan *Client
	unregister      chan *Client
	mu              sync.RWMutex
	messageService  *service.MessageService
	presenceService *service.PresenceService
}

type BroadcastMessage struct {
	ConversationID string
	Message        interface{}
	ExcludeClient  *Client
}

type WSMessage struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
}

func NewHub(messageService *service.MessageService, presenceService *service.PresenceService) *Hub {
	return &Hub{
		clients:         make(map[string]*Client),
		conversations:   make(map[string]map[*Client]bool),
		broadcast:       make(chan *BroadcastMessage, 256),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		messageService:  messageService,
		presenceService: presenceService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case msg := <-h.broadcast:
			h.broadcastToConversation(msg)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client.UserID] = client
	ctx := context.Background()
	h.presenceService.SetOnline(ctx, client.UserID)

	log.Printf("Client registered: %s", client.UserID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client.UserID]; ok {
		delete(h.clients, client.UserID)
		close(client.Send)

		// Remove from all conversations
		for _, clients := range h.conversations {
			delete(clients, client)
		}

		ctx := context.Background()
		h.presenceService.SetOffline(ctx, client.UserID)

		log.Printf("Client unregistered: %s", client.UserID)
	}
}

func (h *Hub) JoinConversation(client *Client, conversationID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.conversations[conversationID] == nil {
		h.conversations[conversationID] = make(map[*Client]bool)
	}

	h.conversations[conversationID][client] = true
	log.Printf("Client %s joined conversation %s", client.UserID, conversationID)
}

func (h *Hub) LeaveConversation(client *Client, conversationID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.conversations[conversationID]; ok {
		delete(clients, client)
	}
}

func (h *Hub) BroadcastToConversation(conversationID string, message interface{}, excludeClient *Client) {
	h.broadcast <- &BroadcastMessage{
		ConversationID: conversationID,
		Message:        message,
		ExcludeClient:  excludeClient,
	}
}

func (h *Hub) broadcastToConversation(msg *BroadcastMessage) {
	h.mu.RLock()
	clients := h.conversations[msg.ConversationID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	data, err := json.Marshal(msg.Message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		if msg.ExcludeClient != nil && client == msg.ExcludeClient {
			continue
		}

		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}

func (h *Hub) HandleClientMessage(client *Client, messageData []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(messageData, &wsMsg); err != nil {
		log.Printf("Error parsing message: %v", err)
		return
	}

	ctx := context.Background()

	switch wsMsg.Type {
	case "join_conversation":
		conversationID, ok := wsMsg.Payload["conversationId"].(string)
		if !ok {
			return
		}
		h.JoinConversation(client, conversationID)

	case "leave_conversation":
		conversationID, ok := wsMsg.Payload["conversationId"].(string)
		if !ok {
			return
		}
		h.LeaveConversation(client, conversationID)

	case "send_message":
		conversationIDStr, _ := wsMsg.Payload["conversationId"].(string)
		content, _ := wsMsg.Payload["content"].(string)

		// Parse UUIDs
		conversationID, err := uuid.Parse(conversationIDStr)
		if err != nil {
			log.Printf("Invalid conversation ID: %v", err)
			return
		}

		senderID, err := uuid.Parse(client.UserID)
		if err != nil {
			log.Printf("Invalid user ID: %v", err)
			return
		}

		msg := &model.Message{
			ConversationID: conversationID,
			SenderID:       senderID,
			Content:        content,
			Type:           "text",
		}

		if err := h.messageService.Create(ctx, msg); err != nil {
			log.Printf("Error saving message: %v", err)
			return
		}

		// Broadcast to conversation
		h.BroadcastToConversation(conversationIDStr, map[string]interface{}{
			"type": "new_message",
			"payload": map[string]interface{}{
				"id":             msg.ID,
				"conversationId": msg.ConversationID,
				"senderId":       msg.SenderID,
				"content":        msg.Content,
				"type":           msg.Type,
				"createdAt":      msg.CreatedAt,
			},
		}, nil)

	case "typing":
		conversationID, _ := wsMsg.Payload["conversationId"].(string)
		isTyping, _ := wsMsg.Payload["isTyping"].(bool)

		if isTyping {
			h.presenceService.StartTyping(ctx, conversationID, client.UserID)
		} else {
			h.presenceService.StopTyping(ctx, conversationID, client.UserID)
		}

		// Broadcast typing status
		h.BroadcastToConversation(conversationID, map[string]interface{}{
			"type": "user_typing",
			"payload": map[string]interface{}{
				"conversationId": conversationID,
				"userId":         client.UserID,
				"isTyping":       isTyping,
			},
		}, client)
	}
}
