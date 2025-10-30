# Tax Calculator — налоговый калькулятор (Go + ConnectRPC)

**Tax Calculator** — это сервис для расчёта подоходного налога (НДФЛ) в России по правилам, действующим с **2025 года**.  
Он учитывает **прогрессивную шкалу**, **северные надбавки**, **районные коэффициенты** и **статус налогового резидентства**.  

## Архитектура проекта

Проект состоит из двух независимых Go-приложений:

### Backend
- Написан на **Go 1.23**  
- Использует **ConnectRPC** (совместимый с gRPC и gRPC-Web стек)  
- В основе - **Protocol Buffers** (`google.golang.org/protobuf`)  
- Отвечает за бизнес-логику, расчёты и API  
- Логирование реализовано через собственный модуль **pkg/logx**  
- Поддерживаются модульные тесты на **testify**

### Frontend
- Реализован на **Go Templates (html/template)**  
- Представляет собой отдельное веб-приложение с собственным `go.mod` и зависимостями  
- Использует **Connect Web Client** для связи с backend через ConnectRPC  
- Стили и интерфейс - CSS / Vanilla JS
---

## Технологический стек

### Backend
- **Go 1.23.10**  
- **ConnectRPC** (gRPC + gRPC-Web совместимый стек)  
- **Protocol Buffers** (`google.golang.org/protobuf v1.34.2`)  
- **pkg/logx** + **log/slog** — кастомная система логирования  
- **testify** — модульное тестирование  
- **oapi-codegen/runtime** — вспомогательные утилиты для сериализации и API  
- **Docker**, **Makefile** — контейнеризация и автоматизация сборки  

### Frontend
- **Go Templates (html/template)** — серверная генерация HTML-страниц  
- **Connect Web Client** (`pb.NewTaxServiceClient(conn)`) — связь с backend через ConnectRPC  
- **godotenv** — загрузка переменных окружения  
- **CSS / Vanilla JS** — фронтенд без React, через шаблоны  
- **Docker**, **Makefile** — сборка и деплой web-приложения 

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

### Зависимости
- **Go 1.23.10**  
- **protoc 3.21.12**  
  - **protoc-gen-go v1.34.2**  
  - **protoc-gen-connect-go v1.16.0**  
- **make 4.4.1**  
- **docker 27.1.2**  
- **docker-compose 2.29.2**  
- *(опционально)* `grpcurl 1.9.1`, `golangci-lint 1.60.3` — для отладки и линтинга  ---


### 🔧 Установка и генерация
```bash
git clone https://github.com/kiselevos/new_tax
cd new_tax
```
#### (Опционально) Установить инструменты разработчика
Если вы хотите запускать линтеры, тесты или генерировать gRPC-код:
```bash
make setup-tools
# Инструменты будут установлены в .tools/bin.
```
#### Подтянуть зависимости
```bash
# Backend
go mod tidy
go mod download

# Frontend
cd web && go mod tidy
go mod download

# Если установлен make:
make setup

```
---

### Конфигурация приложения (пример .env)
- Файл .env в корне проекта:
```bash
# === Backend ===
BACKEND_PORT=50051

# Уровень и формат логов
LOG_LEVEL=info
LOG_MODE=json
```
- Файл /web/.env:
```bash
# === Frontend ===
WEB_PORT=:8080
# === Conditions salary ===
MIN_ALLOWED_SALARY=3000 # Минимальный оклад.
MIN_LIVING_WAGE=2440000 # МРОТ в копейках.
# === Logger ===
LOG_LEVEL=info
LOG_MODE=json
```

## Запуск

### 🐳 Через Docker

```bash
# собрать и запустить все контейнеры
docker compose up --build

# или через Makefile
make docker-build
```
>Все сервисы поднимаются автоматически: web доступен на 8080, backend — на 50051.

### 💻 Без Docker (локально)

```bash
go run ./cmd/main.go &
cd web && go run ./cmd/web.go

# или через Makefile
make run # запуск backend
cd web && make run # запуск frontend
make run-all # запуск одной командой 
```
> После запуска:
  Frontend: http://localhost:8080
  Backend: http://localhost:50051

### Проверка доступности
```bash
grpcurl -plaintext localhost:50051 tax.TaxService/Healthz
```

### Пример вызова
```bash
grpcurl -plaintext -d '{
  "gross_salary": 20000000,
  "territorial_multiplier": 110,
  "northern_coefficient": 130,
  "start_date": "2025-06-01T00:00:00Z",
  "has_tax_privilege": false,
  "is_not_resident": false
}' localhost:50051 tax.TaxService/CalculatePrivate
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
├── cmd/ # Точка входа backend-приложения
│ └── main.go
├── internal/ # Внутренняя логика (не экспортируется наружу)
│ ├── calculate/ # Модуль бизнес-логики и расчёта НДФЛ
│ └── server/ # Инициализация и запуск gRPC / ConnectRPC-сервера
├── pkg/ # Переиспользуемые пакеты (общие утилиты)
│ ├── helpers/ # Вспомогательные функции
│ └── logx/ # Кастомный логгер на основе slog
├── gen/ # Сгенерированные файлы из .proto
│ ├── grpc/ # gRPC и ConnectRPC артефакты
│ ├── go.mod, go.sum # Отдельный модуль для генерации
├── web/ # Отдельное frontend-приложение
│ ├── cmd/ # Точка входа веб-сервера (Go templates)
│ ├── handlers/ # HTTP-обработчики и маршруты
│ ├── internal/ # Внутренняя логика фронта (middleware, модели)
│ ├── static/ # CSS, JS, иконки
│ ├── templates/ # HTML-шаблоны интерфейса
│ ├── template_funcs.go # Встроенные функции для шаблонов
│ ├── Dockerfile, Makefile
│ ├── go.mod, go.sum
├── test/ # Интеграционные и e2e-тесты
│ └── server_integ_test.go
├── tools.go # Фиксация CLI-инструментов в go.mod (protoc, lint и др.)
├── Makefile # Общие команды сборки и генерации
├── Dockerfile # Сборка backend-сервера
├── docker-compose.yaml # Запуск всех сервисов (backend + frontend)
├── docs/grpc/ # Документация по gRPC и API
├── contribute.md # Правила контрибьюции
├── go.mod, go.sum # Основные зависимости backend-а
└── README.md
```

### Осеновные команды Make
```bash
make setup          # установка зависимостей (Go + frontend)
make codegen        # генерация gRPC/Connect-кода
make run-all            # запуск backend + frontend локально
make docker-build   # сборка и запуск контейнеров
make test-all       # запуск всех тестов
```
