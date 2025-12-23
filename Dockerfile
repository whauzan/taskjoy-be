# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o todo-api ./cmd/api

# Final stage - use distroless for smaller, more secure image
FROM alpine:3.19

# Install ca-certificates for HTTPS and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/todo-api .

# Copy migrations (optional, for manual migration runs)
COPY --from=builder /app/db/migrations ./db/migrations

# Use non-root user
USER appuser

# Expose port (Render uses PORT env var)
EXPOSE 8080

# Run the application
CMD ["./todo-api"]
