#!/bin/bash
# Windows-compatible server run script
# Run this in Git Bash

echo "Starting Mattermost server..."

# Navigate to server directory if not already there
cd "$(dirname "$0")"

# Check if Docker services are running
if ! docker ps | grep -q mattermost-postgres; then
    echo "❌ Docker services not running!"
    echo "Please run: ./start-docker-windows.sh"
    exit 1
fi

echo "✅ Docker services detected"
echo ""

# Set up Go workspace if needed
if [ ! -f "go.work" ]; then
    echo "Setting up Go workspace..."
    go work init
    go work use .
    go work use ./public
fi

# Build prepackaged binaries if needed
if [ ! -f "bin/mmctl" ]; then
    echo "Building mmctl..."
    make mmctl-build
fi

# Create necessary directories
mkdir -p ../webapp/channels/dist/files

echo ""
echo "Starting Mattermost server..."
echo "Access at: http://localhost:8065"
echo ""

# Run the server
go run \
  -ldflags '-X "github.com/mattermost/mattermost/server/public/model.BuildNumber=dev" -X "github.com/mattermost/mattermost/server/public/model.BuildDate=$(date -u)" -X "github.com/mattermost/mattermost/server/public/model.BuildHash=dev"' \
  ./cmd/mattermost

