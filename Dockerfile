# syntax=docker/dockerfile:1

# Step 1: Build the executable
FROM golang:alpine AS builder

WORKDIR /app

# Install git and certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Download dependencies (including local replace module lib/goshazam)
COPY go.mod go.sum ./
COPY lib/ ./lib/
RUN go mod download

# Copy source code and compile
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bot .

# Step 2: Build the minimal runtime image
FROM alpine:latest

WORKDIR /app

# Install SSL certificates, timezone data and ffmpeg (for audio/video decoding)
RUN apk --no-cache add ca-certificates tzdata ffmpeg

# Copy the binary from the builder
COPY --from=builder /app/bot /app/bot

# Run as non-root user for security
RUN adduser -D -g '' appuser && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app
USER appuser

ENTRYPOINT ["/app/bot"]
