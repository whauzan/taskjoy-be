#!/bin/bash

echo "🚀 Starting Todo API..."
echo ""

# Download dependencies
echo "📥 Downloading Go dependencies..."
go mod download
go mod tidy

# Start PostgreSQL
echo "📦 Starting PostgreSQL..."
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be ready..."
sleep 5

# Create the database if it doesn't exist
echo "🗄️  Creating database..."
docker exec $(docker-compose ps -q postgres) psql -U postgres -tc "SELECT 1 FROM pg_database WHERE datname = 'todo_db'" | grep -q 1 || docker exec $(docker-compose ps -q postgres) psql -U postgres -c "CREATE DATABASE todo_db"

# Run migrations manually via docker exec
echo "🔄 Running database migrations..."
docker exec -i $(docker-compose ps -q postgres) psql -U postgres -d todo_db << 'EOF'
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create todos table
CREATE TABLE IF NOT EXISTS todos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index on user_id for faster lookups
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);

-- Create index on user_id and completed for filtering
CREATE INDEX IF NOT EXISTS idx_todos_user_id_completed ON todos(user_id, completed);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS \$\$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
\$\$ language 'plpgsql';

-- Drop existing triggers if they exist
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;

-- Trigger to automatically update updated_at on users table
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger to automatically update updated_at on todos table
CREATE TRIGGER update_todos_updated_at BEFORE UPDATE ON todos
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
EOF

echo "✅ Database setup complete!"
echo ""
echo "🏃 Starting the API server..."
echo ""

# Run the app with environment variables loaded
./run.sh
