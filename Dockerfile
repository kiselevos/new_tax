FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/tax ./cmd

# runtime
FROM debian:bookworm-slim

WORKDIR /app

# Устанавливаем базовые зависисмости в систему
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/tax .

EXPOSE 50051

CMD ["./tax"]

