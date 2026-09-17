# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary from ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server

# Final minimal runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy executable binary from builder stage
COPY --from=builder /app/server .

# Expose container port
EXPOSE 8080

# Environment defaults
ENV PORT=8080

# Entrypoint command
CMD ["./server"]
