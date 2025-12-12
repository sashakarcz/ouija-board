#!/bin/bash

set -e

echo "==================================="
echo "Ouija Board Go Setup Script"
echo "==================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "Found Go version: $GO_VERSION"

# Download dependencies
echo ""
echo "Downloading dependencies..."
go mod download
go mod tidy
echo "Dependencies downloaded successfully"

# Build the application
echo ""
echo "Building application..."
go build -o ouija-board .
echo "Build successful: ouija-board"

# Create run script
cat > run.sh << 'EOF'
#!/bin/bash

# Default configuration
export SERVER_ADDR="${SERVER_ADDR:-0.0.0.0:8080}"
export OLLAMA_URL="${OLLAMA_URL:-http://localhost:11434/api/generate}"
export OLLAMA_MODEL="${OLLAMA_MODEL:-qwen3}"
export OLLAMA_TIMEOUT="${OLLAMA_TIMEOUT:-30s}"
export MAX_TOKENS="${MAX_TOKENS:-10}"
export MAX_HISTORY_SIZE="${MAX_HISTORY_SIZE:-1000}"
export RATE_LIMIT="${RATE_LIMIT:-10}"

echo "Starting Ouija Board application..."
echo "Server will be available at: http://${SERVER_ADDR}"
echo "Using Ollama at: ${OLLAMA_URL}"
echo ""

./ouija-board
EOF

chmod +x run.sh

echo ""
echo "==================================="
echo "Setup Complete!"
echo "==================================="
echo ""
echo "To run the application:"
echo "  ./run.sh"
echo ""
echo "Or with custom configuration:"
echo "  OLLAMA_URL=http://your-ollama:11434/api/generate ./run.sh"
echo ""
echo "To build for Docker:"
echo "  docker build -f Dockerfile-go -t ouija-board:go ."
echo ""
echo "To run with docker-compose:"
echo "  docker-compose -f docker-compose-go.yaml up -d"
echo ""
echo "See README.md for more information"
echo "==================================="
