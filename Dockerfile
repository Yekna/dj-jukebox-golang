# ----------------------------------------------------
# Stage 1 — Build the Go binary
# ----------------------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# ----------------------------------------------------
# Stage 2 — Run with minimal Alpine image
# ----------------------------------------------------
FROM alpine:latest

WORKDIR /app

# Add certificates for HTTPS (YouTube API, Mongo Atlas, etc.)
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .
# COPY .env .env

EXPOSE 8080

CMD ["./server"]
