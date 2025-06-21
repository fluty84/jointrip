#!/bin/bash

# Development script for JoinTrip application
# This script runs the application in development mode

set -e

echo "🚀 Starting JoinTrip in development mode..."

# Check if React build exists, if not build it
if [ ! -d "web/dist" ]; then
    echo "📦 React build not found, building frontend..."
    cd web
    source ~/.nvm/nvm.sh
    npm run build
    cd ..
fi

echo "🔨 Starting Go server..."
echo "📱 Frontend will be served at: http://localhost:8080"
echo "🔌 API will be available at: http://localhost:8080/api/v1"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Run the Go application
go run main.go embed.go
