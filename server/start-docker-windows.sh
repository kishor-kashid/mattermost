#!/bin/bash
# Windows-compatible Docker start script
# Run this in Git Bash instead of 'make start-docker'

echo "Starting Docker services for Mattermost development..."

# Navigate to server directory if not already there
cd "$(dirname "$0")"

# Start essential services
docker compose -f docker-compose.makefile.yml up -d \
  postgres \
  minio \
  inbucket \
  redis \
  prometheus \
  grafana \
  loki \
  promtail

echo ""
echo "Waiting for services to be ready..."
sleep 15

# Check if PostgreSQL is ready
echo "Checking PostgreSQL..."
until docker exec mattermost-postgres pg_isready -h localhost > /dev/null 2>&1; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo ""
echo "✅ Docker services started successfully!"
echo ""
echo "Services running:"
echo "  - PostgreSQL:  localhost:5432"
echo "  - MinIO:       localhost:9000"
echo "  - Inbucket:    localhost:9001"
echo "  - Redis:       localhost:6379"
echo "  - Prometheus:  localhost:9090"
echo "  - Grafana:     localhost:3000"
echo ""
echo "You can now run the server with:"
echo "  make run-server"
echo ""

