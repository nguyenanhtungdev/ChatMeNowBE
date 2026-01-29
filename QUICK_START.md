# 🚀 Quick Start Guide

## Chạy Project trong 5 phút

### Bước 1: Clone và Prepare

```bash
cd /home/nguyenanhtung/Documents/ChatMeNow
cp .env.example .env
```

### Bước 2: Start với Docker Compose

```bash
# Cách 1: Dùng script
./start.sh

# Cách 2: Dùng Make
make start

# Cách 3: Dùng Docker Compose trực tiếp
docker-compose up -d
```

### Bước 3: Kiểm tra Services đang chạy

```bash
docker-compose ps

# Hoặc
make status
```

Bạn sẽ thấy:

```
chatmenow-gateway      running   0.0.0.0:3000->3000/tcp
chatmenow-auth         running   0.0.0.0:3001->3001/tcp
chatmenow-blog         running   0.0.0.0:3002->3002/tcp
chatmenow-chat         running   0.0.0.0:8080->8080/tcp
chatmenow-postgres     running   0.0.0.0:5432->5432/tcp
chatmenow-mongodb      running   0.0.0.0:27017->27017/tcp
chatmenow-redis        running   0.0.0.0:6379->6379/tcp
```

### Bước 4: Test API

```bash
# Chạy automated test
./test-api.sh

# Hoặc
make test
```

### Bước 5: Test WebSocket

1. Mở file `websocket-test.html` trong browser
2. Register user để lấy token:

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "123456"
  }'
```

3. Copy `accessToken` từ response
4. Paste vào WebSocket Test UI
5. Tạo conversation và chat!

## 🎯 Các Lệnh Hữu Ích

### Quản lý Services

```bash
make start      # Khởi động tất cả
make stop       # Dừng tất cả
make restart    # Restart tất cả
make logs       # Xem logs
make build      # Build lại images
make clean      # Xóa hết container + volume
```

### Xem Logs

```bash
# Tất cả services
make logs

# Service cụ thể
make logs SERVICE=gateway
make logs SERVICE=chat-service
docker-compose logs -f auth-service
```

### Kết nối Database

```bash
# PostgreSQL
make db-postgres
# Hoặc
docker exec -it chatmenow-postgres psql -U chatmenow -d chatmenow

# MongoDB
make db-mongo
# Hoặc
docker exec -it chatmenow-mongodb mongosh -u chatmenow -p chatmenow123

# Redis
make db-redis
# Hoặc
docker exec -it chatmenow-redis redis-cli
```

### Backup Database

```bash
make backup-db
```

## 📝 Test Flow Đầy Đủ

### 1. Register 2 users

```bash
# User 1: Alice
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "password": "123456"
  }'

# Lưu lại accessToken và user.id

# User 2: Bob
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "bob",
    "email": "bob@example.com",
    "password": "123456"
  }'
```

### 2. Alice tạo conversation

```bash
TOKEN_ALICE="<paste-token-alice>"
USER_ID_BOB="<paste-user-id-bob>"

curl -X POST http://localhost:3000/api/chat/conversations \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Alice & Bob Chat\",
    \"type\": \"direct\",
    \"memberIds\": [\"$USER_ID_BOB\"]
  }"

# Lưu lại conversation.id
```

### 3. Alice gửi message (REST)

```bash
CONV_ID="<paste-conversation-id>"

curl -X POST http://localhost:3000/api/chat/messages \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d "{
    \"conversationId\": \"$CONV_ID\",
    \"content\": \"Hello Bob!\",
    \"type\": \"text\"
  }"
```

### 4. Bob lấy messages

```bash
TOKEN_BOB="<paste-token-bob>"

curl http://localhost:3000/api/chat/conversations/$CONV_ID/messages \
  -H "Authorization: Bearer $TOKEN_BOB"
```

### 5. Test Real-time WebSocket

**Option 1: Dùng WebSocket Test UI**

1. Mở `websocket-test.html` trong 2 browser tab
2. Tab 1: Paste TOKEN_ALICE + CONV_ID → Connect
3. Tab 2: Paste TOKEN_BOB + CONV_ID → Connect
4. Gõ message ở Tab 1 → Thấy xuất hiện ngay ở Tab 2!

**Option 2: Dùng wscat (CLI)**

```bash
# Terminal 1 (Alice)
npm install -g wscat
wscat -c "ws://localhost:8080/ws?token=$TOKEN_ALICE"

# Sau khi connect, gửi:
{"type":"join_conversation","payload":{"conversationId":"YOUR_CONV_ID"}}
{"type":"send_message","payload":{"conversationId":"YOUR_CONV_ID","content":"Hi!"}}

# Terminal 2 (Bob)
wscat -c "ws://localhost:8080/ws?token=$TOKEN_BOB"
{"type":"join_conversation","payload":{"conversationId":"YOUR_CONV_ID"}}

# Bob sẽ nhận được message của Alice real-time!
```

## 🔧 Development Mode (Local)

Nếu muốn develop local (không dùng Docker):

### 1. Chỉ chạy databases

```bash
docker-compose up -d postgres mongodb redis
```

### 2. Chạy từng service

```bash
# Terminal 1: Gateway
cd gateway
npm install
npm run start:dev

# Terminal 2: Auth Service
cd auth-service
npm install
npm run start:dev

# Terminal 3: Blog Service
cd blog-service
npm install
npm run start:dev

# Terminal 4: Chat Service
cd chat-service
go mod download
go run cmd/server/main.go
```

Hoặc dùng Make:

```bash
make install          # Install tất cả dependencies
make dev-gateway      # Chạy gateway dev mode
make dev-auth         # Chạy auth-service dev mode
make dev-blog         # Chạy blog-service dev mode
make dev-chat         # Chạy chat-service dev mode
```

## 🐛 Troubleshooting

### Services không start được

```bash
# Check logs
docker-compose logs

# Xem service cụ thể bị lỗi gì
docker-compose logs gateway
docker-compose logs chat-service
```

### Port đã được dùng

```bash
# Tìm process đang dùng port
sudo lsof -i :3000
sudo lsof -i :8080

# Kill process
sudo kill -9 <PID>
```

### Database connection failed

```bash
# Restart databases
docker-compose restart postgres mongodb redis

# Hoặc recreate
docker-compose down
docker-compose up -d
```

### Clean start (xóa hết data)

```bash
# Cẩn thận: Sẽ xóa tất cả data!
docker-compose down -v
docker-compose up -d
```

## 📚 Tài liệu khác

- **README.md** - Tổng quan kiến trúc
- **API_TESTING.md** - Chi tiết tất cả API endpoints
- **DEVELOPMENT.md** - Hướng dẫn develop, add features
- **websocket-test.html** - WebSocket test UI

## 🎉 Chúc mừng!

Bạn đã setup thành công ChatMeNow!

Giờ có thể:

- ✅ Register/Login users
- ✅ Tạo conversations
- ✅ Gửi/nhận messages (REST)
- ✅ Real-time chat (WebSocket)
- ✅ Tạo blog posts
- ✅ Typing indicators
- ✅ Online presence

Happy coding! 🚀💬
