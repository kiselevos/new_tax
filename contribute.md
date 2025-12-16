# Руководство по контрибьюции

Спасибо за интерес к проекту **Tax Calculator** - сервису для расчёта подоходного налога (НДФЛ) в России по правилам 2025 года.  

Проект реализован на **Go + ConnectRPC**, использует **proto-first подход** и состоит из backend-сервиса, frontend (BFF) и отдельного инфраструктурного слоя для observability.
---

## Архитектура проекта

Проект организован как монорепозиторий и включает несколько логических слоёв.

### Backend
- Основной сервис расчёта НДФЛ
- Реализован на **Go + ConnectRPC**
- Содержит бизнес-логику, валидацию и API
- Использует proto-контракты как единый источник истины
- Основные директории: `cmd/`, `internal/`, `pkg/`, `gen/`

### Frontend (BFF)
- Веб-приложение на **Go templates**
- Отдельный Go-модуль со своим `go.mod`
- Работает как thin client к backend через ConnectRPC
- Не содержит бизнес-логики
- Основные директории: `web/cmd/`, `web/internal/`, `web/templates/`, `web/static/`

### Observability & Infrastructure
- Вынесена в отдельную директорию `infra/`
- Используется для мониторинга и анализа логов
- Не влияет на бизнес-логику приложения
---

## ⚙️ Требования к окружению

| Компонент | Версия |
|------------|---------|
| Go | **1.23.10** |
| protoc | **3.21.12** |
| protoc-gen-go | **v1.34.2** |
| protoc-gen-connect-go | **v1.16.0** |
| make | **4.4.1** |
| docker | **27.1.2** |
| docker-compose | **2.29.2** |

_Опционально:_  
`grpcurl 1.9.1`, `golangci-lint 1.60.3` - для тестирования и линтинга.

---

## 🚀 Быстрый старт

1. **Клонируйте репозиторий:**
   ```bash
   git clone https://github.com/kiselevos/new_tax
   cd new_tax
   ```

2. **Установите зависимости:**
   ```bash
   make setup-tools # (опционально для линтера, тестов и кодгена)
   
   make setup
   ```

3. **Сгенерируйте код (если нужно):**
   ```bash
   make codegen
   ```

4. **Запустите проект:**
   ```bash
   make run-all
   ```
   или через Docker:
   ```bash
   make docker-build
   ```

5. **Observability (опционально)**
Для мониторинга логов и метрик используется отдельный инфраструктурный стек.
Запуск observability-стека:
  ```bash
  make docker-infra-up
  ```
---

## 🧰 Основные команды Make

| Команда | Назначение |
|----------|-------------|
| `make setup` | Установка зависимостей backend и frontend |
| `make codegen` | Генерация gRPC / ConnectRPC кода из `.proto` |
| `make run` | Запуск backend локально |
| `make run-all` | Запуск backend и frontend одновременно |
| `make test-all` | Запуск тестов |
| `make lint-all` | Проверка линтером |
| `make docker-build` | Сборка и запуск контейнеров |
| `make docker-down` | Остановка контейнеров |

---

## 🧪 Тестирование и проверка

- **Линтинг:**
  ```bash
  make lint-all
  ```
- **Юнит- и интеграционные тесты:**
  ```bash
  make test-all
  ```
- **Проверка форматирования:**
  ```bash
  make check-fmt
  ```
- **Проверка генерации кода:**
  ```bash
  make check-generated
  ```

---

## 🧱 Структура проекта

```bash
├── cmd/            # Точка входа backend
├── internal/       # Бизнес-логика и сервер
├── pkg/            # Общие утилиты и логгер
├── gen/            # Сгенерированный gRPC-код
├── web/            # Frontend (BFF)
├── infra/          # Prometheus / Loki / Grafana
├── test/           # Интеграционные тесты
├── docs/grpc/      # Proto-контракты
├── Makefile        # Команды разработки
└── Dockerfile      # Сборка backend
```
