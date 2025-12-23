# Quick Start Guide

Get the Todo API running in 2 minutes!

## Choose Your Setup

### Option 1: Docker (Recommended)

**Requirements:** Docker Desktop + Go 1.25.5+

```bash
# 1. Clone the repo
git clone <repository-url>
cd todo-api

# 2. Download dependencies
go mod tidy

# 3. Run!
./start.sh
```

Done! API running at `http://localhost:8080`

To stop: `./stop.sh`

---

### Option 2: Manual Setup

**Requirements:** PostgreSQL 16 + Go 1.25.5+

```bash
# 1. Clone the repo
git clone <repository-url>
cd taskjoy-be

# 2. Create database
psql -U postgres -c "CREATE DATABASE todo_db;"

# 3. Setup tables
psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql

# 4. Configure environment
cp .env.example .env
# Edit .env:
# - DATABASE_URL: Change port from 5433 to 5432 (local PostgreSQL uses 5432)
# - JWT_SECRET: Set a secure string (min 32 characters)

# 5. Run!
go mod download
go run cmd/api/main.go
```

Done! API running at `http://localhost:8080`

---

## Test It

```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Copy the "token" from the response data field, then:
# Note: All responses come in envelope format: {"success": true, "data": {...}}

# Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"My first todo","description":"Testing the API"}'

# Update a todo (use the ID from the previous response)
curl -X PATCH http://localhost:8080/api/v1/todos/YOUR_TODO_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"completed":true}'
```

---

## What's Next?

- Read the full [README.md](README.md) for detailed documentation
- Check [API_DESIGN.md](API_DESIGN.md) for complete API reference
- Explore the code structure in `internal/`

## Common Issues

**Docker not starting?**
- Make sure Docker Desktop is running
- Try: `./stop.sh` then `./start.sh`

**Database connection failed?**
- Docker: Check `docker-compose ps`
- No Docker: Check PostgreSQL is running

**Port 8080 in use?**
- Change port: `PORT=8081 go run cmd/api/main.go`

---

Need help? Check the [Troubleshooting section](README.md#troubleshooting) in README.md
