FROM node:20-alpine AS builder
WORKDIR /app

COPY package*.json ./
RUN npm install --legacy-peer-deps

# Устанавливаем конкретные версии
RUN npm install -g @bufbuild/buf@1.38.0 @bufbuild/protoc-gen-es@1.10.1 @connectrpc/protoc-gen-connect-es@1.7.0

COPY . .

RUN rm -rf src/gen/ && buf generate --template buf.gen.frontend.yaml ./api
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]