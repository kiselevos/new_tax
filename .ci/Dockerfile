FROM golang:1.23-bookworm

LABEL maintainer="Oleg Kiselev <olegsergeevichkiselev@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/kiselevos/new_tax"

# --- System deps ---
RUN apt-get update && apt-get install -y --no-install-recommends \
    unzip curl git ca-certificates protobuf-compiler make bash jq \
    && rm -rf /var/lib/apt/lists/*

# --- Node.js + npm (для фронтенда) ---
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest && \
    node -v && npm -v

# --- Go tools ---
ENV PATH="/go/bin:${PATH}"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.5 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0 && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

# --- Buf + grpcurl ---
RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.57.2/buf-Linux-x86_64 \
    -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf

RUN curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.8.9/grpcurl_1.8.9_linux_x86_64.tar.gz \
    -o grpcurl.tar.gz && tar -xzf grpcurl.tar.gz && mv grpcurl /usr/local/bin/ && rm grpcurl.tar.gz

# --- Prefetch Go modules ---
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

RUN go version && node -v && npm -v && buf --version && echo "✅ Full CI image ready"
