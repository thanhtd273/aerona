# Use official Golang image as base (Debian-based)
FROM golang:1.24 AS builder

# Install necessary dependencies for confluent-kafka-go
RUN apt-get update && apt-get install -y librdkafka-dev pkg-config

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
RUN GOOS=linux go build -a -o server ./cmd/api/main.go

# Use minimal Debian-based image for runtime
FROM debian:bookworm-slim

# Add some basic security and librdkafka for runtime
RUN apt-get update && apt-get install -y ca-certificates librdkafka1 && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port (adjust as needed)
EXPOSE 8080

# Run the server
CMD ["./server"]