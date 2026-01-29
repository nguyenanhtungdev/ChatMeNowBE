# ChatMeNow API Testing Guide

## 1. Register a New User

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "password123"
  }'
```

Response:

```json
{
  "message": "Registration successful. Please login to continue.",
  "user": {
    "id": "uuid",
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

**Note**: You must login after registration to get access tokens.

## 2. Login

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "emailOrUsername": "alice@example.com",
    "password": "password123"
  }'
```

Response:

```json
{
  "accessToken": "eyJhbGc...",
  "refreshToken": "eyJhbGc...",
  "user": {
    "id": "uuid",
    "username": "alice",
    "email": "alice@example.com"
  }
}
```

**Note**: Save the `accessToken` to use in subsequent API calls.

## 3. Get Current User

```bash
curl http://localhost:3000/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 4. Create a Conversation

### Direct Conversation (1-on-1 chat)

```bash
curl -X POST http://localhost:3000/api/chat/conversations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "direct",
    "memberIds": ["other-user-uuid"]
  }'
```

**Note**: For direct conversation, only pass 1 member ID (the other person). The system automatically adds you, making it 2 members total.

### Group Conversation

```bash
curl -X POST http://localhost:3000/api/chat/conversations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "General Chat",
    "type": "group",
    "memberIds": ["user-id-1", "user-id-2", "user-id-3"]
  }'
```

**Note**: For group conversation, `name` is required. You can add any number of members.

## 5. Get All Conversations

```bash
curl http://localhost:3000/api/chat/conversations \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 6. Send a Message (REST)

```bash
curl -X POST http://localhost:3000/api/chat/messages \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "conversationId": "conv-uuid",
    "content": "Hello World!",
    "type": "text"
  }'
```

## 7. Get Messages from Conversation

```bash
curl http://localhost:3000/api/chat/conversations/CONV_ID/messages \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 8. Create a Blog Post

```bash
curl -X POST http://localhost:3000/api/blog/posts \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Post",
    "content": "This is my first blog post!",
    "excerpt": "Introduction to my blog",
    "tags": ["tutorial", "tech"]
  }'
```

## 9. Get All Published Posts

```bash
curl http://localhost:3000/api/blog/posts
```

## 10. Publish a Post

```bash
curl -X PATCH http://localhost:3000/api/blog/posts/POST_ID/publish \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## WebSocket Connection

### Using JavaScript

```javascript
const token = "YOUR_ACCESS_TOKEN";
const ws = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

ws.onopen = () => {
  console.log("Connected");

  // Join a conversation
  ws.send(
    JSON.stringify({
      type: "join_conversation",
      payload: { conversationId: "your-conversation-id" },
    }),
  );

  // Send a message
  ws.send(
    JSON.stringify({
      type: "send_message",
      payload: {
        conversationId: "your-conversation-id",
        content: "Hello from WebSocket!",
      },
    }),
  );

  // Start typing
  ws.send(
    JSON.stringify({
      type: "typing",
      payload: {
        conversationId: "your-conversation-id",
        isTyping: true,
      },
    }),
  );
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log("Received:", data);
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

ws.onclose = () => {
  console.log("Disconnected");
};
```

### Using wscat (CLI)

```bash
# Install wscat
npm install -g wscat

# Connect
wscat -c "ws://localhost:8080/ws?token=YOUR_ACCESS_TOKEN"

# Then send messages:
{"type":"join_conversation","payload":{"conversationId":"conv-123"}}
{"type":"send_message","payload":{"conversationId":"conv-123","content":"Hello!"}}
{"type":"typing","payload":{"conversationId":"conv-123","isTyping":true}}
```

## Database Access

### PostgreSQL

```bash
docker exec -it chatmenow-postgres psql -U chatmenow -d chatmenow

# Useful queries
SELECT * FROM users;
SELECT * FROM conversations;
SELECT * FROM blog_posts;
```

### MongoDB

```bash
docker exec -it chatmenow-mongodb mongosh -u chatmenow -p chatmenow123

use chatmenow
db.messages.find().pretty()
db.messages.find({conversationId: "your-id"}).sort({createdAt: -1}).limit(10)
```

### Redis

```bash
docker exec -it chatmenow-redis redis-cli

# Check online users
KEYS user:*:online

# Check typing users
SMEMBERS conversation:conv-123:typing
```

## Full User Flow Example

```bash
# 1. Register two users
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@test.com","password":"123456"}'

curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","email":"bob@test.com","password":"123456"}'

# Save the tokens from responses

# 2. Alice creates a conversation
TOKEN_ALICE="eyJhbGc..."
curl -X POST http://localhost:3000/api/chat/conversations \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice & Bob","type":"direct","memberIds":["alice-id","bob-id"]}'

# 3. Alice sends a message
CONV_ID="conversation-uuid"
curl -X POST http://localhost:3000/api/chat/messages \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d "{\"conversationId\":\"$CONV_ID\",\"content\":\"Hi Bob!\",\"type\":\"text\"}"

# 4. Bob gets messages
TOKEN_BOB="eyJhbGc..."
curl http://localhost:3000/api/chat/conversations/$CONV_ID/messages \
  -H "Authorization: Bearer $TOKEN_BOB"

# 5. Connect via WebSocket for real-time
# Open two browser tabs with the JavaScript code above
# Use TOKEN_ALICE in one, TOKEN_BOB in another
# Send messages and see them appear in real-time!
```

## Troubleshooting

### Check service status

```bash
docker-compose ps
```

### View logs

```bash
docker-compose logs -f gateway
docker-compose logs -f auth-service
docker-compose logs -f chat-service
```

### Restart a service

```bash
docker-compose restart gateway
```

### Rebuild after code changes

```bash
docker-compose up -d --build
```

### Clean everything

```bash
docker-compose down -v
docker-compose up -d
```
