# ── Build Stage ────────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o bot .

# ── Runtime Stage ──────────────────────────────────────────────────────────────
FROM alpine:3.20

RUN apk add --no-cache sqlite-libs ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/bot .

RUN mkdir -p /app/data

ENV CONNECT_METHOD=qr \
    DB_PATH=/app/data/sessions.db \
    BOT_PREFIX=! \
    BOT_NAME=GoBot \
    AUTO_RECONNECT=true

VOLUME ["/app/data"]

CMD ["./bot"]
