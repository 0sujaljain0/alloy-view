# =========================================================
# Stage 1: The Builder
# =========================================================
FROM golang:1.25.1-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

# Dependency Caching
COPY go.mod go.sum ./
RUN go mod download

# Copy Source
COPY . .

# Build Static Binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main.go

# =========================================================
# Stage 2: The Runner
# =========================================================
FROM alpine:latest

# 1. Set the folder where everything will live
WORKDIR /app

# 2. Install Tools
RUN apk add --no-cache \
    curl \
    bind-tools \
    iproute2 \
    ca-certificates

# 3. Copy Binary AND Config to the SAME LEVEL (/app)
COPY --from=builder /app/main .
COPY config.yaml .
COPY static static

EXPOSE 8093

# 5. Run with "./" to ensure it looks in the current folder
ENTRYPOINT ["./main"]
