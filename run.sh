#!/bin/bash

# Set environment variables for the backend
export SERVER_PORT=8081
export FRONTEND_URL=http://localhost:3000

# Start from the script directory
cd "$(dirname "$0")"

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Start the Go backend
echo "Starting Go backend on port $SERVER_PORT..."
go run cmd/main.go &

# Set environment variables for the frontend
export REACT_APP_API_BASE_URL=http://localhost:8081

# Move to frontend directory
cd ../web-analyzer-frontend

# Install Node.js dependencies
echo "Installing frontend dependencies..."
npm install

# Start the React frontend
echo "Starting React frontend on port 3000..."
npm start &

# Wait for both processes
wait