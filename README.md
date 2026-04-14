# new_tax — калькулятор НДФЛ 2025

Сервис для расчёта подоходного налога в России по правилам 2025 года. Подходит для сотрудников, HR и разработчиков, которым нужен точный программируемый расчёт.

## Попробовать

Сервис в разработке — ссылка появится здесь.

## Что считает

- Прогрессивная шкала НДФЛ 2025 года (пять ставок: 13%, 15%, 18%, 20%, 22%)
- Районный коэффициент (РК) — облагается по общей прогрессивной шкале
- Северная надбавка (СН) — отдельная налоговая база, облагается по упрощённой шкале (13% / 15%)
- Взносы работодателя: ПФР (22% / 10% сверх лимита), ФОМС (5,1%), ФСС (2,9% до лимита)
- Льготники силовых ведомств — упрощённая шкала 13% / 15%
- Нерезиденты РФ — единая ставка 30%
- Помесячная детализация с накопительным итогом (YTD) за весь налоговый период

## Быстрый запуск

```bash
# Собрать и запустить все контейнеры
docker compose up --build

# Или через Makefile
make docker-build
```

После запуска:
- Frontend: http://localhost:8080
- Backend: http://localhost:50051

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

Пример вызова API:

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

Основные команды Make:

```bash
make setup          # установка зависимостей (Go + frontend)
make codegen        # генерация gRPC/Connect-кода
make run-all        # запуск backend + frontend локально
make docker-build   # сборка и запуск контейнеров
make test-all       # запуск всех тестов
make ci             # запуск проверки перед push
```

## Архитектура

Проект состоит из двух независимых Go-приложений и инфраструктурного слоя наблюдаемости.

### Backend
- Написан на **Go 1.23**
- Использует **ConnectRPC** (совместимый с gRPC и gRPC-Web стек)
- В основе — **Protocol Buffers** (`google.golang.org/protobuf`)
- Отвечает за бизнес-логику, расчёты и API
- Логирование реализовано через собственный модуль **pkg/logx**
- Поддерживаются модульные тесты на **testify**

### Frontend (BFF)
- Реализован на Go Templates (html/template)
- Отдельное веб-приложение с собственным go.mod
- Выступает в роли BFF (Backend-for-Frontend)
- Использует Connect Web Client для связи с backend через ConnectRPC
- Не содержит бизнес-логики — все вычисления выполняются в backend
- Экспортирует Prometheus-метрики и пишет структурированные логи
- UI: CSS / Vanilla JS

### API: public и private методы
- Backend предоставляет два типа API: public — для веб-интерфейса и публичного доступа, private — для внутренних интеграций и автоматизации
- Доступ к private-методам защищён API-ключом, проверяемым в middleware backend-сервера
- Для public и private API настроен раздельный rate-limit
- Все методы реализованы поверх ConnectRPC и описаны в `.proto` (tax.TaxService)

### Observability
- Логи backend и frontend пишутся в stdout в формате JSON
- Promtail собирает container logs и отправляет их в Loki
- Prometheus собирает метрики сервисов
- Grafana — единая точка просмотра логов, анализа метрик и диагностики

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
├── ABOUT.md
├── README.md
├── contribute.md
├── Dockerfile
├── docker-compose.yaml
├── Makefile

├── cmd/
│   └── main.go              # Точка входа backend gRPC / ConnectRPC сервиса

├── docs/
│   └── grpc/
│       └── tax.proto        # Proto-контракты (proto-first подход)

├── gen/                     # Сгенерированный код из .proto

├── internal/                # Внутренняя логика backend
│   ├── calculate/           # Доменная логика расчёта НДФЛ
│   ├── config/
│   ├── middleware/
│   └── server/

├── pkg/
│   └── logx/                # Кастомный структурированный логгер (slog)

├── test/
│   └── server_integ_test.go

├── infra/                   # Observability и инфраструктура
│   ├── docker-compose.yaml
│   ├── prometheus/
│   ├── loki/
│   ├── promtail/
│   └── grafana/

├── web/                     # Frontend (BFF слой)
│   ├── Dockerfile
│   ├── Makefile
│   ├── cmd/
│   │   └── web.go
│   ├── handlers/
│   ├── internal/
│   │   ├── api/
│   │   ├── client/
│   │   ├── config/
│   │   ├── geoip/
│   │   ├── metrics/
│   │   ├── middleware/
│   │   └── server/
│   ├── static/
│   └── templates/
```

## Разработка

Правила контрибьюции — см. [contribute.md](contribute.md).

## Как работает расчёт

Подробная документация алгоритма — см. [project-docs/how-calculation-works.md](project-docs/how-calculation-works.md).
