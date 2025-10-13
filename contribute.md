# Contributing to Tax Calculator

Спасибо, что решили внести вклад в проект **Tax Calculator** — gRPC-сервис на Go для расчёта прогрессивного подоходного налога в России начиная с 2025 года. Этот документ описывает, как локально поднять окружение, запустить тесты и внести изменения безопасно для продакшена.

---

## Требования

- Go 1.23.10.
- gRPC runtime:
    * google.golang.org/grpc v1.67.1
    * google.golang.org/protobuf v1.34.2
- (Для генерации кода) protoc или buf + плагины:
    * protoc ≥ 3.21 или buf (CLI).
    * protoc-gen-go (генератор protobuf-типов).
    * protoc-gen-go-grpc (генератор gRPC-стабов).
- (Опционально) инструменты разработчика:
    * grpcurl — ручная проверка RPC.
    * golangci-lint — линтер.
    * make — Makefile.

---

## Клонирование и базовая настройка

```bash
git clone <repo_url>
cd <repo_dir>
go mod download
```
---

## Установка зависимостей

```bash
make tidy
```
> Выполняет go mod tidy и проверяет, что go.mod и go.sum не содержат непроиндексированных изменений.

## Генерация gRPC-кода

```bash
make codegen
```
>Сгенерирует код из `api/tax.proto` в `gen/grpc/api`.

Требуемые инструменты:

```bash
# buf (CLI)
go install github.com/bufbuild/buf/cmd/buf@v1.45.0

# генераторы (совместимы с protobuf v1.34.2 и grpc v1.67.1)
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

# полезные утилиты (опционально)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@v1.8.9
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
```

---

## Сборка и запуск сервера

- Переменные окружения (.env)
    * LOG_MODE — json (по умолчанию) или text
    * LOG_LEVEL — debug|info|warn|error (по умолчанию info)
    * PORT — адрес прослушивания, например :50051

### Локально через Go

```bash
make run
```
> Запускает `cmd/main.go` с gRPC-сервером.

### Через Docker Compose

```bash
make docker-build     # собрать и запустить
make docker-up        # только запустить
make docker-down      # остановить
```

### gRPC-сервер доступен по адресу:
```
localhost:50051
```
Можно проверить с помощью `grpcurl` или теста:

```bash
grpcurl -plaintext localhost:50051 list
go test -run Test_Server_Healthz ./test
```


## Локальный запуск CI
Если установлен [`act`](https://github.com/nektos/act):

```bash
make local-CI
```

---

## 📁 Структура проекта

```
.ci                — Dockerfile (образ для ci)
.github/           — конфигурация CI (ci.yaml)
cmd/               — точка входа (main.go)
internal/          — server, calculator
api/               — (tax.proto, buf.yaml)
gen/               — сгенерированный gRPC-код
test/              — интеграционные тесты
pkg/               — логирование и утилиты
Dockerfile         — Dockerfile
Makefile           — сборка, тесты, генерация
```

---