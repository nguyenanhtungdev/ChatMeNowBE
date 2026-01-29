# 🚀 GORM Quick Reference - ChatMeNow

## 📦 Import Packages

```go
import (
    "github.com/google/uuid"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)
```

---

## 1️⃣ Basic CRUD

### Create

```go
message := &model.Message{
    ConversationID: uuid.New(),
    SenderID:       uuid.New(),
    Content:        "Hello, World!",
    Type:           "text",
}
db.Create(message)
// message.ID được auto-generate
```

### Read

```go
// Find by ID
var message model.Message
db.First(&message, "id = ?", messageUUID)

// Find all
var messages []model.Message
db.Find(&messages)

// Find with condition
db.Where("sender_id = ?", userUUID).Find(&messages)
```

### Update

```go
// Update single field
db.Model(&message).Update("content", "Updated content")

// Update multiple fields
db.Model(&message).Updates(map[string]interface{}{
    "content": "New content",
    "type":    "edited",
})

// Save (update all fields)
message.Content = "Changed"
db.Save(&message)
```

### Delete

```go
// Soft delete (set deleted_at)
db.Delete(&message, "id = ?", messageUUID)

// Permanent delete
db.Unscoped().Delete(&message)
```

---

## 2️⃣ Advanced Queries

### Where Conditions

```go
// Simple where
db.Where("type = ?", "text").Find(&messages)

// Multiple conditions
db.Where("sender_id = ? AND type = ?", userUUID, "text").Find(&messages)

// IN clause
db.Where("id IN ?", []uuid.UUID{id1, id2, id3}).Find(&messages)

// Like
db.Where("content LIKE ?", "%hello%").Find(&messages)
```

### Order & Limit

```go
db.Order("created_at DESC").Limit(50).Find(&messages)

// With offset (pagination)
db.Order("created_at DESC").
   Limit(20).
   Offset(40).
   Find(&messages)
```

### Count

```go
var count int64
db.Model(&model.Message{}).
   Where("conversation_id = ?", convUUID).
   Count(&count)
```

---

## 3️⃣ Relationships

### Preload (Eager Loading)

```go
// Load all members
var conversation model.Conversation
db.Preload("Members").First(&conversation, "id = ?", id)

// Selective preload
db.Preload("Members", "role = ?", "admin").First(&conversation)

// Nested preload
db.Preload("Members.User").First(&conversation)
```

### Joins

```go
var conversations []model.Conversation
db.Joins("JOIN conversation_members ON conversation_members.conversation_id = conversations.id").
   Where("conversation_members.user_id = ?", userUUID).
   Find(&conversations)
```

---

## 4️⃣ Transactions

### Basic Transaction

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Create conversation
    if err := tx.Create(&conversation).Error; err != nil {
        return err
    }

    // Add members
    for _, memberID := range memberIDs {
        member := &model.ConversationMember{
            ConversationID: conversation.ID,
            UserID:         memberID,
        }
        if err := tx.Create(member).Error; err != nil {
            return err
        }
    }

    return nil
})
```

### Manual Transaction

```go
tx := db.Begin()

if err := tx.Create(&conversation).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&member).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

---

## 5️⃣ Raw SQL

### Raw Query

```go
type Result struct {
    ID    uuid.UUID
    Count int64
}

var results []Result
db.Raw("SELECT conversation_id as id, COUNT(*) as count FROM messages GROUP BY conversation_id").
   Scan(&results)
```

### Exec

```go
db.Exec("UPDATE messages SET content = ? WHERE id = ?", "Updated", messageUUID)
```

---

## 6️⃣ Context & Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

db.WithContext(ctx).Find(&messages)
```

---

## 7️⃣ Scopes (Reusable Queries)

```go
// Define scope
func ActiveConversations(db *gorm.DB) *gorm.DB {
    return db.Where("deleted_at IS NULL")
}

func RecentMessages(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(50)
}

// Use scope
db.Scopes(ActiveConversations, RecentMessages).Find(&conversations)
```

---

## 8️⃣ Hooks

```go
type Message struct {
    // ... fields
}

// Before create
func (m *Message) BeforeCreate(tx *gorm.DB) error {
    if m.Type == "" {
        m.Type = "text"
    }
    return nil
}

// After create
func (m *Message) AfterCreate(tx *gorm.DB) error {
    // Notify WebSocket hub
    return nil
}
```

---

## 9️⃣ Common Patterns

### Find or Create

```go
var conversation model.Conversation
db.Where(model.Conversation{Type: "direct", CreatedBy: userUUID}).
   Attrs(model.Conversation{Name: "Chat"}).
   FirstOrCreate(&conversation)
```

### Update or Create

```go
db.Where(model.ConversationMember{ConversationID: convUUID, UserID: userUUID}).
   Assign(model.ConversationMember{Role: "admin"}).
   FirstOrCreate(&member)
```

### Batch Insert

```go
messages := []model.Message{
    {Content: "Msg 1", Type: "text"},
    {Content: "Msg 2", Type: "text"},
}
db.Create(&messages)
```

---

## 🔟 Error Handling

```go
// Check if record not found
if err := db.First(&message, id).Error; err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return fmt.Errorf("message not found")
    }
    return err
}

// Check affected rows
result := db.Delete(&message, id)
if result.RowsAffected == 0 {
    return fmt.Errorf("no rows deleted")
}
```

---

## 🎯 Best Practices

### 1. Always use Context

```go
db.WithContext(ctx).Find(&messages)
```

### 2. Use Prepared Statements

```go
// Already enabled in config
PrepareStmt: true
```

### 3. Index Common Queries

```go
type Message struct {
    ConversationID uuid.UUID `gorm:"index"`
    CreatedAt      time.Time `gorm:"index"`
}
```

### 4. Pagination

```go
func GetMessages(page, pageSize int) ([]model.Message, error) {
    var messages []model.Message
    offset := (page - 1) * pageSize

    err := db.Order("created_at DESC").
        Limit(pageSize).
        Offset(offset).
        Find(&messages).Error

    return messages, err
}
```

### 5. Select Specific Fields

```go
// Don't load all fields
db.Select("id", "content", "created_at").Find(&messages)
```

---

## 📊 Performance Tips

1. **Use Indexes** - Already defined in migrations
2. **Preload wisely** - Only load what you need
3. **Pagination** - Always use LIMIT/OFFSET
4. **Connection Pooling** - Already configured
5. **Prepared Statements** - Enabled by default
6. **Select specific fields** - Don't use `SELECT *`

---

## 🔧 Debugging

### Enable SQL Logging

```go
db.Debug().Find(&messages)
```

### Print SQL

```go
db.ToSQL(func(tx *gorm.DB) *gorm.DB {
    return tx.Find(&messages)
})
```

---

## 📚 References

- [GORM Docs](https://gorm.io/docs/)
- [UUID Package](https://github.com/google/uuid)
- Project: `/home/nguyenanhtung/Documents/ChatMeNow/chat-service`

---

**Quick Commands:**

```bash
# Build
go build -o bin/chat-service cmd/server/main.go

# Run
go run cmd/server/main.go

# Test
go test ./...

# Tidy dependencies
go mod tidy
```
