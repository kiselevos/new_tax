FROM golang:1.23-bookworm AS builder

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.16.0

RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.38.0/buf-Linux-x86_64 -o /usr/local/bin/buf && \
    chmod +x /usr/local/bin/buf

ENV PATH=$PATH:/go/bin

WORKDIR /app
COPY . .

RUN buf generate --template buf.gen.backend.yaml ./api

# СОЗДАЕМ ПРАВИЛЬНУЮ СТРУКТУРУ ДЛЯ ИМПОРТОВ
RUN mkdir -p gen/grpc/api && \
    mv gen/grpc/tax.pb.go gen/grpc/api/ && \
    mv gen/grpc/taxconnect gen/grpc/api/

RUN go mod download
RUN go build -o /tax ./cmd

FROM gcr.io/distroless/base-debian12
COPY --from=builder /tax /
EXPOSE 8081
CMD ["/tax"]