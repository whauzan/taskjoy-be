# API Design Documentation

This document provides detailed information about all API endpoints, including request/response formats, status codes, and error handling.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require authentication using JWT tokens. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Response Format

All API responses use a consistent envelope format.

### Success Response Format

```json
{
  "success": true,
  "data": { /* response data here */ }
}
```

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": ["Optional array of detailed error messages"]
  }
}
```

### Error Codes

- `VALIDATION_ERROR` - Input validation failed
- `INVALID_CREDENTIALS` - Invalid email or password
- `USER_EXISTS` - User with email already exists
- `NOT_FOUND` - Resource not found
- `FORBIDDEN` - Access denied
- `UNAUTHORIZED` - Authentication required
- `INTERNAL_ERROR` - Internal server error
- `BAD_REQUEST` - Bad request

## Endpoints

### Health Check

#### GET /health

Check the health status of the API and database.

**Authentication:** Not required

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "database": "healthy",
    "time": "2025-12-23T10:00:00Z"
  }
}
```

**Response:** 503 Service Unavailable (if database is down)

```json
{
  "success": true,
  "data": {
    "status": "unhealthy",
    "database": "unhealthy",
    "time": "2025-12-23T10:00:00Z"
  }
}
```

---

## Authentication Endpoints

### Register User

#### POST /api/v1/auth/register

Register a new user account.

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Validation Rules:**

- `email`: Required, valid email format, max 255 characters
- `password`: Required, min 8 characters, max 72 characters
- `name`: Required, min 1 character, max 255 characters

**Response:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-12-23T10:00:00Z"
  }
}
```

**Error Response:** 400 Bad Request

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      "email: is required",
      "password: must be at least 8 characters"
    ]
  }
}
```

**Error Response:** 409 Conflict

```json
{
  "success": false,
  "error": {
    "code": "USER_EXISTS",
    "message": "User with this email already exists"
  }
}
```

---

### Login

#### POST /api/v1/auth/login

Authenticate a user and receive a JWT token.

**Authentication:** Not required

**Request Body:**

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Validation Rules:**

- `email`: Required, valid email format
- `password`: Required

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-12-26T10:00:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

**Error Response:** 401 Unauthorized

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password"
  }
}
```

---

### Refresh Token

#### POST /api/v1/auth/refresh

Refresh an existing JWT token to extend the session. This endpoint is useful for mobile apps to keep users logged in without requiring re-authentication. The token must still be valid (not expired) to be refreshed.

**Authentication:** Required (Bearer token in Authorization header)

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:** None

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2025-12-27T10:00:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "John Doe"
    }
  }
}
```

**Error Response:** 401 Unauthorized

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  }
}
```

**Use Case (Mobile Apps):**

In your Flutter app, check if the token expires soon before making API calls:

```dart
// Before each API call
if (tokenExpiresIn < 5.minutes && userIsActive) {
  final newToken = await refreshToken(currentToken);
  // Update stored token
}
```

This ensures users remain logged in during active use without seeing session expiration warnings.

---

## Todo Endpoints

All todo endpoints require authentication.

### List Todos

#### GET /api/v1/todos

Get all todos for the authenticated user.

**Authentication:** Required

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response:** 200 OK

```json
{
  "success": true,
  "data": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Buy groceries",
      "description": "Milk, eggs, bread",
      "completed": false,
      "created_at": "2025-12-22T10:00:00Z",
      "updated_at": "2025-12-22T10:00:00Z"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Write documentation",
      "description": null,
      "completed": true,
      "created_at": "2025-12-22T09:00:00Z",
      "updated_at": "2025-12-22T11:00:00Z"
    }
  ]
}
```

**Response:** 200 OK (empty list)

```json
{
  "success": true,
  "data": []
}
```

**Error Response:** 401 Unauthorized

```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Authentication required"
  }
}
```

---

### Create Todo

#### POST /api/v1/todos

Create a new todo item.

**Authentication:** Required

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Request Body:**

```json
{
  "title": "Buy groceries",
  "description": "Milk, eggs, bread"
}
```

**Validation Rules:**

- `title`: Required, min 1 character, max 255 characters
- `description`: Optional, max 2000 characters

**Response:** 201 Created

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false,
    "created_at": "2025-12-22T10:00:00Z",
    "updated_at": "2025-12-22T10:00:00Z"
  }
}
```

**Error Response:** 400 Bad Request

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": [
      "title: is required"
    ]
  }
}
```

---

### Get Single Todo

#### GET /api/v1/todos/{id}

Get a specific todo by ID.

**Authentication:** Required

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**URL Parameters:**

- `id`: UUID of the todo

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": false,
    "created_at": "2025-12-22T10:00:00Z",
    "updated_at": "2025-12-22T10:00:00Z"
  }
}
```

**Error Response:** 400 Bad Request

```json
{
  "success": false,
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid todo ID"
  }
}
```

**Error Response:** 404 Not Found

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Todo not found"
  }
}
```

**Error Response:** 403 Forbidden

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

---

### Update Todo

#### PATCH /api/v1/todos/{id}

Partially update a todo item. All fields are optional - only send the fields you want to change.

**Authentication:** Required

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**URL Parameters:**

- `id`: UUID of the todo

**Request Body:**

```json
{
  "title": "Buy groceries and cook dinner",
  "description": "Milk, eggs, bread, chicken",
  "completed": true
}
```

**Validation Rules:**

- `title`: Optional, min 1 character, max 255 characters
- `description`: Optional, max 2000 characters
- `completed`: Optional, boolean

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries and cook dinner",
    "description": "Milk, eggs, bread, chicken",
    "completed": true,
    "created_at": "2025-12-22T10:00:00Z",
    "updated_at": "2025-12-22T11:30:00Z"
  }
}
```

**Partial Update Example:**

```json
{
  "completed": true
}
```

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Buy groceries",
    "description": "Milk, eggs, bread",
    "completed": true,
    "created_at": "2025-12-22T10:00:00Z",
    "updated_at": "2025-12-22T11:35:00Z"
  }
}
```

**Error Response:** 404 Not Found

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Todo not found"
  }
}
```

**Error Response:** 403 Forbidden

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

---

### Delete Todo

#### DELETE /api/v1/todos/{id}

Delete a todo item.

**Authentication:** Required

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**URL Parameters:**

- `id`: UUID of the todo

**Response:** 200 OK

```json
{
  "success": true,
  "data": {
    "message": "Todo deleted successfully"
  }
}
```

**Error Response:** 404 Not Found

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Todo not found"
  }
}
```

**Error Response:** 403 Forbidden

```json
{
  "success": false,
  "error": {
    "code": "FORBIDDEN",
    "message": "You don't have permission to access this resource"
  }
}
```

---

## HTTP Status Codes

The API uses the following HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Request successful, no content to return
- `400 Bad Request` - Invalid request format or validation error
- `401 Unauthorized` - Authentication required or invalid token
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - Service temporarily unavailable

## Rate Limiting

Currently, there is no rate limiting implemented. Consider adding rate limiting for production use.

## CORS

CORS is configured to allow requests from origins specified in the `CORS_ALLOWED_ORIGINS` environment variable.

Default allowed origins in development:
- `http://localhost:3000`
- `http://localhost:8080`

## Request ID

Every request receives a unique Request ID in the `X-Request-ID` response header. This can be used for debugging and tracing requests through logs.

## Timestamps

All timestamps are in UTC and follow the RFC3339 format:

```
2025-12-22T10:00:00Z
```

## Testing with cURL

### Complete Flow Example

```bash
# 1. Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# 2. Login and save token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.token')

# 3. Refresh token (optional - extends session)
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.token')

# 4. Create a todo
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"title":"Test Todo","description":"Test Description"}'

# 5. List todos
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer $TOKEN"

# 6. Update todo (replace TODO_ID with actual ID)
curl -X PATCH http://localhost:8080/api/v1/todos/TODO_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"completed":true}'

# 7. Delete todo
curl -X DELETE http://localhost:8080/api/v1/todos/TODO_ID \
  -H "Authorization: Bearer $TOKEN"
```
