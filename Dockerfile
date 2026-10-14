# syntax=docker/dockerfile:1

# Step 1: Build the executable
FROM golang:alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bot .

# Step 2: Build the minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install SSL certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder
COPY --from=builder /app/bot /app/bot

# Run as non-root user for security
RUN adduser -D -g '' appuser
USER appuser

ENTRYPOINT ["/app/bot"]
