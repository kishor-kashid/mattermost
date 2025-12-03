#!/bin/bash
# Complete Mattermost startup script for Windows (Git Bash)
# This replaces 'make run' on Windows

echo "========================================"
echo " Mattermost Windows Development Setup"
echo "========================================"
echo ""

# Navigate to server directory
cd "$(dirname "$0")"

# Function to check if a Docker container is running
container_running() {
    docker ps --format '{{.Names}}' | grep -q "$1"
}

# Step 1: Start Docker services
echo "Step 1/3: Starting Docker services..."
if container_running "mattermost-postgres"; then
    echo "✅ Docker services already running"
else
    echo "Starting PostgreSQL, MinIO, Redis, and other services..."
    docker compose -f docker-compose.makefile.yml up -d \
      postgres \
      minio \
      inbucket \
      redis \
      prometheus \
      grafana \
      loki \
      promtail

    echo "Waiting for PostgreSQL to be ready..."
    sleep 5
    until docker exec mattermost-postgres pg_isready -h localhost > /dev/null 2>&1; do
      echo -n "."
      sleep 2
    done
    echo ""
    echo "✅ Docker services started"
fi

echo ""

# Step 2: Set up Go workspace
echo "Step 2/3: Setting up Go workspace..."
if [ ! -f "go.work" ]; then
    go work init
    go work use .
    go work use ./public
    echo "✅ Go workspace created"
else
    echo "✅ Go workspace already exists"
fi

echo ""

# Step 3: Build mmctl if needed
echo "Step 3/3: Checking mmctl..."
if [ ! -f "bin/mmctl" ]; then
    echo "Building mmctl (first time only, may take a few minutes)..."
    go build -trimpath -ldflags '-X "github.com/mattermost/mattermost/server/v8/cmd/mmctl/commands.buildDate='"$(date -u +'%Y-%m-%dT%H:%M:%SZ')"'"' -o bin/mmctl ./cmd/mmctl
    echo "✅ mmctl built"
else
    echo "✅ mmctl already built"
fi

echo ""

# Create necessary directories
mkdir -p ../webapp/channels/dist/files

echo "========================================"
echo " Starting Mattermost Server"
echo "========================================"
echo ""
echo "Server will be available at:"
echo "  🌐 http://localhost:8065"
echo ""
echo "To stop the server: Press Ctrl+C"
echo ""
echo "Logs will appear below:"
echo "========================================" 
echo ""

# Run the server
exec go run \
  -ldflags '-X "github.com/mattermost/mattermost/server/public/model.BuildNumber=dev" -X "github.com/mattermost/mattermost/server/public/model.BuildDate='"$(date -u)"'" -X "github.com/mattermost/mattermost/server/public/model.BuildHash=dev"' \
  ./cmd/mattermost

