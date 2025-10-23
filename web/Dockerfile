# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . /app/web
COPY ../go.mod ../go.sum /app/

RUN go mod download
RUN go build -o /web cmd/web.go

# --- final image ---
FROM alpine:3.20
WORKDIR /web
COPY --from=builder /web /web
COPY templates/ templates/
COPY static/ static/

EXPOSE 8080
CMD ["./web"]
