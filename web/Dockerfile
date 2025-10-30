# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем только нужные каталоги
COPY web/ web/
COPY pkg/ pkg/
COPY gen/ gen/

# Сборка веб-модуля
RUN cd web && go mod tidy && go build -o web ./cmd/web.go

# --- финальный образ ---
FROM alpine:3.20
WORKDIR /app/web

COPY --from=builder /app/web/web .
COPY --from=builder /app/web/templates ./templates
COPY --from=builder /app/web/static ./static
COPY web/.env .env

EXPOSE 8080
CMD ["./web"]