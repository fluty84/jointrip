#!/bin/bash

# Build script for JoinTrip application
# This script builds the React frontend and Go backend

set -e

echo "🚀 Building JoinTrip application..."

# Build React frontend
echo "📦 Building React frontend..."
cd web
source ~/.nvm/nvm.sh
npm run build
cd ..

# Build Go backend with embedded React files
echo "🔨 Building Go backend..."
go build -o bin/jointrip

echo "✅ Build completed successfully!"
echo "📁 Binary location: bin/jointrip"
echo ""
echo "To run the application:"
echo "  ./bin/jointrip"
echo ""
echo "Or run in development mode:"
echo "  go run main.go embed.go"
