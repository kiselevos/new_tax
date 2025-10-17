# 🤝 Contributing Guide

Спасибо, что решили внести вклад в проект **Tax Calculator** — gRPC-сервис на Go для расчёта прогрессивного подоходного налога в России начиная с 2025 года. Этот документ описывает, как локально поднять окружение, запустить тесты и внести изменения безопасно для продакшена.

---

## 🧩 Общая структура проекта

```text
├── api/              # .proto контракты (Buf)
├── cmd/              # точка входа Go-приложения
├── internal/         # бизнес-логика (calculate, server)
├── pkg/              # вспомогательные пакеты (helpers, logx)
├── gen/              # сгенерированный код (Connect/gRPC)
├── web/              # фронтенд (React + Vite)
├── Dockerfile        # backend dockerfile
├── docker-compose.yaml
└── Makefile
```

---

## ⚙️ Требования к окружению

### Backend
- **Go 1.23.10+**
- **Buf CLI** (или `protoc ≥ 3.21` + `protoc-gen-go` + `protoc-gen-connect-go`)
- **golangci-lint** для статического анализа
- **grpcurl** (опционально, для проверки RPC)

### Frontend
- **Node.js 20+**
- **npm 10+** или **pnpm**
- **Vite** (через `npm run dev`)
- **eslint**, **typescript**

### Опционально
- **Docker** и **docker-compose** — для контейнерной сборки
- **Make** — для удобного запуска команд

---

## 🚀 Быстрый старт для разработки

```bash
git clone https://github.com/kiselevos/new_tax
cd new_tax
```

### Установка зависимостей

```bash
# Backend
go mod tidy
buf generate

# Frontend
cd web && npm ci
```

или одной командой:

```bash
make setup
```

---

### Запуск локально

```bash
# Запустить backend
go run ./cmd/main.go

# Запустить frontend (в отдельном окне)
cd web && npm run dev
```

или одной командой:

```bash
make run
```

> После запуска:  
> - Backend: [http://localhost:8081](http://localhost:8081)  
> - Frontend: [http://localhost:8080](http://localhost:8080)

---

## 🧬 Генерация кода (Buf / ConnectRPC)

Все `.proto`-контракты хранятся в директории `api/`.  
При изменении этих файлов нужно сгенерировать код:

```bash
make codegen
```

или вручную:

```bash
buf generate
```

Сгенерированные файлы попадают в `gen/grpc/`.
> ⚠️ Никогда не редактируйте файлы в `gen/` вручную — они пересоздаются автоматически.

---

## 🧹 Линтинг и форматирование

```bash
# Backend (Go)
make vet           # go vet
make check-fmt     # проверка форматирования
make gofmt         # автоформатирование

# Frontend (React)
make lint          # cd web && npm run lint

# All
make lint-all
```
---

## 🧪 Тестирование

### Запуск всех тестов
```bash
make test-all
```

### Пример интеграционного теста Healthz
```bash
go test -run Test_Server_Healthz ./test
```
---

## 🐳 Работа с Docker

### Собрать и запустить контейнеры
```bash
docker compose up --build
# или
make docker-build
```

### Остановить контейнеры
```bash
make docker-down
```

> После сборки сервис будет доступен:  
> - Backend → http://localhost:8081  
> - Frontend → http://localhost:8080

---

## 🧱 CI/CD

CI выполняет следующие шаги:
- Генерацию gRPC/Connect-кода (`make codegen`)
- Проверку зависимостей (`make tidy`)
- Линтинг (`make lint-all`)
- Тестирование (`make test-all`)
- Сборку frontend и backend

Используется кастомный CI-образ:  
`ghcr.io/kiselevos/tax-ci:v1.5.0`

---

## 🧭 Полезные команды

| Команда | Назначение |
|----------|-------------|
| `make setup` | установить все зависимости |
| `make codegen` | сгенерировать protobuf-код |
| `make run` | запустить backend и frontend |
| `make test-all` | запустить все тесты |
| `make docker-build` | собрать и запустить контейнеры |
| `make docker-down` | остановить контейнеры |

---
