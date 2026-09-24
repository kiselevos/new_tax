# Калькулятор НДФЛ

[English](README.md) | **Русский**

Сервис для помесячного расчёта подоходного налога (НДФЛ) и страховых взносов работодателя в России. Подходит для сотрудников, HR и разработчиков, которым нужен точный программируемый расчёт.

## Попробовать

[calculator-ndfl.ru](https://calculator-ndfl.ru)

## Что считает

- **Прогрессивная шкала НДФЛ** (с 2025 года): пять ставок — 13%, 15%, 18%, 20%, 22%
- **Три типа занятости**: трудовой договор (ТД), договор ГПХ, самозанятость (НПД, 4% / 6%)
- **Налоговые вычеты**: стандартные на детей (ст. 218 НК РФ), социальные — лечение и обучение (ст. 219), имущественные — покупка жилья и проценты по ипотеке (ст. 220)
- **Районный коэффициент (РК)** — часть основного дохода, облагается по общей прогрессивной шкале
- **Северная надбавка (СН)** — отдельная налоговая база, облагается по упрощённой шкале 13% / 15%
- **Взносы работодателя**: ПФР (22% до предельной базы / 10% сверх), ФОМС (5,1%), ФСС (2,9% до предельной базы)
- **Разовые премии** по месяцам с корректным переходом через пороги шкалы
- **Льготники силовых ведомств** — упрощённая шкала 13% / 15%
- **Нерезиденты РФ** — единая ставка 30%
- **Помесячная детализация с накопительным итогом (YTD)** за весь налоговый период

Все денежные суммы хранятся в `uint64` в копейках — в расчётах нет чисел с плавающей точкой.
Налог округляется до полных рублей по п. 6 ст. 52 НК РФ (менее 50 копеек отбрасывается, 50 копеек и более — округляется вверх).

## Быстрый запуск

```bash
# Создать локальный конфиг
cp .env.example .env

# Собрать и запустить все контейнеры
docker compose up --build

# Или через Makefile
make docker-build
```

После запуска:
- Веб-интерфейс и REST API: http://localhost:8080
- Backend gRPC-сервер: `localhost:50051` (gRPC без TLS, reflection включён)

Запуск без Docker:

```bash
go run ./cmd/main.go &
cd web && go run ./cmd/web.go

# Или через Makefile
make run-all
```

Проверка доступности:

```bash
grpcurl -plaintext localhost:50051 tax.TaxService/Healthz
```

Пример вызова gRPC API (private-метод требует API-ключ, заданный в `API_KEY` в `.env`):

```bash
grpcurl -plaintext \
  -H 'x-api-key: <API_KEY>' \
  -d '{
    "gross_salary": 20000000,
    "territorial_multiplier": 110,
    "northern_coefficient": 130,
    "start_date": "2025-06-01T00:00:00Z",
    "has_tax_privilege": false,
    "is_not_resident": false
  }' \
  localhost:50051 tax.TaxService/CalculatePrivate
```

Пример вызова REST (JSON) API через веб-сервис:

```bash
curl -X POST http://localhost:8080/api/v1/calc \
  -H 'Content-Type: application/json' \
  -d '{"gross_salary": 20000000, "territorial_multiplier": 120}'
```

Документация OpenAPI (Swagger UI) доступна по адресу `/api-docs`.

### Публичный демо-ключ

Проект открытый и учебный, поэтому ключ private API публичен намеренно. Это `API_KEY` из `.env.example`, и он работает на проде:

```bash
curl -X POST https://calculator-ndfl.ru/api/v1/private-calc \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: api_7f28c3a4b1ef49d9a6c1d742e91f35c2' \
  -d '{"gross_salary": 20000000, "employment_type": "TD", "start_date": "2026-01-01"}'
```

Действует rate limit. Пользуйтесь и проверяйте.

Основные команды Make:

```bash
make setup          # установка зависимостей (backend + web)
make codegen        # генерация Go gRPC-кода из .proto
make run-all        # запуск backend + web локально
make docker-build   # сборка и запуск контейнеров
make test-all       # запуск тестов backend
make ci             # линтер + тесты перед push
```

## Архитектура

```
Браузер ──HTTP──▶ web (BFF, :8080) ──gRPC──▶ backend (:50051)
                   │  HTML-страницы               │  бизнес-логика
                   │  REST API /api/v1/*          │  auth, rate limit
                   │  /metrics                    │
                   ▼                              ▼
             Prometheus ◀── scrape          JSON-логи в stdout ──▶ Promtail ──▶ Loki ──▶ Grafana
```

Проект состоит из двух независимых Go-приложений и инфраструктурного слоя наблюдаемости.

### Backend
- **Go 1.23**, gRPC-сервер на `google.golang.org/grpc`
- **Proto-first**: контракт в `docs/grpc/tax.proto`, Go-код генерируется в `gen/`
- Отвечает за бизнес-логику: расчёт налога, валидацию входных данных, API
- Цепочка interceptor'ов: recovery → логирование (request ID) → проверка API-ключа → rate limit
- Структурированное логирование через `log/slog` (тонкая обёртка `pkg/logx`)
- Модульные и интеграционные тесты на **testify**

### Web (BFF)
- Отдельный Go-модуль со своим `go.mod`
- Серверный рендеринг на Go `html/template`, vanilla CSS и JavaScript
- Выступает в роли BFF (Backend-for-Frontend): обращается к backend через gRPC-клиент
- Не содержит логики расчёта — все вычисления выполняются в backend
- Публичный JSON REST API (`/api/v1/calc`, `/api/v1/private-calc`) с описанием в OpenAPI
- Экспортирует Prometheus-метрики (`/metrics`) и пишет структурированные JSON-логи
- GeoIP (CSV-база) для разметки метрик по регионам

### API: public и private методы
- `CalculatePublic` — упрощённый расчёт, доступен всем
- `CalculatePrivate` — полный расчёт (дата начала, премии, вычеты, тип занятости, взносы работодателя); требует заголовок `x-api-key`, который проверяется в interceptor'е backend. Запросы из веб-интерфейса помечаются как внутренний трафик
- Раздельный rate limit для public и private вызовов (token bucket, `golang.org/x/time/rate`)
- Все методы описаны в `.proto` (`tax.TaxService`)

### Observability
- Логи backend и web пишутся в stdout в формате JSON
- Promtail собирает логи контейнеров и отправляет их в Loki
- Prometheus собирает метрики веб-сервиса
- Grafana — единая точка просмотра логов, метрик и бизнес-дашбордов

Запуск observability-стека:

```bash
docker compose -f infra/docker-compose.yaml up -d

# Или через Makefile
make docker-infra-up
```

После запуска:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100

### Структура проекта

```
├── README.md / README.ru.md
├── ABOUT.md                 # Описание продукта
├── contribute.md            # Правила контрибьюции
├── Dockerfile               # Образ backend
├── docker-compose.yaml
├── Makefile
│
├── cmd/
│   └── main.go              # Точка входа backend (gRPC-сервер)
│
├── docs/
│   └── grpc/
│       └── tax.proto        # Proto-контракт (единый источник истины для API)
│
├── gen/                     # Сгенерированный из .proto код (отдельный Go-модуль)
│
├── internal/                # Внутренняя логика backend
│   ├── calculate/           # Доменная логика расчёта НДФЛ
│   ├── config/              # Конфигурация из env
│   ├── middleware/          # gRPC interceptor'ы: auth, логи, recovery, rate limit
│   └── server/              # gRPC-сервер и реализация сервиса
│
├── pkg/
│   └── logx/                # Структурированный логгер на slog
│
├── test/
│   └── server_integ_test.go # Интеграционные тесты gRPC-сервера
│
├── infra/                   # Observability-стек
│   ├── docker-compose.yaml
│   ├── prometheus/
│   ├── loki/
│   ├── promtail/
│   └── grafana/
│
├── web/                     # BFF (отдельный Go-модуль)
│   ├── Dockerfile
│   ├── Makefile
│   ├── cmd/web.go           # Точка входа web
│   ├── handlers/            # Обработчики HTML-страниц
│   ├── internal/
│   │   ├── api/             # REST API: обработчики и DTO
│   │   ├── client/          # gRPC-клиент к backend
│   │   ├── config/
│   │   ├── geoip/
│   │   ├── metrics/
│   │   ├── middleware/
│   │   └── server/
│   ├── static/              # CSS, JS, спецификация OpenAPI
│   └── templates/           # Go HTML-шаблоны
│
└── project-docs/            # Документация расчёта, реестр налоговых констант, роадмап
```

## Разработка

Правила контрибьюции — см. [contribute.md](contribute.md).

## Как работает расчёт

Подробная документация алгоритма — см. [project-docs/how-calculation-works.md](project-docs/how-calculation-works.md).
