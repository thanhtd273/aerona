# Use official Golang image as base
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/api/main.go

# Use minimal alpine image for runtime
FROM alpine:3.20

# Add some basic security
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port (adjust as needed)
EXPOSE 8080

# Run the server
CMD ["./server"]