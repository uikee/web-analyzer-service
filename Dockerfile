FROM golang:1.23.5-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy application source code
COPY . .

# Build the application
RUN go build -o web-analyzer-service cmd/main.go

# Use a minimal base image
FROM alpine:latest

# Set working directory in container
WORKDIR /root/

# Copy the built binary from builder stage
COPY --from=builder /app/web-analyzer-service .

# Expose the application port
EXPOSE 8081

# Command to run the application
CMD ["./web-analyzer-service"]