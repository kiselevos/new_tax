# --- Build stage ---
FROM golang:1.23-bookworm AS builder
WORKDIR /app

# Кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код и собираем бинарь
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o tax ./cmd

# --- Runtime stage ---
FROM gcr.io/distroless/base-debian12 AS runtime
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/tax .

# Backend слушает gRPC (50051) и gRPC-Web (8081)
EXPOSE 50051
EXPOSE 8081

# Запускаем
ENTRYPOINT ["/app/tax"]
