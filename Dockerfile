# Build Stage: Go Backend
FROM golang:1.25-bookworm as builder

WORKDIR /app

# Copy backend code
COPY backend/ ./backend/

# Build Go Server
WORKDIR /app/backend
RUN go mod download
RUN go build -o server main.go

# Runtime Stage
FROM debian:bookworm-slim

WORKDIR /app

# Install CA certificates for external API calls
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Copy built backend from builder
COPY --from=builder /app/backend/server /app/server

# Copy Frontend Assets 
# NOTE: This assumes 'flutter build web' has been run and exists in frontend/build/web
COPY frontend/build/web /app/frontend/build/web

# Environment Variables
ENV PORT=8080

# Run Server
WORKDIR /app
CMD ["./server"]
