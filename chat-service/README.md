# Chat Service (Go + GORM)

Real-time messaging service built with **Go**, **GORM**, **WebSocket**, và **PostgreSQL**.

## 🚀 Features

- ✅ **GORM ORM** - Type-safe database queries với PostgreSQL
- ✅ **UUID** - Sử dụng UUID cho tất cả các ID
- ✅ **WebSocket** - Real-time messaging
- ✅ **Redis** - Presence & typing indicators
- ✅ **Auto-migration** - Tự động tạo bảng khi khởi động
- ✅ **JWT Authentication** - Secure authentication
- ✅ **RESTful API** - CRUD operations

## 📦 Tech Stack

```
- Go 1.21
- GORM (PostgreSQL ORM)
- Gorilla WebSocket
- UUID v4
- Redis (presence)
- PostgreSQL
```

## 🗄️ Database Schema

### Tables

#### `conversations`

```sql
id UUID PRIMARY KEY
name VARCHAR(255)
type VARCHAR(20)  -- 'direct' | 'group'
avatar_url VARCHAR(500)
created_by UUID
created_at TIMESTAMP
updated_at TIMESTAMP
deleted_at TIMESTAMP
```

#### `conversation_members`

```sql
id UUID PRIMARY KEY
conversation_id UUID REFERENCES conversations(id)
user_id UUID
role VARCHAR(20)  -- 'admin' | 'member'
joined_at TIMESTAMP
deleted_at TIMESTAMP
```

#### `messages`

```sql
id UUID PRIMARY KEY
conversation_id UUID REFERENCES conversations(id)
sender_id UUID
content TEXT
type VARCHAR(20)  -- 'text' | 'image' | 'file' | 'video'
metadata JSONB
created_at TIMESTAMP
updated_at TIMESTAMP
deleted_at TIMESTAMP
```

## 🏗️ Project Structure

```
chat-service/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # GORM connection setup
│   ├── model/
│   │   └── model.go          # GORM models với UUID
│   ├── repository/
│   │   ├── message.go        # Message GORM repo
│   │   ├── conversation.go   # Conversation GORM repo
│   │   └── redis.go          # Redis client
│   ├── service/
│   │   ├── message.go        # Business logic
│   │   ├── conversation.go
│   │   └── presence.go
│   ├── handler/
│   │   └── handler.go        # HTTP handlers
│   ├── websocket/
│   │   ├── hub.go            # WebSocket hub
│   │   ├── client.go         # WebSocket client
│   │   └── register.go       # Connection registry
│   └── middleware/
│       ├── auth.go           # JWT middleware
│       └── jwt.go            # JWT utils
├── migrations/
│   └── 001_init_schema.sql  # SQL migrations
└── go.mod
```

## 🔧 Installation

1. **Install dependencies**:

```bash
cd chat-service
go mod tidy
```

2. **Set environment variables**:

```bash
export PORT=8080
export POSTGRES_URL="postgresql://chatmenow:chatmenow123@localhost:5432/chatmenow"
export REDIS_URL="localhost:6379"
export JWT_SECRET="your-secret-key"
```

3. **Run the service**:

```bash
go run cmd/server/main.go
```

## 📝 API Endpoints

### REST API

#### Create Conversation

```http
POST /conversations
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "name": "Team Chat",
  "type": "group",
  "memberIds": ["uuid1", "uuid2", "uuid3"]
}
```

#### Get User Conversations

```http
GET /conversations
Authorization: Bearer <JWT>
```

#### Get Conversation Messages

```http
GET /conversations/{id}/messages?limit=50&offset=0
Authorization: Bearer <JWT>
```

#### Send Message

```http
POST /messages
Authorization: Bearer <JWT>
Content-Type: application/json

{
  "conversationId": "uuid",
  "content": "Hello!",
  "type": "text"
}
```

### WebSocket

#### Connect

```javascript
const ws = new WebSocket("ws://localhost:8080/ws?token=JWT_TOKEN");

// Join room
ws.send(
  JSON.stringify({
    type: "join",
    payload: {
      conversationId: "uuid",
      userId: "uuid",
    },
  }),
);

// Send message
ws.send(
  JSON.stringify({
    type: "message",
    payload: {
      conversationId: "uuid",
      content: "Hello!",
      type: "text",
    },
  }),
);

// Typing indicator
ws.send(
  JSON.stringify({
    type: "typing",
    payload: {
      conversationId: "uuid",
      userId: "uuid",
      isTyping: true,
    },
  }),
);
```

## 🔍 GORM Usage Examples

### Create Message

```go
message := &model.Message{
    ConversationID: conversationUUID,
    SenderID:       userUUID,
    Content:        "Hello, World!",
    Type:           "text",
}
err := db.Create(message).Error
```

### Query with Preload

```go
var conversation model.Conversation
db.Preload("Members").First(&conversation, "id = ?", id)
```

### Complex Query

```go
var messages []model.Message
db.Where("conversation_id = ?", conversationID).
   Order("created_at DESC").
   Limit(50).
   Offset(0).
   Find(&messages)
```

### Transaction

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Create conversation
    if err := tx.Create(conv).Error; err != nil {
        return err
    }

    // Add members
    for _, memberID := range memberIDs {
        member := &model.ConversationMember{
            ConversationID: conv.ID,
            UserID:         memberID,
        }
        if err := tx.Create(member).Error; err != nil {
            return err
        }
    }

    return nil
})
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/repository -run TestMessageRepository
```

## 🐳 Docker

```bash
# Build image
docker build -t chat-service .

# Run container
docker run -p 8080:8080 \
  -e POSTGRES_URL="postgresql://chatmenow:chatmenow123@postgres:5432/chatmenow" \
  -e REDIS_URL="redis:6379" \
  chat-service
```

## 📊 Performance Tips

1. **Use Indexes**: Đã tạo indexes cho các trường thường query
2. **Connection Pooling**: GORM tự động quản lý connection pool
3. **Prepared Statements**: Enable trong GORM config
4. **Pagination**: Luôn sử dụng `LIMIT` và `OFFSET`

## 🔒 Security

- JWT authentication cho tất cả endpoints
- UUID thay vì auto-increment ID
- Soft delete với `deleted_at`
- Input validation
- CORS middleware

## 📚 References

- [GORM Documentation](https://gorm.io/docs/)
- [UUID Package](https://github.com/google/uuid)
- [Gorilla WebSocket](https://github.com/gorilla/websocket)

## 🤝 Contributing

Contributions welcome! Please read CONTRIBUTING.md first.

## 📄 License

MIT License - see LICENSE file for details.
