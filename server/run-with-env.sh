#!/bin/bash
# Load environment variables from .env file and run the server

# Load .env file if it exists
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(cat .env | grep -v '^#' | xargs)
    echo "✓ Environment variables loaded"
else
    echo "Warning: .env file not found"
fi

# Run the server
MM_NO_DOCKER=true make run-server

