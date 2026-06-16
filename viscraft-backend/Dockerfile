# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install git for fetching dependencies (some may use git)
RUN apk add --no-cache git

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o viscraft-backend .

# Stage 2: Production image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS calls (e.g., Gemini API)
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from builder
COPY --from=builder /build/viscraft-backend /app/viscraft-backend

# Create persistent storage directory
RUN mkdir -p /app/storage/images

# Expose the application port
EXPOSE 8080

# Define volume for persistent image storage
VOLUME ["/app/storage/images"]

# Run the application
CMD ["/app/viscraft-backend"]
