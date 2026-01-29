#!/bin/bash

echo "=== Quick Conversation Creator ==="
echo ""
read -p "Enter User 1 Email: " EMAIL1
read -sp "Enter User 1 Password: " PASS1
echo ""
read -p "Enter User 2 Email: " EMAIL2
read -sp "Enter User 2 Password: " PASS2
echo ""
echo ""

# Login User 1
echo "=== Logging in User 1... ==="
LOGIN1=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL1\",
    \"password\": \"$PASS1\"
  }")

TOKEN1=$(echo "$LOGIN1" | jq -r '.accessToken')
USER1_ID=$(echo "$LOGIN1" | jq -r '.user.id')

if [ "$TOKEN1" == "null" ] || [ -z "$TOKEN1" ]; then
  echo "❌ Login User 1 failed!"
  echo "$LOGIN1" | jq '.'
  exit 1
fi

echo "✅ User 1 logged in"
echo "   ID: $USER1_ID"

# Login User 2
echo ""
echo "=== Logging in User 2... ==="
LOGIN2=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$EMAIL2\",
    \"password\": \"$PASS2\"
  }")

TOKEN2=$(echo "$LOGIN2" | jq -r '.accessToken')
USER2_ID=$(echo "$LOGIN2" | jq -r '.user.id')

if [ "$TOKEN2" == "null" ] || [ -z "$TOKEN2" ]; then
  echo "❌ Login User 2 failed!"
  echo "$LOGIN2" | jq '.'
  exit 1
fi

echo "✅ User 2 logged in"
echo "   ID: $USER2_ID"

# Create conversation
echo ""
echo "=== Creating Conversation... ==="
CONV=$(curl -s -X POST http://localhost:8080/api/conversations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN1" \
  -d "{
    \"participantIds\": [\"$USER2_ID\"],
    \"type\": \"direct\",
    \"name\": \"Chat Test\"
  }")

CONV_ID=$(echo "$CONV" | jq -r '.id')

if [ "$CONV_ID" == "null" ] || [ -z "$CONV_ID" ]; then
  echo "❌ Create conversation failed!"
  echo "$CONV" | jq '.'
  exit 1
fi

echo "✅ Conversation created!"
echo ""
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║                  🎉 SUCCESS! COPY THESE VALUES 🎉            ║"
echo "╠══════════════════════════════════════════════════════════════╣"
echo "║ Conversation ID:                                             ║"
echo "║ $CONV_ID"
echo "║                                                              ║"
echo "║ User 1 ($EMAIL1) Token:                        ║"
echo "║ $TOKEN1"
echo "║                                                              ║"
echo "║ User 2 ($EMAIL2) Token:                        ║"
echo "║ $TOKEN2"
echo "╠══════════════════════════════════════════════════════════════╣"
echo "║                     HOW TO TEST                              ║"
echo "╠══════════════════════════════════════════════════════════════╣"
echo "║ 1. Open http://localhost:5500/websocket-test.html           ║"
echo "║    in NORMAL Chrome window:                                  ║"
echo "║    - Paste User 1 Token                                      ║"
echo "║    - Paste Conversation ID                                   ║"
echo "║    - Click Connect                                           ║"
echo "║                                                              ║"
echo "║ 2. Open http://localhost:5500/websocket-test.html           ║"
echo "║    in INCOGNITO Chrome window:                               ║"
echo "║    - Paste User 2 Token                                      ║"
echo "║    - Paste Conversation ID (SAME as above)                   ║"
echo "║    - Click Connect                                           ║"
echo "║                                                              ║"
echo "║ 3. Start chatting! Messages will sync in real-time! 🚀      ║"
echo "╚══════════════════════════════════════════════════════════════╝"
