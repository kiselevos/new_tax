# syntax=docker/dockerfile:1
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./

# Добавляем protobuf-модуль
COPY pkg/logx/ ./pkg/logx/
COPY gen/ ./gen/

RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o /app/bin/tax ./cmd/main.go

# --- финальный образ ---
FROM debian:bookworm-slim

WORKDIR /app

# Минимальный runtime
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/tax .

EXPOSE 50051
CMD ["./tax"]