package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chatmenow/chat-service/internal/config"
	"github.com/chatmenow/chat-service/internal/middleware"
	"github.com/chatmenow/chat-service/internal/model"
	"github.com/chatmenow/chat-service/internal/service"
	"github.com/chatmenow/chat-service/internal/websocket"
	"github.com/google/uuid"
	ws "github.com/gorilla/websocket"
)

type Handler struct {
	config              *config.Config
	messageService      *service.MessageService
	conversationService *service.ConversationService
	hub                 *websocket.Hub
	upgrader            ws.Upgrader
}

func New(
	cfg *config.Config,
	messageService *service.MessageService,
	conversationService *service.ConversationService,
	hub *websocket.Hub,
) *Handler {
	return &Handler{
		config:              cfg,
		messageService:      messageService,
		conversationService: conversationService,
		hub:                 hub,
		upgrader: ws.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins (configure properly in production)
			},
		},
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "chat-service",
	})
}

func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Extract token from query or header
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 {
				token = parts[1]
			}
		}
	}

	if token == "" {
		http.Error(w, "No token provided", http.StatusUnauthorized)
		return
	}

	// Verify JWT
	claims, err := middleware.VerifyJWT(token, h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade connection
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create client
	client := websocket.NewClient(h.hub, conn, claims.Sub)
	h.hub.RegisterClient(client)

	// Start pumps
	go client.WritePump()
	go client.ReadPump()
}

func (h *Handler) ConversationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getConversations(w, r)
	case http.MethodPost:
		h.createConversation(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getConversations(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse user ID from JWT
	userID, err := uuid.Parse(user.Sub)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conversations, err := h.conversationService.GetByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func (h *Handler) createConversation(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(user.Sub)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req model.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	conversation, err := h.conversationService.Create(r.Context(), &req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversation)
}

func (h *Handler) ConversationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract conversation ID from path
	path := strings.TrimPrefix(r.URL.Path, "/conversations/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	conversationID := parts[0]

	if len(parts) >= 2 && parts[1] == "messages" {
		h.getMessages(w, r, conversationID)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

func (h *Handler) getMessages(w http.ResponseWriter, r *http.Request, conversationIDStr string) {
	// Parse conversation ID
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	// Get query parameters for pagination
	limit := 50
	offset := 0
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := h.messageService.GetMessages(r.Context(), conversationID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *Handler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(user.Sub)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req model.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	msg := &model.Message{
		ConversationID: req.ConversationID,
		SenderID:       userID,
		Content:        req.Content,
		Type:           req.Type,
		Metadata:       req.Metadata,
	}

	if err := h.messageService.Create(r.Context(), msg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast via WebSocket (use string representation for hub)
	h.hub.BroadcastToConversation(req.ConversationID.String(), map[string]interface{}{
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msg)
}
