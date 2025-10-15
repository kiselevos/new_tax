FROM golang:1.23-bookworm

LABEL maintainer="Oleg Kiselev <olegsergeevichkiselev@gmail.com>"
LABEL org.opencontainers.image.source="https://github.com/kiselevos/new_tax"

# --- System deps ---
RUN apt-get update && apt-get install -y --no-install-recommends \
    unzip curl git ca-certificates make bash jq \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# --- Node.js + npm ---
RUN curl -fsSL https://deb.nodesource.com/setup_22.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest pnpm@latest

# --- Go tools с совместимыми версиями для Go 1.23 ---
ENV PATH="/go/bin:${PATH}"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0 && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3 && \
    go install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.16.0

# --- Buf + grpcurl ---
RUN curl -sSL https://github.com/bufbuild/buf/releases/download/v1.57.2/buf-Linux-x86_64 \
    -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf

RUN curl -sSL https://github.com/fullstorydev/grpcurl/releases/download/v1.8.9/grpcurl_1.8.9_linux_x86_64.tar.gz \
    | tar -xz -C /usr/local/bin/ grpcurl && chmod +x /usr/local/bin/grpcurl

# --- Prefetch Go modules ---
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# --- Prefetch Node deps ---
WORKDIR /app/web
COPY web/package*.json web/pnpm-lock.yaml* ./
RUN if [ -f "pnpm-lock.yaml" ]; then pnpm install --frozen-lockfile; \
    else npm ci --legacy-peer-deps || npm install; fi

# --- Verify toolchain ---
WORKDIR /app
RUN go version && \
    node -v && npm -v && pnpm -v && \
    buf --version && grpcurl --version && \
    echo "✅ Full CI image ready with Go, Node, Buf, ConnectRPC"