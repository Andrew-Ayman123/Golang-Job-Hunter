# Single-stage build for Go application
FROM golang:1.24.4-alpine

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag for swagger generation and goose for migrations
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy source code
COPY . .

# Generate swagger documentation
RUN swag init -g cmd/jobhunter/main.go -o ./docs

# Build the application from cmd/jobhunter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/jobhunter

# Create a startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Running database migrations..."' >> /app/start.sh && \
    echo 'goose up' >> /app/start.sh && \
    echo 'echo "Starting application..."' >> /app/start.sh && \
    echo 'exec ./main' >> /app/start.sh && \
    chmod +x /app/start.sh

# Expose port (default, can be overridden)
EXPOSE 8080

# Run migrations then start the application
CMD ["/app/start.sh"]