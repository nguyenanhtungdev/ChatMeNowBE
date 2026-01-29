# Development Guide - ChatMeNow

## 📁 Project Structure

```
ChatMeNow/
├── gateway/              # NestJS API Gateway (Port 3000)
│   ├── src/
│   │   ├── main.ts
│   │   ├── app.module.ts
│   │   ├── controllers/
│   │   ├── services/
│   │   └── guards/
│   ├── package.json
│   └── Dockerfile
│
├── auth-service/         # NestJS Auth Service (Port 3001)
│   ├── src/
│   │   ├── main.ts
│   │   ├── app.module.ts
│   │   ├── auth.controller.ts
│   │   ├── auth.service.ts
│   │   ├── entities/
│   │   ├── dto/
│   │   └── guards/
│   ├── package.json
│   └── Dockerfile
│
├── blog-service/         # NestJS Blog Service (Port 3002)
│   ├── src/
│   │   ├── main.ts
│   │   ├── app.module.ts
│   │   ├── post.controller.ts
│   │   ├── post.service.ts
│   │   ├── entities/
│   │   ├── dto/
│   │   └── guards/
│   ├── package.json
│   └── Dockerfile
│
├── chat-service/         # Go Chat Service (Port 8080)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/
│   │   └── websocket/
│   ├── go.mod
│   └── Dockerfile
│
├── docker-compose.yml
├── init-db.sql
└── .env.example
```

## 🔧 Development Setup

### Prerequisites

- Docker & Docker Compose
- Node.js 18+ (for local dev)
- Go 1.21+ (for local dev)
- PostgreSQL client (optional)
- MongoDB client (optional)

### Local Development (without Docker)

#### 1. Start Databases

```bash
# Start only databases
docker-compose up -d postgres mongodb redis
```

#### 2. Gateway Development

```bash
cd gateway
npm install
npm run start:dev
# Runs on http://localhost:3000
```

#### 3. Auth Service Development

```bash
cd auth-service
npm install
npm run start:dev
# Runs on http://localhost:3001
```

#### 4. Blog Service Development

```bash
cd blog-service
npm install
npm run start:dev
# Runs on http://localhost:3002
```

#### 5. Chat Service Development

```bash
cd chat-service
go mod download
go run cmd/server/main.go
# Runs on http://localhost:8080
```

## 🎯 Key Concepts

### Authentication Flow

1. **Register/Login** → Returns `accessToken` (15min) + `refreshToken` (7 days)
2. **Authenticated Requests** → Include `Authorization: Bearer {accessToken}` header
3. **Token Refresh** → Use refresh token to get new access token
4. **Logout** → Invalidate refresh token

### WebSocket Protocol

**Client → Server:**

```javascript
// Join conversation
{ type: "join_conversation", payload: { conversationId: "..." } }

// Send message
{ type: "send_message", payload: { conversationId: "...", content: "..." } }

// Typing indicator
{ type: "typing", payload: { conversationId: "...", isTyping: true/false } }
```

**Server → Client:**

```javascript
// New message
{ type: "new_message", payload: { id, conversationId, senderId, content, createdAt } }

// User typing
{ type: "user_typing", payload: { conversationId, userId, isTyping } }
```

### Database Models

**PostgreSQL** - Relational data

- `users` - User accounts
- `refresh_tokens` - Session management
- `conversations` - Chat rooms/DMs
- `conversation_members` - Many-to-many relationship
- `blog_posts` - Blog content

**MongoDB** - Chat messages (high-write throughput)

- `messages` collection with indexes on `conversationId` + `createdAt`

**Redis** - Ephemeral data

- `user:{userId}:online` - Online presence
- `conversation:{convId}:typing` - Typing indicators

## 🚀 Deployment Checklist

### Environment Variables

**Production `.env`:**

```bash
# Change these!
JWT_SECRET=<strong-random-secret>
POSTGRES_PASSWORD=<strong-password>
MONGO_INITDB_ROOT_PASSWORD=<strong-password>

# Rate limiting
RATE_LIMIT_TTL=60
RATE_LIMIT_MAX=100

# Token expiry
JWT_EXPIRES_IN=15m
REFRESH_TOKEN_EXPIRES_IN=7d
```

### Security Best Practices

1. **HTTPS Only** in production
2. **Strong JWT Secret** (64+ characters)
3. **CORS Configuration** - whitelist allowed origins
4. **Rate Limiting** - prevent abuse
5. **Input Validation** - all DTOs validated
6. **SQL Injection** - using parameterized queries
7. **XSS Protection** - sanitize user input
8. **Password Hashing** - bcrypt with salt

### Database Migrations

```bash
# PostgreSQL migrations (use TypeORM or raw SQL)
cd auth-service
npm run migration:generate -- -n CreateUsersTable
npm run migration:run

# MongoDB indexes (created in repository code)
# See chat-service/internal/repository/message.go
```

### Monitoring

**Health Checks:**

- Gateway: `GET http://localhost:3000/health`
- Auth: `GET http://localhost:3001/auth/health`
- Blog: `GET http://localhost:3002/posts/health`
- Chat: `GET http://localhost:8080/health`

**Logging:**

```bash
# Aggregate logs
docker-compose logs -f --tail=100

# Service-specific
docker-compose logs -f gateway
docker-compose logs -f chat-service
```

**Database Backups:**

```bash
# PostgreSQL
docker exec chatmenow-postgres pg_dump -U chatmenow chatmenow > backup.sql

# MongoDB
docker exec chatmenow-mongodb mongodump --username=chatmenow --password=chatmenow123 --out=/backup

# Restore
docker exec -i chatmenow-postgres psql -U chatmenow chatmenow < backup.sql
```

## 🧪 Testing

### Unit Tests

```bash
# NestJS services
cd auth-service
npm test

# Go services
cd chat-service
go test ./...
```

### Integration Tests

```bash
# Run test script
chmod +x test-api.sh
./test-api.sh
```

### Load Testing (with k6)

```javascript
// load-test.js
import ws from "k6/ws";
import { check } from "k6";

export let options = {
  vus: 100, // 100 concurrent users
  duration: "30s",
};

export default function () {
  const url = "ws://localhost:8080/ws?token=YOUR_TOKEN";

  const res = ws.connect(url, function (socket) {
    socket.on("open", () => {
      socket.send(
        JSON.stringify({
          type: "send_message",
          payload: { conversationId: "test", content: "Hello" },
        }),
      );
    });

    socket.on("message", (data) => console.log(data));
    socket.setTimeout(() => socket.close(), 5000);
  });

  check(res, { "status is 101": (r) => r && r.status === 101 });
}
```

```bash
k6 run load-test.js
```

## 📚 Adding Features

### Add New REST Endpoint (NestJS)

1. Create DTO: `src/dto/feature.dto.ts`
2. Update Service: `src/feature.service.ts`
3. Add Controller method: `src/feature.controller.ts`
4. Test: `curl http://localhost:3000/api/...`

### Add New WebSocket Event (Go)

1. Update model: `internal/model/model.go`
2. Add handler case: `internal/websocket/hub.go` → `HandleClientMessage()`
3. Broadcast: `hub.BroadcastToConversation(...)`

### Add New Database Table

1. Update SQL: `init-db.sql`
2. Create Entity (NestJS): `src/entities/table.entity.ts`
3. Register in Module: `TypeOrmModule.forFeature([NewEntity])`
4. Rebuild: `docker-compose down -v && docker-compose up -d`

## 🐛 Common Issues

### Port Already in Use

```bash
# Find process
sudo lsof -i :3000
kill -9 <PID>
```

### Database Connection Failed

```bash
# Check if containers are running
docker-compose ps

# Restart database
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### WebSocket Connection Refused

```bash
# Check if chat-service is running
docker-compose logs chat-service

# Verify JWT token is valid
# Token expires after 15 minutes
```

### TypeScript Errors (node_modules not found)

```bash
cd gateway  # or auth-service, blog-service
rm -rf node_modules package-lock.json
npm install
```

### Go Module Errors

```bash
cd chat-service
go mod tidy
go mod download
```

## 📞 Support

- GitHub Issues: [Report bugs]
- Documentation: See README.md
- API Testing: See API_TESTING.md

Happy coding! 🚀
