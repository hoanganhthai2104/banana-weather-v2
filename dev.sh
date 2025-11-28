#!/bin/bash
#
# Banana Weather Local Development Script
#
# This script builds the Flutter frontend and starts the Go backend server locally.
# It loads configuration from .env and manages the build process.
#
# Usage:
#   ./dev.sh [options]
#
# Options:
#   --quick    Skip Flutter build (Use this if you haven't changed frontend code)
#

# Load environment variables from .env
if [ -f .env ]; then
  echo "Loading configuration from .env..."
  set -o allexport
  source .env
  set +o allexport
else
  echo "Warning: .env file not found. Ensure environment variables are set manually."
fi

# Check for required keys
if [ -z "$GOOGLE_MAPS_API_KEY" ]; then
  echo "Error: GOOGLE_MAPS_API_KEY in environment."
  echo "Please create a .env file or export these variables."
  exit 1
fi

if [ -z "$PROJECT_ID" ] && [ -z "$GOOGLE_CLOUD_PROJECT" ]; then
    echo "Error: PROJECT_ID or GOOGLE_CLOUD_PROJECT not set."
    exit 1
fi

# 1. Frontend Build (Flutter)
if [[ "$*" == *"--quick"* ]]; then
  echo "------------------------------------------------"
  echo "⚡ Skipping Frontend build (--quick selected)"
else
  echo "------------------------------------------------"
  echo "🎨 Building Frontend (Flutter Web)..."
  (cd frontend && flutter build web)
  if [ $? -ne 0 ]; then
    echo "Error: Flutter build failed."
    exit 1
  fi
fi

# 2. Backend Build & Run (Go)
echo "------------------------------------------------"
echo "🍌 Building and Starting Backend Server..."
cd backend
go build -o server
if [ $? -eq 0 ]; then
  port="${PORT:-8080}"
  echo "✅ Server starting on port $port..."
  echo "   - Web App: http://localhost:$port"
  echo "   - API:     http://localhost:$port/api/weather"
  ./server
else
  echo "❌ Backend build failed."
  exit 1
fi
