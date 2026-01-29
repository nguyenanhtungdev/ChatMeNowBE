# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in ChatMeNow, please follow these steps:

### 1. Do NOT Open a Public Issue

Security vulnerabilities should not be publicly disclosed until they have been addressed.

### 2. Report Privately

Send an email to: **security@chatmenow.example.com** (or create a private Security Advisory on GitHub)

Include:

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity

### 4. Disclosure Policy

- We will acknowledge your report
- We will investigate and validate the issue
- We will develop and test a fix
- We will release a patch
- We will publicly disclose (with credit to reporter if desired)

## Security Best Practices for Deployment

### 1. Environment Variables

**Never commit sensitive data!**

```bash
# Change these in production!
JWT_SECRET=<use-strong-64-char-random-string>
POSTGRES_PASSWORD=<strong-unique-password>
MONGO_INITDB_ROOT_PASSWORD=<strong-unique-password>
```

Generate strong secrets:

```bash
openssl rand -base64 64
```

### 2. HTTPS Only

Always use HTTPS in production:

- Use Let's Encrypt for free SSL certificates
- Configure reverse proxy (nginx, Caddy)
- Enable HSTS headers

### 3. CORS Configuration

Update CORS in production to whitelist specific origins:

```typescript
// gateway/src/main.ts
app.enableCors({
  origin: ["https://yourapp.com", "https://www.yourapp.com"],
  credentials: true,
});
```

### 4. Rate Limiting

Adjust based on your needs:

```env
RATE_LIMIT_TTL=60      # seconds
RATE_LIMIT_MAX=100     # requests per TTL
```

### 5. Database Security

**PostgreSQL:**

```sql
-- Create separate users with minimal permissions
CREATE USER app_user WITH PASSWORD 'strong_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO app_user;
```

**MongoDB:**

```javascript
// Enable authentication
// Use role-based access control
db.createUser({
  user: "app_user",
  pwd: "strong_password",
  roles: ["readWrite"],
});
```

### 6. JWT Token Security

- **Access Token**: Short-lived (15 minutes)
- **Refresh Token**: Rotate on use, blacklist on logout
- **Secret**: Use strong random secret (64+ characters)
- **Algorithm**: HS256 or RS256 (asymmetric recommended for multi-service)

### 7. Input Validation

All DTOs use class-validator:

```typescript
export class RegisterDto {
  @IsString()
  @MinLength(3)
  username: string;

  @IsEmail()
  email: string;

  @IsString()
  @MinLength(6)
  password: string;
}
```

### 8. SQL Injection Prevention

Always use parameterized queries:

```typescript
// ✅ Good
const user = await this.userRepo.findOne({ where: { email } });

// ❌ Bad
const query = `SELECT * FROM users WHERE email = '${email}'`;
```

### 9. XSS Protection

Sanitize user input before displaying:

```typescript
import { escape } from "html-escaper";
const safe = escape(userInput);
```

### 10. Dependencies

Regularly update dependencies:

```bash
npm audit
npm audit fix

go list -m all | go-mod-outdated
go get -u ./...
```

## Known Security Considerations

### Current Implementation

1. **JWT Secret**: Shared across services (convenient for development)
   - **Production**: Use asymmetric keys (RS256) or separate secrets

2. **CORS**: Currently allows all origins
   - **Production**: Whitelist specific domains

3. **Rate Limiting**: Basic implementation
   - **Production**: Consider distributed rate limiting (Redis)

4. **WebSocket**: No rate limiting
   - **Future**: Add per-connection message rate limit

5. **File Upload**: Not implemented yet
   - **Future**: Add virus scanning, size limits, type validation

## Compliance

ChatMeNow follows security best practices for:

- OWASP Top 10
- CWE/SANS Top 25
- GDPR considerations (data encryption, user deletion)

## Security Checklist for Production

- [ ] Change all default passwords
- [ ] Use strong random JWT secret
- [ ] Enable HTTPS only
- [ ] Configure CORS whitelist
- [ ] Set up firewall rules
- [ ] Enable database authentication
- [ ] Use environment variables for secrets
- [ ] Set up log monitoring
- [ ] Enable rate limiting
- [ ] Regular security updates
- [ ] Backup encryption
- [ ] Network segmentation
- [ ] Security headers (helmet.js)

## Contact

For security concerns: **security@chatmenow.example.com**

Thank you for helping keep ChatMeNow secure! 🔒
