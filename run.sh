#!/bin/bash

# Simple script to run the app with environment variables loaded

# Check if .env exists
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please copy .env.example to .env and configure it."
    exit 1
fi

# Load environment variables and run the app
env $(cat .env | grep -v '^#' | grep -v '^$' | xargs) go run cmd/api/main.go
