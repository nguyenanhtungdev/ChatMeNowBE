# ChatMeNow - Realtime Chat Platform

Ứng dụng chat real-time kiểu Zalo với kiến trúc microservices hybrid **Go + Node.js**.

## 🏗️ Architecture

```
┌─────────┐
│ Client  │
└────┬────┘
     │
     v
┌──────────────┐
│   Gateway    │  (NestJS - Port 3000)
│  + JWT Auth  │
│  + RateLimit │
└──────┬───────┘
       │
       ├─────────────────┬─────────────────┐
       v                 v                 v
┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│Auth Service │   │Blog Service │   │Chat Service │
│  (NestJS)   │   │  (NestJS)   │   │    (Go)     │
│  Port 3001  │   │  Port 3002  │   │  Port 8080  │
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │                 │                 │
       v                 v                 v
  PostgreSQL         PostgreSQL         MongoDB
                                      + Redis (presence)
```

## 📦 Services

### 1. Gateway (NestJS)

- Public API cho client
- JWT verification
- Rate limiting
- Request logging & tracing
- Proxy requests to services

### 2. Auth Service (NestJS + PostgreSQL)

- User registration/login
- JWT token generation
- Refresh token rotation
- Device session management
- User profile

### 3. Blog Service (NestJS + PostgreSQL)

- CRUD blog posts
- Publish/Draft status
- Tags & Categories
- User posts

### 4. Chat Service (Go + MongoDB + Redis)

- WebSocket real-time messaging
- REST API for message history
- Conversation management
- Online presence (Redis)
- Typing indicators

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Node.js 18+ (cho development)
- Go 1.21+ (cho development)

### Run với Docker Compose

```bash
# Clone và start tất cả services
docker-compose up -d

# Check logs
docker-compose logs -f

# Stop
docker-compose down
```

### API Endpoints

#### Gateway (`:3000`)

```
GET  /health
POST /api/auth/*      → auth-service
POST /api/blog/*      → blog-service
POST /api/chat/*      → chat-service
WS   /ws              → chat-service (WebSocket)
```

#### Auth Service (`:3001`)

```
POST /auth/register
POST /auth/login
POST /auth/refresh
GET  /auth/me
POST /auth/logout
```

#### Blog Service (`:3002`)

```
GET    /posts
POST   /posts
GET    /posts/:id
PUT    /posts/:id
DELETE /posts/:id
PATCH  /posts/:id/publish
```

#### Chat Service (`:8080`)

```
GET  /conversations
POST /conversations
GET  /conversations/:id/messages
POST /messages
WS   /ws?token=<jwt>
```

## 📊 Database Schema

### PostgreSQL (Auth & Blog)

**users**

```sql
id, username, email, password_hash, avatar_url, created_at, updated_at
```

**refresh_tokens**

```sql
id, user_id, token_hash, device_id, expires_at, created_at
```

**conversations**

```sql
id, name, type (direct/group), created_by, created_at, updated_at
```

**conversation_members**

```sql
id, conversation_id, user_id, role, joined_at
```

**blog_posts**

```sql
id, user_id, title, content, status, tags, created_at, updated_at
```

### MongoDB (Chat)

**messages**

```javascript
{
  _id: ObjectId,
  conversationId: string,
  senderId: string,
  content: string,
  type: "text" | "image" | "file",
  metadata: {},
  createdAt: Date,
  updatedAt: Date
}
```

### Redis (Presence)

```
user:{userId}:online → timestamp
conversation:{convId}:typing → Set of userIds
```

## 🔐 Authentication Flow

1. Client → `POST /api/auth/register` → JWT access + refresh token
2. Client → `POST /api/auth/login` → JWT tokens
3. Client → `GET /api/chat/conversations` (Authorization: Bearer {token})
4. Gateway verify JWT → Forward to services
5. Client → WebSocket `/ws?token={jwt}` → Realtime

## 📝 WebSocket Protocol

### Client → Server

```json
{
  "type": "join_conversation",
  "payload": { "conversationId": "123" }
}

{
  "type": "send_message",
  "payload": {
    "conversationId": "123",
    "content": "Hello!"
  }
}

{
  "type": "typing",
  "payload": { "conversationId": "123", "isTyping": true }
}
```

### Server → Client

```json
{
  "type": "new_message",
  "payload": {
    "id": "msg123",
    "conversationId": "123",
    "senderId": "user456",
    "content": "Hello!",
    "createdAt": "2026-01-29T..."
  }
}

{
  "type": "user_typing",
  "payload": {
    "conversationId": "123",
    "userId": "user789",
    "isTyping": true
  }
}
```

## 🛠️ Development

### Gateway

```bash
cd gateway
npm install
npm run start:dev
```

### Auth Service

```bash
cd auth-service
npm install
npm run start:dev
```

### Blog Service

```bash
cd blog-service
npm install
npm run start:dev
```

### Chat Service

```bash
cd chat-service
go mod download
go run cmd/server/main.go
```

## 🧪 Testing

```bash
# Test auth
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"123456"}'

# Test WebSocket (sử dụng wscat)
npm install -g wscat
wscat -c "ws://localhost:8080/ws?token=YOUR_JWT_TOKEN"
```

## 📂 Project Structure

```
ChatMeNow/
├── gateway/              # NestJS API Gateway
├── auth-service/         # NestJS Auth Service
├── blog-service/         # NestJS Blog Service
├── chat-service/         # Go Chat Service
├── docker-compose.yml    # Orchestration
├── .env.example          # Environment variables
└── README.md
```

## 🔮 Future Enhancements

- [ ] Message encryption (E2E)
- [ ] File upload (S3/MinIO)
- [ ] Voice/Video call (WebRTC)
- [ ] Push notifications (FCM)
- [ ] Message reactions
- [ ] Read receipts
- [ ] Group admin features
- [ ] User blocking
- [ ] Search messages
- [ ] Elasticsearch for full-text search

## 📄 License

MIT
