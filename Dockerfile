# Stage 1: Build the Go binary
FROM golang:alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk add --no-cache git ca-certificates

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Final minimal image
FROM alpine:3.21

WORKDIR /app

# Install tzdata for TimeZone resolution and copy ca-certificates
RUN apk add --no-cache tzdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from the builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8082

# Command to run the application
CMD ["./main"]
