#!/bin/bash

# ChatMeNow Quick Start Script

echo "Starting ChatMeNow..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Please start Docker first."
    exit 1
fi

# Copy env file if not exists
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
fi

# Start services
echo "Starting Docker containers..."
docker-compose up -d

# Wait for databases to be ready
echo "Waiting for databases to be ready..."
sleep 10

echo ""
echo "ChatMeNow is running!"
echo ""
echo "📍 Services:"
echo "  Gateway:       http://localhost:3000"
echo "  Auth Service:  http://localhost:3001"
echo "  Blog Service:  http://localhost:3002"
echo "  Chat Service:  http://localhost:8080"
echo "  PostgreSQL:    localhost:5432"
echo "  MongoDB:       localhost:27017"
echo "  Redis:         localhost:6379"
echo ""
echo "Check logs: docker-compose logs -f"
echo "Stop:       docker-compose down"
echo ""
