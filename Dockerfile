# Build stage
FROM golang:1.22-alpine AS builder

# Install ca-certificates and git
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o turso-migrate ./cmd/turso-migrate

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/turso-migrate .

# Create migrations directory
RUN mkdir -p /migrations

# Set default migrations directory
ENV MIGRATIONS_DIR=/migrations

# Expose the binary
ENTRYPOINT ["./turso-migrate"]
CMD ["--help"]