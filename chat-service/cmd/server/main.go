package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chatmenow/chat-service/internal/config"
	"github.com/chatmenow/chat-service/internal/handler"
	"github.com/chatmenow/chat-service/internal/middleware"
	"github.com/chatmenow/chat-service/internal/model"
	"github.com/chatmenow/chat-service/internal/repository"
	"github.com/chatmenow/chat-service/internal/service"
	"github.com/chatmenow/chat-service/internal/websocket"
	"gorm.io/gorm"
)

// autoMigrate runs database migrations
func autoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&model.Conversation{},
		&model.ConversationMember{},
		&model.Message{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func main() {

	cfg := config.Load()

	// Auto-migrate database schema
	if err := autoMigrate(cfg.DB); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repositories
	messageRepo := repository.NewMessageRepository(cfg.DB)
	conversationRepo := repository.NewConversationRepository(cfg.DB)
	redisClient := repository.NewRedisClient(cfg.RedisURL)
	defer redisClient.Close()

	// Initialize services
	messageService := service.NewMessageService(messageRepo)
	conversationService := service.NewConversationService(conversationRepo)
	presenceService := service.NewPresenceService(redisClient)

	hub := websocket.NewHub(messageService, presenceService)
	go hub.Run()

	h := handler.New(cfg, messageService, conversationService, hub)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.HealthCheck)
	mux.HandleFunc("/ws", h.WebSocketHandler)

	// Protected routes - require JWT authentication
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)
	mux.Handle("/conversations", authMiddleware(http.HandlerFunc(h.ConversationsHandler)))
	mux.Handle("/conversations/", authMiddleware(http.HandlerFunc(h.ConversationHandler)))
	mux.Handle("/messages", authMiddleware(http.HandlerFunc(h.SendMessageHandler)))

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
