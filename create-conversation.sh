#!/bin/bash

# Script to create a conversation between existing users

# User 1 (alice3)
TOKEN1="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI5MWQzNmU4Yy01MmYwLTQ2MzAtODk5Ny01N2ZmMDkyYWQ5ODQiLCJ1c2VybmFtZSI6ImFsaWNlMyIsImVtYWlsIjoiYWxpY2UzQGV4YW1wbGUuY29tIiwiaWF0IjoxNzY5Njc0MTU0LCJleHAiOjE3Njk2NzUwNTR9.mAhRc_SKyz4Qs3XIc6QxsiEsOKtZZp1dUQ2ky8G2PlA"
USER1_ID="91d36e8c-52f0-4630-8997-57ff092ad984"

# User 2 (alice2)
TOKEN2="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzYzdjOTBjZC04ZDFiLTRhNTItOTZmNS0xZjliOTgyNmIzZjUiLCJ1c2VybmFtZSI6ImFsaWNlMiIsImVtYWlsIjoiYWxpY2UyQGV4YW1wbGUuY29tIiwiaWF0IjoxNzY5Njc0MDMwLCJleHAiOjE3Njk2NzQ5MzB9.x1PxViaLz5QUVUYw569qAUR-P7qM_aLFQc6y439inKA"
USER2_ID="3c7c90cd-8d1b-4a52-96f5-1f9b9826b3f5"

echo "Alice Token: $TOKEN1"
echo "Alice ID: $USER1_ID"

# Login User 2
echo -e "\n=== Login User 2 (Bob) ==="
LOGIN2=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "bob@example.com",
    "password": "password123"
  }')

echo "$LOGIN2" | jq '.'
TOKEN2=$(echo "$LOGIN2" | jq -r '.accessToken')
USER2_ID=$(echo "$LOGIN2" | jq -r '.user.id')

echo "Bob Token: $TOKEN2"
echo "Bob ID: $USER2_ID"

# Create conversation
echo -e "\n=== Creating Conversation ==="
CONV_RESPONSE=$(curl -s -X POST http://localhost:8080/api/conversations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN1" \
  -d "{
    \"participantIds\": [\"$USER2_ID\"],
    \"type\": \"direct\",
    \"name\": \"Alice and Bob Chat\"
  }")

echo "$CONV_RESPONSE" | jq '.'
CONV_ID=$(echo "$CONV_RESPONSE" | jq -r '.id')

echo -e "\n╔════════════════════════════════════════════════════════════╗"
echo "║                    WEBSOCKET TEST INFO                     ║"
echo "╠════════════════════════════════════════════════════════════╣"
echo "║ Conversation ID:                                           ║"
echo "║ $CONV_ID"
echo "║                                                            ║"
echo "║ Alice (User 1) Token:                                      ║"
echo "║ $TOKEN1"
echo "║                                                            ║"
echo "║ Bob (User 2) Token:                                        ║"
echo "║ $TOKEN2"
echo "╠════════════════════════════════════════════════════════════╣"
echo "║                   TESTING INSTRUCTIONS                     ║"
echo "╠════════════════════════════════════════════════════════════╣"
echo "║ 1. Open websocket-test.html in Chrome (normal mode)       ║"
echo "║    - Paste Alice's token                                   ║"
echo "║    - Paste Conversation ID                                 ║"
echo "║    - Click Connect                                         ║"
echo "║                                                            ║"
echo "║ 2. Open websocket-test.html in Chrome (incognito mode)    ║"
echo "║    - Paste Bob's token                                     ║"
echo "║    - Paste the SAME Conversation ID                        ║"
echo "║    - Click Connect                                         ║"
echo "║                                                            ║"
echo "║ 3. Start chatting between the two windows!                ║"
echo "╚════════════════════════════════════════════════════════════╝"
