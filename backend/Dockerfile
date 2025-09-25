# 1) Build stage
FROM golang:1.23 AS builder
WORKDIR /app


# Install build deps using apt (for Debian-based golang image)
RUN apt-get update && apt-get install -y --no-install-recommends git && apt-get clean


COPY go.mod go.sum ./
RUN go mod download


COPY . .


# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/server ./cmd/server


# 2) Runtime stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server /app/server

# Non-root user
RUN adduser -D appuser
USER appuser

# Environment variables (optional, adjust as needed)
ENV PORT=8080

# Expose port
EXPOSE ${PORT}

# Healthcheck (optional, add if you have a health endpoint)
HEALTHCHECK --interval=30s --timeout=3s \
    CMD wget --no-check-certificate --quiet --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Entry point
ENTRYPOINT ["/app/server"]