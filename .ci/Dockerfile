FROM golang:1.23-bookworm

LABEL maintainer="Oleg Kiselev <olegsergeevichkiselev@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/kiselevos/new_tax"

# --- System deps ---
RUN apt-get update && apt-get install -y --no-install-recommends \
    unzip curl git ca-certificates protobuf-compiler make bash jq \
    && rm -rf /var/lib/apt/lists/*

# --- Install Go tools ---
ENV PATH="/go/bin:${PATH}"

# protoc plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0

# linters — фиксированные версии
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

# modern proto tool (buf)
RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.57.2/buf-Linux-x86_64 \
    -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf

# gRPC test tool — версия 1.8.9
RUN curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.8.9/grpcurl_1.8.9_linux_x86_64.tar.gz \
    -o grpcurl.tar.gz && tar -xzf grpcurl.tar.gz && mv grpcurl /usr/local/bin/ && rm grpcurl.tar.gz

# --- Copy and prefetch Go modules ---
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# --- Sanity checks ---
RUN go version && protoc --version && \
    golangci-lint --version && grpcurl --version && buf --version && \
    echo "✅ CI image ready"
