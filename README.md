# Todo API

A production-grade RESTful API for a Todo application built with Go, following Clean Architecture principles.

**🚀 New to this project?** Check out the [Quick Start Guide](QUICKSTART.md) to get running in 2 minutes!

## Features

- User authentication with JWT
- CRUD operations for todos
- Clean Architecture with clear separation of concerns
- PostgreSQL database with connection pooling
- Database migrations with golang-migrate
- Type-safe SQL queries with sqlc
- Structured logging with slog
- Graceful shutdown
- Docker support
- Comprehensive error handling
- Input validation
- CORS support

## Tech Stack

- **Go**: 1.25.5
- **Router**: Chi v5
- **Database**: PostgreSQL with pgx v5
- **Query Generation**: sqlc
- **Migrations**: golang-migrate
- **JWT**: golang-jwt v5
- **Password Hashing**: bcrypt
- **Config**: env v11
- **Validation**: validator v10
- **UUID**: google/uuid
- **Logging**: log/slog (stdlib)
- **CORS**: go-chi/cors
- **Testing**: testify

## Project Structure

```
todo-api/
├── cmd/api/              # Application entrypoint
├── internal/
│   ├── config/          # Configuration loading
│   ├── domain/          # Domain entities
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # HTTP middleware
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic
│   └── pkg/             # Shared utilities
├── db/
│   ├── migrations/      # Database migrations
│   └── queries/         # SQL queries for sqlc
├── start.sh             # Quick start script (Docker)
├── stop.sh              # Stop script (Docker)
├── docker-compose.yml   # Docker configuration
├── Dockerfile           # Docker image definition
├── fly.toml             # Fly.io deployment config
└── README.md            # This file
```

## Prerequisites

**Minimum Requirements:**
- Go 1.25.5 or higher
- PostgreSQL 16 (or use Docker)

**Choose Your Path:**
- **With Docker** (Easiest): Just install Docker Desktop
- **Without Docker**: You'll need PostgreSQL installed locally

## Quick Start

### Option 1: Docker (Recommended)

**Prerequisites:**
- [Docker Desktop](https://www.docker.com/products/docker-desktop) installed and running
- Go 1.25.5 or higher

**Steps:**

```bash
# 1. Clone the repository
git clone <repository-url>
cd todo-api

# 2. Download dependencies and start
go mod tidy  # First time only
./start.sh
```

That's it! The API will be running at `http://localhost:8080`

**Note:** The script will automatically download dependencies on first run.

**To stop:**
```bash
./stop.sh
```

The `start.sh` script automatically:
- Starts PostgreSQL in Docker
- Creates the database and tables
- Starts the API server

---

### Option 2: Manual Setup

**Prerequisites:**
- Go 1.25.5 or higher
- PostgreSQL 16 installed and running locally

**Steps:**

```bash
# 1. Clone the repository
git clone <repository-url>
cd todo-api

# 2. Create the database
psql -U postgres -c "CREATE DATABASE todo_db;"

# 3. Run the migration SQL
psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql

# 4. Copy and configure environment variables
cp .env.example .env

# Edit .env file and update:
# - DATABASE_URL: Change port from 5433 to 5432 (local PostgreSQL uses 5432)
#   DATABASE_URL=postgres://YOUR_USERNAME:YOUR_PASSWORD@localhost:5432/todo_db?sslmode=disable
# - JWT_SECRET: Set a secure random string (min 32 characters)

# 5. Install Go dependencies
go mod download

# 6. Run the application
go run cmd/api/main.go
```

The API will be running at `http://localhost:8080`

---

## Verify It's Working

After starting the application, test the health check:

```bash
curl http://localhost:8080/health
```

You should see:
```json
{
  "status": "healthy",
  "database": "healthy",
  "time": "2025-12-23T..."
}
```

## API Endpoints

### Health Check

```
GET /health
```

### Authentication

```
POST /api/v1/auth/register  - Register new user
POST /api/v1/auth/login     - Login and get JWT token
POST /api/v1/auth/refresh   - Refresh JWT token
```

### Todos (Authenticated)

```
GET    /api/v1/todos        - Get all todos
POST   /api/v1/todos        - Create a new todo
GET    /api/v1/todos/{id}   - Get a specific todo
PATCH  /api/v1/todos/{id}   - Update a todo (partial update)
DELETE /api/v1/todos/{id}   - Delete a todo
```

## Usage Examples

### Register a User

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "name": "John Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGc...",
    "expires_at": "2025-12-26T10:00:00Z",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
} 
```

### Refresh Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Response:
```json
{
  "success": true,
  "data": {
    "token": "new-jwt-token...",
    "expires_at": "2025-12-27T10:00:00Z",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

### Create a Todo

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Buy groceries",
    "description": "Milk, eggs, bread"
  }'
```

### Get All Todos

```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Update a Todo

```bash
curl -X PATCH http://localhost:8080/api/v1/todos/{id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Buy groceries and cook",
    "completed": true
  }'
```

### Delete a Todo

```bash
curl -X DELETE http://localhost:8080/api/v1/todos/{id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Deployment

### Deploy to Fly.io

This project includes a `fly.toml` for deployment to Fly.io.

**Prerequisites:**
- Install [flyctl](https://fly.io/docs/hands-on/install-flyctl/)
- Run `fly auth login`

**Steps:**

```bash
# 1. Create the app (first time only)
fly launch --no-deploy

# 2. Create PostgreSQL database
fly postgres create --name taskjoy-db

# 3. Attach database to app (sets DATABASE_URL automatically)
fly postgres attach taskjoy-db

# 4. Set secrets
fly secrets set JWT_SECRET="your-secure-secret-min-32-characters"
fly secrets set CORS_ALLOWED_ORIGINS="https://your-frontend.com"

# 5. Deploy
fly deploy
```

**Run Migrations on Fly.io:**
```bash
# Connect to your Fly Postgres
fly postgres connect -a taskjoy-db

# Then run in psql:
\i db/migrations/000001_init.up.sql
```

**Useful Commands:**
```bash
fly status              # Check app status
fly logs                # View logs
fly ssh console         # SSH into the machine
fly secrets list        # List secrets
```

## Database Migrations

### With Docker
If you used `start.sh`, migrations are already applied automatically!

To manually manage migrations:
```bash
# Run migrations via Docker
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql

# Rollback migrations
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d todo_db -f db/migrations/000001_init.down.sql
```

### Without Docker
```bash
# Run migrations
psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql

# Rollback migrations
psql -U postgres -d todo_db -f db/migrations/000001_init.down.sql
```

## Development

### Common Commands

```bash
go fmt ./...                                      # Format code
go vet ./...                                      # Run linter
go test -v ./...                                  # Run tests
go build -o bin/todo-api cmd/api/main.go         # Build binary
rm -rf bin/ coverage.out coverage.html           # Clean
```

## Environment Variables

See `.env.example` for all available environment variables:

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development, staging, production)
- `DATABASE_URL` - PostgreSQL connection string
- `JWT_SECRET` - Secret key for JWT (min 32 characters)
- `JWT_EXPIRY_HOURS` - JWT token expiry in hours (default: 72)
- `CORS_ALLOWED_ORIGINS` - Comma-separated allowed origins
- `LOG_LEVEL` - Log level (debug, info, warn, error)

## Troubleshooting

### With Docker

**PostgreSQL not starting:**
```bash
# Check if Docker is running
docker ps

# Check PostgreSQL container status
docker-compose ps

# View PostgreSQL logs
docker-compose logs postgres

# Restart everything
./stop.sh
./start.sh
```

**Database connection failed:**
```bash
# Make sure containers are running
docker-compose ps

# Restart PostgreSQL
docker-compose restart postgres

# Check if port 5432 is already in use
lsof -i :5432
```

**Start fresh (remove all data):**
```bash
./stop.sh
docker-compose down -v  # This removes volumes/data
./start.sh
```

### Without Docker

**Database connection failed:**
```bash
# Check if PostgreSQL is running
# macOS:
brew services list | grep postgresql

# Linux:
sudo systemctl status postgresql

# Start PostgreSQL if not running
# macOS:
brew services start postgresql

# Linux:
sudo systemctl start postgresql
```

**Tables not found:**
```bash
# Run migrations again
psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql
```

**Reset database:**
```bash
# Drop and recreate database
psql -U postgres -c "DROP DATABASE IF EXISTS todo_db;"
psql -U postgres -c "CREATE DATABASE todo_db;"
psql -U postgres -d todo_db -f db/migrations/000001_init.up.sql
```

### General Issues

**Port 8080 already in use:**
```bash
# Find what's using port 8080
lsof -i :8080

# Kill the process or change PORT in .env
PORT=8081 go run cmd/api/main.go
```

**JWT errors:**
- Make sure `JWT_SECRET` in `.env` is at least 32 characters long
- Generate a secure secret: `openssl rand -base64 32`

**Port 5432/5433 conflicts (PostgreSQL):**
- If you have local PostgreSQL installed, it uses port 5432
- Docker PostgreSQL is configured to use port 5433 to avoid conflicts
- **Option 1:** Stop local PostgreSQL: `brew services stop postgresql` (macOS) or `sudo systemctl stop postgresql` (Linux)
- **Option 2:** Use the Docker PostgreSQL on port 5433 (already configured in `.env`)
- If you see "database does not exist" but tables are in Docker, you're connecting to the wrong PostgreSQL instance

## License

MIT

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
