# Tax Calculator — налоговый калькулятор (Go + React + ConnectRPC)

**Tax Calculator** — это сервис для расчёта подоходного налога (НДФЛ) в России по правилам, действующим с 2025 года.  
Он учитывает прогрессивную шкалу, северные надбавки, районные коэффициенты и статус налогового резидентства.  

Backend реализован на **Go (ConnectRPC)**, а frontend — на **React (Vite)**  
с типизированным клиентом, сгенерированным из `.proto`-контрактов.

> Backend использует **ConnectRPC**, совместимый с gRPC и gRPC-Web,  
> всё работает на одном порту `8081` без отдельного gRPC-порта.  
> Для браузера разрешены CORS-запросы (через middleware `withCORS()`).

---

## Технологический стек

### Backend
- **Go 1.23.10**
- **ConnectRPC** (gRPC + gRPC-Web)
- **Protobuf** (`google.golang.org/protobuf v1.34.2`)
- **log/slog**
- **Buf** для генерации `.proto`
- **Docker**, **Makefile**

### Frontend
- **React + TypeScript**
- **Vite**
- **Connect Web Client** (`@connectrpc/connect`, `@bufbuild/protobuf`)
- **Node.js 20+**
- **Nginx** (в Docker)

---

## Возможности

- Расчёт по новой прогрессивной шкале (5 уровней ставок).  
- Отдельная шкала для северных надбавок.  
- Помесячное округление налога по правилам НК РФ (п. 6 ст. 52).  
- Поддержка статуса нерезидента (единая ставка 30 %).  
- Подробная детализация помесячно (доход, ставка, налог, нетто).  
- Встроенные middleware: логирование, graceful shutdown, recovery.  

---

## Как считается налог

- **Оклад** (`gross_salary`) указывается в копейках.  
- **Коэффициенты** (`territorial_multiplier`, `northern_coefficient`) — в процентах (100–200).
- **Округление** выполняется каждый месяц: 50 коп — отбрасывается, ≥ 50 коп — округляется вверх (п. 6 ст. 52 НК РФ).  
- **Период**: расчёт с указанного месяца → до декабря текущего года.  
- **Льготы** (`has_tax_privilege = true`) — применяется упрощённая шкала (13 % / 15 %).  
- **Нерезиденты** (`is_not_resident = true`) — единая ставка 30 %.  

## Быстрый старт

### 📦 Зависимости
- **Go 1.23.10+**
- **Node.js 20+ и npm 10+**
- **Buf CLI** *(или protoc ≥ 3.21 + protoc-gen-go + protoc-gen-connect-go)*
- (Опционально) `grpcurl`, `golangci-lint`, `make`, `docker`, `docker-compose`
---

### 🔧 Установка и генерация
```bash
git clone https://github.com/kiselevos/new_tax
cd new_tax

# Backend
go mod tidy
buf generate

# Frontend
cd web && npm ci

# Если установлен make:
make setup
```
---

### Настройки окружения
- Создайте файл .env в корне проекта (если отсутствует — применятся значения по умолчанию):
```bash
# === Backend ===
# Порт для ConnectRPC (HTTP + gRPC-Web)
BACKEND_PORT=8081

# Уровень и формат логов
LOG_LEVEL=info
LOG_MODE=json

# === Frontend ===
# Порт для React dev-сервера
FRONTEND_PORT=8080
```

## Запуск

### 🐳 Через Docker

```bash
# собрать и запустить все контейнеры
docker compose up --build

# или через Makefile
make docker-build
```
### 💻 Без Docker (локально)

```bash
go run ./cmd/main.go &
cd web && npm run dev

# или через Makefile
make run
```
> После запуска:
  Frontend: http://localhost:8080
  Backend (ConnectRPC): http://localhost:8081

### Проверка доступности
```bash
curl -X POST http://localhost:8081/tax.TaxService/Healthz \
  -H "Content-Type: application/json" \
  -H "Connect-Protocol-Version: 1" \
  -d '{}'
```

### Пример вызова
```bash
curl -X POST http://localhost:8081/tax.TaxService/CalculatePrivate \
  -H "Content-Type: application/json" \
  -H "Connect-Protocol-Version: 1" \
  -d '{
    "gross_salary": 20000000,
    "territorial_multiplier": 110,
    "northern_coefficient": 130,
    "start_date": "2025-06-01T00:00:00Z",
    "has_tax_privilege": false,
    "is_not_resident": false
  }'
```

### gRPC API (из .proto)
#### Сервис: tax.TaxService
  - CalculatePublic(CalculatePublicRequest) → CalculatePublicResponse
  - CalculatePrivate(CalculatePrivateRequest) → CalculatePrivateResponse
  - Healthz(HealthzRequest) → HealthzResponse
#### Основные поля:
  - gross_salary — оклад в копейках
  - territorial_multiplier, northern_coefficient — 100–200 %
  - start_date — дата начала периода
  - has_tax_privilege, is_not_resident — флаги статуса


### Проектная структура
```bash
├── api/              # .proto контракты
├── cmd/              # точка входа Go-приложения
├── internal/         # бизнес-логика (calculate, server)
├── pkg/              # общие пакеты (helpers, logx)
├── gen/              # сгенерированный код (Connect/gRPC)
├── web/              # фронтенд (React + Vite)
├── Dockerfile        # backend dockerfile
├── docker-compose.yaml
└── Makefile
```

### Осеновные команды Make
```bash
make setup          # установка зависимостей (Go + frontend)
make codegen        # генерация gRPC/Connect-кода
make run            # запуск backend + frontend локально
make docker-build   # сборка и запуск контейнеров
make test-all       # запуск всех тестов
```

