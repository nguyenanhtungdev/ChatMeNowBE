#!/bin/bash

echo "🧪 Testing ChatMeNow APIs..."

API_BASE="http://localhost:3000"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Register user 1
echo ""
echo "1️⃣  Registering user Alice..."
RESPONSE=$(curl -s -X POST "$API_BASE/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"alice_test","email":"alice_test@example.com","password":"123456"}')

if echo "$RESPONSE" | grep -q "accessToken"; then
    echo -e "${GREEN}✅ Alice registered successfully${NC}"
    TOKEN_ALICE=$(echo "$RESPONSE" | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
    USER_ID_ALICE=$(echo "$RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
else
    echo -e "${RED}❌ Failed to register Alice${NC}"
    echo "$RESPONSE"
fi

# Register user 2
echo ""
echo "2️⃣  Registering user Bob..."
RESPONSE=$(curl -s -X POST "$API_BASE/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"username":"bob_test","email":"bob_test@example.com","password":"123456"}')

if echo "$RESPONSE" | grep -q "accessToken"; then
    echo -e "${GREEN}✅ Bob registered successfully${NC}"
    TOKEN_BOB=$(echo "$RESPONSE" | grep -o '"accessToken":"[^"]*' | cut -d'"' -f4)
    USER_ID_BOB=$(echo "$RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
else
    echo -e "${RED}❌ Failed to register Bob${NC}"
    echo "$RESPONSE"
fi

# Get current user
echo ""
echo "3️⃣  Getting Alice's profile..."
RESPONSE=$(curl -s "$API_BASE/api/auth/me" \
  -H "Authorization: Bearer $TOKEN_ALICE")

if echo "$RESPONSE" | grep -q "alice_test"; then
    echo -e "${GREEN}✅ Profile retrieved${NC}"
else
    echo -e "${RED}❌ Failed to get profile${NC}"
fi

# Create conversation
echo ""
echo "4️⃣  Alice creating a conversation..."
RESPONSE=$(curl -s -X POST "$API_BASE/api/chat/conversations" \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"Test Chat\",\"type\":\"group\",\"memberIds\":[\"$USER_ID_ALICE\",\"$USER_ID_BOB\"]}")

if echo "$RESPONSE" | grep -q '"id"'; then
    echo -e "${GREEN}✅ Conversation created${NC}"
    CONV_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
    echo "Conversation ID: $CONV_ID"
else
    echo -e "${RED}❌ Failed to create conversation${NC}"
    echo "$RESPONSE"
fi

# Send message
echo ""
echo "5️⃣  Alice sending a message..."
RESPONSE=$(curl -s -X POST "$API_BASE/api/chat/messages" \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d "{\"conversationId\":\"$CONV_ID\",\"content\":\"Hello from test!\",\"type\":\"text\"}")

if echo "$RESPONSE" | grep -q "Hello from test"; then
    echo -e "${GREEN}✅ Message sent${NC}"
else
    echo -e "${RED}❌ Failed to send message${NC}"
    echo "$RESPONSE"
fi

# Get messages
echo ""
echo "6️⃣  Bob getting messages..."
RESPONSE=$(curl -s "$API_BASE/api/chat/conversations/$CONV_ID/messages" \
  -H "Authorization: Bearer $TOKEN_BOB")

if echo "$RESPONSE" | grep -q "Hello from test"; then
    echo -e "${GREEN}✅ Messages retrieved${NC}"
else
    echo -e "${RED}❌ Failed to get messages${NC}"
    echo "$RESPONSE"
fi

# Create blog post
echo ""
echo "7️⃣  Alice creating a blog post..."
RESPONSE=$(curl -s -X POST "$API_BASE/api/blog/posts" \
  -H "Authorization: Bearer $TOKEN_ALICE" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Post","content":"This is a test post from API test script","tags":["test"]}')

if echo "$RESPONSE" | grep -q "Test Post"; then
    echo -e "${GREEN}✅ Blog post created${NC}"
    POST_ID=$(echo "$RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
else
    echo -e "${RED}❌ Failed to create blog post${NC}"
    echo "$RESPONSE"
fi

# Publish post
echo ""
echo "8️⃣  Publishing the blog post..."
RESPONSE=$(curl -s -X PATCH "$API_BASE/api/blog/posts/$POST_ID/publish" \
  -H "Authorization: Bearer $TOKEN_ALICE")

if echo "$RESPONSE" | grep -q "published"; then
    echo -e "${GREEN}✅ Post published${NC}"
else
    echo -e "${RED}❌ Failed to publish post${NC}"
fi

# Get all posts
echo ""
echo "9️⃣  Getting all published posts..."
RESPONSE=$(curl -s "$API_BASE/api/blog/posts")

if echo "$RESPONSE" | grep -q "Test Post"; then
    echo -e "${GREEN}✅ Posts retrieved${NC}"
else
    echo -e "${RED}❌ Failed to get posts${NC}"
fi

echo ""
echo "============================================"
echo -e "${GREEN}✅ All tests completed!${NC}"
echo "============================================"
echo ""
echo "📝 Saved credentials for manual testing:"
echo "Alice Token: $TOKEN_ALICE"
echo "Bob Token: $TOKEN_BOB"
echo "Conversation ID: $CONV_ID"
echo ""
echo "🔗 Test WebSocket connection:"
echo "wscat -c \"ws://localhost:8080/ws?token=$TOKEN_ALICE\""
echo ""
