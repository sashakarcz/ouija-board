# Ouija Board Web Application

A mystical web application that uses AI (via Ollama) to answer questions like a Ouija board. The application animates a planchette across a digital Ouija board to spell out answers.

## Overview

This is a Go-based web application that:
- Accepts user questions via a web interface
- Generates mystical responses using the Ollama AI API
- Animates a planchette to spell out answers on a digital Ouija board
- Maintains a history of questions and answers in memory
- Provides robust security, rate limiting, and error handling

## Features

- **AI-Powered Responses**: Connects to Ollama API for generating mystical answers
- **Interactive UI**: Animated planchette movement across the Ouija board
- **Security Hardened**: Multiple security headers, input validation, and XSS protection
- **Rate Limiting**: Per-IP rate limiting to prevent abuse
- **Graceful Shutdown**: Proper signal handling and connection draining
- **Configurable**: All settings configurable via environment variables
- **Production Ready**: Built with best practices for production deployment
- **Observability**: Optional OpenTelemetry integration for distributed tracing

## Architecture

The application is structured as follows:

```
ouija-board/
├── main.go           # Application entry point and server setup
├── config.go         # Configuration management
├── handlers.go       # HTTP request handlers
├── middleware.go     # HTTP middleware (logging, rate limiting, security)
├── storage.go        # In-memory storage for Q&A history
├── ollama.go         # Ollama API client
├── static/           # Static assets (CSS, JavaScript, images)
├── templates/        # HTML templates
└── Dockerfile.go     # Docker build configuration
```

## Security Features

### Implemented Protections

1. **Input Validation**
   - Question length limited to 1000 characters
   - Empty questions rejected
   - Input sanitization removes control characters

2. **Rate Limiting**
   - Per-IP rate limiting (default: 10 requests/second)
   - Automatic cleanup of old limiters

3. **Security Headers**
   - X-Frame-Options: DENY (prevents clickjacking)
   - X-Content-Type-Options: nosniff (prevents MIME sniffing)
   - X-XSS-Protection: enabled
   - Content-Security-Policy: restricts resource loading
   - Referrer-Policy: strict-origin-when-cross-origin

4. **Error Handling**
   - No internal information leakage in error messages
   - Graceful degradation on API failures
   - Proper timeout handling

5. **Memory Management**
   - Maximum history size enforced (default: 1000 entries)
   - Automatic removal of oldest entries
   - Thread-safe operations with proper locking

6. **Network Security**
   - Request timeouts (30 seconds default)
   - Server timeouts configured
   - Graceful shutdown on termination signals

## Fixes from Python Version

The Go version addresses the following issues from the original Python implementation:

| Issue | Python Version | Go Version |
|-------|----------------|------------|
| Race Conditions | File writes without locking | Thread-safe in-memory storage with RWMutex |
| Unbounded Growth | Unlimited file and memory growth | Enforced maximum size (1000 entries) |
| Input Validation | None | Length limits and sanitization |
| Error Handling | Minimal error handling | Comprehensive error handling throughout |
| Request Timeouts | No timeout on Ollama API | 30-second timeout (configurable) |
| Rate Limiting | None | Per-IP rate limiting |
| Security Headers | None | Comprehensive security headers |
| Configuration | Hardcoded values | Environment-based configuration |
| Graceful Shutdown | None | Proper signal handling and connection draining |
| JSON Parsing Errors | Unhandled | Handled with graceful fallback |

## Configuration

All configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_ADDR` | `0.0.0.0:8080` | Server address and port |
| `OLLAMA_URL` | `http://localhost:11434/api/generate` | Ollama API endpoint |
| `OLLAMA_MODEL` | `qwen3` | Ollama model to use |
| `OLLAMA_TIMEOUT` | `30s` | Timeout for Ollama API requests |
| `MAX_HISTORY_SIZE` | `1000` | Maximum number of Q&A pairs to keep in memory |
| `MAX_TOKENS` | `10` | Maximum tokens for AI response |
| `RATE_LIMIT` | `10` | Maximum requests per second per IP |
| `ENABLE_OTEL` | `false` | Enable OpenTelemetry tracing |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `http://localhost:4317` | OpenTelemetry collector endpoint |

## Installation

### Prerequisites

- Go 1.21 or higher
- Access to an Ollama instance
- Docker and Docker Compose (for containerized deployment)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/sashakarcz/ouija-board.git
cd ouija-board
```

2. Install dependencies:
```bash
go mod download
```

3. Set environment variables (optional):
```bash
export OLLAMA_URL="http://your-ollama-instance:11434/api/generate"
export OLLAMA_MODEL="qwen3"
```

4. Run the application:
```bash
go run .
```

5. Open your browser to `http://localhost:8080`

### Docker Deployment

1. Build the Docker image:
```bash
docker build -f Dockerfile-go -t ouija-board:go .
```

2. Run the container:
```bash
docker run -p 8080:8080 \
  -e OLLAMA_URL="http://your-ollama-instance:11434/api/generate" \
  ouija-board:go
```

### Docker Compose Deployment

1. Update `docker-compose-go.yaml` with your configuration

2. Start the service:
```bash
docker-compose -f docker-compose-go.yaml up -d
```

3. View logs:
```bash
docker-compose -f docker-compose-go.yaml logs -f
```

4. Stop the service:
```bash
docker-compose -f docker-compose-go.yaml down
```

## API Endpoints

### GET /
Returns the main HTML interface.

### POST /ask
Submit a question to the Ouija board.

**Request:**
```json
{
  "question": "What is the meaning of life?"
}
```

**Response:**
```json
{
  "answer": "The answer lies within you."
}
```

**Error Response:**
```json
{
  "error": "Question too long (max 1000 characters)"
}
```

### GET /history
Retrieve all Q&A history.

**Response:**
```json
[
  {
    "question": "What is the meaning of life?",
    "answer": "The answer lies within you."
  }
]
```

### GET /static/*
Serves static assets (CSS, JavaScript, images).

## Building for Production

### Build Binary

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ouija-board .
```

### Build Docker Image (Multi-stage)

The provided `Dockerfile-go` uses a multi-stage build:

1. **Builder stage**: Compiles the Go application
2. **Runtime stage**: Creates a minimal scratch-based image

Final image size: Approximately 15-20 MB

```bash
docker build -f Dockerfile-go -t ouija-board:go .
```

## Testing

### Manual Testing

1. Start the application
2. Navigate to `http://localhost:8080`
3. Enter a question
4. Observe the planchette animation and answer

### API Testing

Test the `/ask` endpoint:
```bash
curl -X POST http://localhost:8080/ask \
  -H "Content-Type: application/json" \
  -d '{"question": "Does this work?"}'
```

Test the `/history` endpoint:
```bash
curl http://localhost:8080/history
```

Test rate limiting:
```bash
for i in {1..20}; do curl -X POST http://localhost:8080/ask -H "Content-Type: application/json" -d '{"question": "test"}'; done
```

## Performance

- **Startup Time**: < 1 second
- **Memory Usage**: ~10-20 MB (idle), scales with history size
- **Request Latency**: Depends on Ollama API response time
- **Concurrency**: Handles multiple concurrent requests safely

## Monitoring and Observability

### Logging

The application logs all HTTP requests with:
- Method
- URI
- Status code
- Response time
- Client IP

Example log entry:
```
2025/12/10 10:30:45 POST /ask 200 1.234s 192.168.1.100
```

### Health Check

A health check can be added to the Docker configuration:
```yaml
healthcheck:
  test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:8080/ || exit 1"]
  interval: 30s
  timeout: 10s
  retries: 3
```

### OpenTelemetry (Optional)

Enable distributed tracing by setting:
```bash
ENABLE_OTEL=true
OTEL_EXPORTER_OTLP_ENDPOINT=http://your-otel-collector:4317
```

## Troubleshooting

### Common Issues

1. **Cannot connect to Ollama**
   - Verify `OLLAMA_URL` is correct
   - Ensure Ollama instance is accessible
   - Check firewall rules

2. **Rate limit errors**
   - Increase `RATE_LIMIT` environment variable
   - Check if multiple clients are using the same IP

3. **Memory usage growing**
   - Decrease `MAX_HISTORY_SIZE`
   - Check for memory leaks (shouldn't happen with proper implementation)

4. **Slow response times**
   - Check Ollama API performance
   - Reduce `MAX_TOKENS` for faster responses
   - Adjust `OLLAMA_TIMEOUT` if needed

## Development

### Code Structure

- **main.go**: Application bootstrapping and server lifecycle
- **config.go**: Environment variable parsing and defaults
- **handlers.go**: HTTP request handlers and business logic
- **middleware.go**: Cross-cutting concerns (logging, security, rate limiting)
- **storage.go**: Data persistence layer (in-memory)
- **ollama.go**: External API integration

### Adding New Features

1. Add configuration to `config.go`
2. Implement handler in `handlers.go`
3. Register route in `main.go`
4. Add middleware if needed in `middleware.go`
5. Update this README

### Code Quality

- **Linting**: Use `golangci-lint`
- **Formatting**: Use `gofmt`
- **Vet**: Use `go vet`

## License

This project is provided as-is for educational and entertainment purposes.

## Credits

- Original Python implementation by Sasha Karcz
- Go conversion with security enhancements by Claude Sonnet 4.5
- Ollama for AI inference

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## Support

For issues or questions:
- Open an issue on GitHub
- Contact: sasha@starnix.net

## Changelog

### Version 2.0 (Go Rewrite)
- Converted from Python/Flask to Go
- Implemented in-memory storage
- Added comprehensive security features
- Implemented rate limiting
- Added graceful shutdown
- Improved error handling
- Made fully configurable via environment variables
- Added proper logging
- Dockerized with multi-stage build

### Version 1.0 (Python)
- Initial Python/Flask implementation
- Basic Ollama integration
- JSON file storage
- Simple web interface
