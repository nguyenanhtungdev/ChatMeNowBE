# Changelog

All notable changes to ChatMeNow will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-01-29

### Added

#### Gateway Service

- API Gateway with request proxying
- JWT authentication middleware
- Rate limiting (100 req/min)
- Request logging with request IDs
- CORS support
- Health check endpoint

#### Auth Service

- User registration with bcrypt password hashing
- Login with JWT token generation
- Refresh token mechanism (7-day expiry)
- Device session tracking
- User profile endpoint
- Logout functionality
- TypeORM integration with PostgreSQL

#### Blog Service

- CRUD operations for blog posts
- Draft/Published/Archived status
- Tags support
- View count tracking
- User-specific post filtering
- Publish/unpublish functionality
- TypeORM integration with PostgreSQL

#### Chat Service (Go)

- WebSocket real-time messaging
- Conversation management (direct/group)
- Message persistence to MongoDB
- Online presence tracking with Redis
- Typing indicators
- Join/leave conversation rooms
- Message history retrieval
- Concurrent connection handling
- Graceful shutdown

#### Infrastructure

- Docker Compose orchestration
- PostgreSQL 15 database
- MongoDB 7 document store
- Redis 7 cache
- Database initialization scripts
- Health checks for all services
- Volume persistence

#### Developer Tools

- Complete API documentation
- WebSocket test UI (HTML)
- Automated API test script
- Quick start script
- Makefile with useful commands
- Development mode support
- Database backup script

#### Documentation

- README.md - Project overview
- QUICK_START.md - Quick start guide
- DEVELOPMENT.md - Development guide
- API_TESTING.md - API documentation
- PROJECT_SUMMARY.md - Complete summary
- CONTRIBUTING.md - Contribution guide

### Security

- JWT-based authentication
- Bcrypt password hashing (10 rounds)
- Refresh token rotation
- Rate limiting
- Input validation with class-validator
- SQL injection prevention
- Environment-based configuration

## [Unreleased]

### Planned Features

- End-to-end message encryption
- File upload support (images, documents)
- Voice/Video calling (WebRTC)
- Push notifications (FCM)
- Message reactions
- Read receipts
- User blocking
- Full-text message search
- Group admin controls
- Kubernetes deployment
- CI/CD pipeline
- Prometheus monitoring
- E2E tests

---

[1.0.0]: https://github.com/chatmenow/chatmenow/releases/tag/v1.0.0
