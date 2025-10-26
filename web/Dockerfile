# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app/web

# Показываем, что Docker реально видит контекст
RUN ls -la /app

# Копируем исходники веба и API
COPY web/ . 
COPY gen/ ../gen/
RUN go mod tidy
RUN go build -o web ./cmd/web.go

# --- финальный образ ---
FROM alpine:3.20
WORKDIR /app/web
COPY web/.env .env 
COPY --from=builder /app/web/web .
COPY web/templates/ templates/
COPY web/static/ static/

EXPOSE 8080
CMD ["./web"]
