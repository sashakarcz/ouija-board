# Quick Start Guide

## The Fastest Way to Get Started

### Option 1: Using the Setup Script (Recommended)

```bash
# Make setup script executable (if not already)
chmod +x setup.sh

# Run setup
./setup.sh

# Start the application
./run.sh
```

The application will be available at `http://localhost:8080`

### Option 2: Manual Build

```bash
# Download dependencies
go mod download

# Build
go build -o ouija-board .

# Run
./ouija-board
```

### Option 3: Using Make

```bash
# Build and run
make run
```

### Option 4: Using Docker

```bash
# Build Docker image
make docker-build

# Run container
make docker-run

# Or use docker-compose
make compose-up
```

## Configuration

Set environment variables before running:

```bash
export OLLAMA_URL="http://your-ollama-host:11434/api/generate"
export OLLAMA_MODEL="qwen3"
./run.sh
```

Or create a `.env` file (not tracked by git):

```bash
SERVER_ADDR=0.0.0.0:8080
OLLAMA_URL=http://your-ollama-host:11434/api/generate
OLLAMA_MODEL=qwen3
OLLAMA_TIMEOUT=30s
MAX_TOKENS=10
MAX_HISTORY_SIZE=1000
RATE_LIMIT=10
```

Then source it:

```bash
set -a
source .env
set +a
./ouija-board
```

## Testing

Test the API:

```bash
# Ask a question
curl -X POST http://localhost:8080/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "Does this work?"}'

# Get history
curl http://localhost:8080/history
```

## Common Commands

```bash
# Build
make build

# Run
make run

# Test
make test

# Format code
make fmt

# Run linter
make lint

# Build Docker image
make docker-build

# Run with docker-compose
make compose-up

# View docker-compose logs
make compose-logs

# Stop docker-compose
make compose-down

# Clean build artifacts
make clean

# See all available commands
make help
```

## Troubleshooting

**Cannot connect to Ollama:**
```bash
# Check if Ollama is running
curl http://your-ollama-host:11434/api/tags

# Update OLLAMA_URL
export OLLAMA_URL="http://correct-host:11434/api/generate"
```

**Port already in use:**
```bash
# Use a different port
export SERVER_ADDR="0.0.0.0:9090"
./ouija-board
```

**Permission denied on binary:**
```bash
chmod +x ouija-board
```

## Next Steps

- Read [README.md](README.md) for comprehensive documentation
- Review [SECURITY.md](SECURITY.md) for security details
- Check [CONVERSION_SUMMARY.md](CONVERSION_SUMMARY.md) for migration information

## Need Help?

- Check logs: Application logs to stdout
- Open an issue on GitHub
- Contact: sasha@starnix.net
