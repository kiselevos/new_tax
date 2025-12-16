# Tax Calculator - налоговый калькулятор (Go + ConnectRPC)

**Tax Calculator** - это сервис для расчёта подоходного налога (НДФЛ) в России по правилам, действующим с **2025 года**.  
Он учитывает **прогрессивную шкалу**, **северные надбавки**, **районные коэффициенты** и **статус налогового резидентства**.
Сервис включает UI (web-приложение), API для программного доступа и систему наблюдения, охватывающая логи, метрики и визуализацию.

## Архитектура проекта

Проект состоит из двух независимых Go-приложений, а также выделенного инфраструктурного слоя для наблюдаемости и мониторинга:

### Backend
- Написан на **Go 1.23**  
- Использует **ConnectRPC** (совместимый с gRPC и gRPC-Web стек)  
- В основе - **Protocol Buffers** (`google.golang.org/protobuf`)  
- Отвечает за бизнес-логику, расчёты и API  
- Логирование реализовано через собственный модуль **pkg/logx**
- Поддерживаются модульные тесты на **testify**

### Frontend (BFF)
- Реализован на Go Templates (html/template)
- Представляет собой отдельное веб-приложение с собственным go.mod
- Выступает в роли BFF (Backend-for-Frontend)
- Использует Connect Web Client для связи с backend через ConnectRPC
- Не содержит бизнес-логики — все вычисления выполняются в backend
- Экспортирует Prometheus-метрики и пишет структурированные логи
- UI: CSS / Vanilla JS

### API: public и private ручки
- Backend предоставляет два типа API:
    public — для веб-интерфейса и публичного доступа
    private — для внутренних интеграций и автоматизации
- Доступ к private-методам защищён API-ключом, проверяемым в middleware backend-сервера
- Для public и private API настроен раздельный rate-limit
- Все методы реализованы поверх ConnectRPC и описаны в .proto (tax.TaxService)

### Observability & Runtime Infrastructure
- Логи backend и frontend пишутся в stdout в формате JSON
- Promtail собирает container logs и отправляет их в Loki
- Prometheus собирает метрики сервисов
- Grafana используется как единая точка:
    просмотра логов
    анализа метрик
    диагностики ошибок и производительности
- Инфраструктура вынесена в отдельный слой (infra/) и не влияет на бизнес-логику
---

## Технологический стек

### Backend
- **Go 1.23.10**  
- **ConnectRPC** (gRPC + gRPC-Web совместимый стек)  
- **Protocol Buffers** (`google.golang.org/protobuf v1.34.2`)  
- **pkg/logx** + **log/slog** - кастомная система логирования  
- **testify** - модульное тестирование  
- **oapi-codegen/runtime** - вспомогательные утилиты для сериализации и API  
- **Docker**, **Makefile** - контейнеризация и автоматизация сборки  

### Frontend
- **Go Templates (html/template)** - серверная генерация HTML-страниц  
- **Connect Web Client** (`pb.NewTaxServiceClient(conn)`) - связь с backend через ConnectRPC  
- **godotenv** - загрузка переменных окружения  
- **CSS / Vanilla JS** - фронтенд без React, через шаблоны  
- **Docker**, **Makefile** - сборка и деплой web-приложения 

### Observability & Infrastructure
- Prometheus — сбор и хранение метрик
- Loki — централизованное хранилище логов
- Promtail — доставка container logs → Loki
- Grafana — визуализация метрик и логов
- Docker Compose — оркестрация runtime и infra-сервисов

## Возможности

### Бизнес-логика расчёта
- Расчёт НДФЛ по прогрессивной шкале, действующей с 2025 года (5 уровней ставок)
- Отдельная шкала для северных надбавок и районных коэффициентов
- Поддержка статуса налогового нерезидента (единая ставка 30 %)
- Помесячное округление налога по правилам НК РФ  
  (п. 6 ст. 52 — 50 коп. отбрасывается, ≥ 50 коп. округляется вверх)
- Подробная помесячная детализация:
  доход / применённая ставка / налог / net-доход

### Платформенные возможности
- Единая бизнес-логика для UI, публичного API и интеграций
- Встроенные middleware:
  - структурированное логирование
  - graceful shutdown
  - panic recovery
- Поддержка public / private API с раздельным rate-limit
---

## Как считается налог

- **Оклад** (`gross_salary`) указывается в копейках.  
- **Коэффициенты** (`territorial_multiplier`, `northern_coefficient`) - в процентах (100–200).
- **Округление** выполняется каждый месяц: 50 коп - отбрасывается, ≥ 50 коп - округляется вверх (п. 6 ст. 52 НК РФ).  
- **Период**: расчёт с указанного месяца → до декабря текущего года.  
- **Льготы** (`has_tax_privilege = true`) - применяется упрощённая шкала (13 % / 15 %).  
- **Нерезиденты** (`is_not_resident = true`) - единая ставка 30 %.  

## Быстрый старт

### Зависимости
- **Go 1.23.10**  
- **protoc 3.21.12**  
  - **protoc-gen-go v1.34.2**  
  - **protoc-gen-connect-go v1.16.0**  
- **make 4.4.1**  
- **docker 27.1.2**  
- **docker-compose 2.29.2**  
- *(опционально)* `grpcurl 1.9.1`, `golangci-lint 1.60.3` - для отладки и линтинга  ---


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

# === Logging ===
LOG_LEVEL=info          # debug | info | warn | error
LOG_MODE=json           # json | text

# === Security ===
# API-ключ для доступа к private API
API_KEY=api_7f28c3a4b1ef49d9a6c1d742e91f35c2

# === Rate limiting ===
# Public API
RATE_LIMIT_PUBLIC_RPS=1        # запросов в секунду
RATE_LIMIT_PUBLIC_BURST=10

# Private API
RATE_LIMIT_PRIVATE_RPS=2       # запросов в секунду
RATE_LIMIT_PRIVATE_BURST=20
```

- Файл /web/.env:
```bash
# === API Version ===
API_VERSION=v1

# === Web server ===
WEB_PORT=8080

# === Backend connection ===
BACKEND_ADDR=localhost:50051

# === Business constraints ===
MIN_ALLOWED_SALARY=3000    # минимальный ввод оклада, в рублях
MIN_LIVING_WAGE=2440000    # МРОТ, в копейках

# === Logging ===
LOG_LEVEL=info
LOG_MODE=json

# === Feedback ===
FEEDBACK_EMAIL=olegsergeevichkiselev@gmail.com

# === Geo / Metrics ===
# Путь к CSV с IP → регион
GEOIP_CSV_PATH=data/ip.csv
```

## Запуск

### 🐳 Через Docker

```bash
# собрать и запустить все контейнеры
docker compose up --build

# или через Makefile
make docker-build
```
>Все сервисы поднимаются автоматически: web доступен на 8080, backend - на 50051.

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

### Observability & Infrastructure (Prometheus, Loki, Grafana)

Для мониторинга и анализа логов используется отдельный инфраструктурный стек, вынесенный в infra/.

Запуск observability-стека:
```bash
docker compose -f infra/docker-compose.yaml up -d

# или через Makefile
make docker-infra-up
```
После запуска будут доступны:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
> Grafana автоматически подхватывает источники данных и дашборды через provisioning.

### gRPC API (из .proto)
#### Сервис: tax.TaxService
  - CalculatePublic(CalculatePublicRequest) → CalculatePublicResponse
  - CalculatePrivate(CalculatePrivateRequest) → CalculatePrivateResponse
  - Healthz(HealthzRequest) → HealthzResponse
#### Основные поля:
  - gross_salary - оклад в копейках
  - territorial_multiplier, northern_coefficient - 100–200 %
  - start_date - дата начала периода
  - has_tax_privilege, is_not_resident - флаги статуса


### Проектная структура
  - Монолит с модульной архитектурой, без жёсткой связности
  - Proto-first подход — единый контракт для UI, API и интеграций
  - BFF-слой (web) — thin client к backend через gRPC
  - Observability вынесена в infra/, может запускаться отдельно
  - Одинаковое поведение логов и метрик в local / Docker / prod

```bash
├── ABOUT.md                 # Описание проекта и его назначения
├── README.md                # Основная документация
├── contribute.md            # Правила контрибьюции
├── Dockerfile               # Docker-сборка backend-сервиса
├── docker-compose.yaml      # Запуск backend + web (runtime окружение)
├── Makefile                 # Сборка, запуск, тесты, codegen

├── cmd/
│   └── main.go              # Точка входа backend gRPC / ConnectRPC сервиса

├── docs/
│   └── grpc/
│       └── tax.proto        # Proto-контракты (proto-first подход)

├── gen/                     # Сгенерированный код из .proto

├── internal/                # Внутренняя логика backend (не экспортируется)
│   ├── calculate/           # Доменная логика расчёта НДФЛ
│   ├── config/              # Конфигурация backend-приложения
│   ├── middleware/          # gRPC middleware (auth, logging, rate-limit, recovery)
│   └── server/              # Инициализация и запуск gRPC / ConnectRPC сервера

├── pkg/                     # Переиспользуемые пакеты
│   └── logx/                # Кастомный структурированный логгер (slog)

├── test/
│   └── server_integ_test.go # Интеграционные тесты backend

├── infra/                   # Observability и инфраструктура
│   ├── docker-compose.yaml  # Отдельный compose для infra (Prometheus, Loki, Grafana)
│   ├── prometheus/
│   │   └── prometheus.yml   # Конфигурация Prometheus (scrape targets)
│   ├── loki/
│   │   └── config.yaml      # Конфигурация Loki
│   ├── promtail/
│   │   └── config.yaml      # Сбор container logs → Loki
│   └── grafana/
│       ├── dashboards/      # Готовые Grafana dashboards
│       └── provisioning/    # Автопровижининг Grafana

├── web/                     # Frontend (BFF слой)
│   ├── Dockerfile           # Docker-сборка web-сервиса
│   ├── Makefile             # Сборка, запуск, тесты
│   ├── go.mod
│   ├── go.sum

│   ├── cmd/
│   │   └── web.go           # Точка входа web-сервера

│   ├── api_docs/            # Swagger / API-документация
│   │   └── swagger.json
│   ├── api_docs_embed.go    # Embed API-доков в бинарник

│   ├── handlers/            # HTTP-хендлеры и маршруты
│   ├── internal/            # Внутренняя логика web-приложения
│   │   ├── api/             # Public / Private HTTP API (thin wrapper)
│   │   ├── client/          # gRPC / ConnectRPC клиент к backend
│   │   ├── config/          # Конфигурация web
│   │   ├── geoip/           # GeoIP / региональная логика
│   │   ├── metrics/         # Экспорт Prometheus-метрик
│   │   ├── middleware/      # HTTP middleware (logging, metrics, CORS)
│   │   └── server/          # HTTP-сервер web-приложения

│   ├── data/                # Вспомогательные данные (CSV, метрики регионов)
│   ├── static/              # CSS, JS, изображения, sitemap, robots.txt
│   ├── templates/           # Go templates (UI)
│   ├── template_funcs.go    # Вспомогательные функции шаблонов
│   └── testutils/           # Тестовые заглушки и fake-клиенты
```

### Осеновные команды Make
```bash
make setup          # установка зависимостей (Go + frontend)
make codegen        # генерация gRPC/Connect-кода
make run-all        # запуск backend + frontend локально
make docker-build   # сборка и запуск контейнеров
make test-all       # запуск всех тестов
make ci             # запуск проверки перед push
```
