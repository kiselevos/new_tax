# Russian Personal Income Tax Calculator (NDFL)

**English** | [Русский](README.ru.md)

A service that calculates Russian personal income tax (NDFL) and employer insurance contributions month by month. Built for employees, HR specialists and developers who need an accurate, programmable calculation.

## Try it

[calculator-ndfl.ru](https://calculator-ndfl.ru)

## What it calculates

- **Progressive NDFL scale** (2025+): five brackets — 13%, 15%, 18%, 20%, 22%
- **Three employment types**: employment contract (TD), civil-law contract (GPH), self-employed (NPD, 4% / 6%)
- **Tax deductions**: standard child deductions (Art. 218 of the Tax Code), social deductions — medical treatment and education (Art. 219), property deductions — housing purchase and mortgage interest (Art. 220)
- **Regional coefficient (RK)** — part of the main income, taxed on the general progressive scale
- **Northern allowance (SN)** — a separate tax base, taxed on the simplified 13% / 15% scale
- **Employer contributions**: pension (22% up to the annual cap, 10% above), medical insurance (5.1%), social insurance (2.9% up to the cap)
- **One-off bonuses** by month, with correct bracket transitions
- **Law-enforcement employees** — simplified 13% / 15% scale
- **Non-residents** — flat 30% rate
- **Monthly breakdown with a year-to-date (YTD) running total** for the whole tax period

All money values are stored as `uint64` in kopecks (1/100 of a ruble) — no floating point anywhere in the calculation.
Tax is rounded to whole rubles per Art. 52 §6 of the Tax Code (less than 50 kopecks is dropped, 50 kopecks or more is rounded up).

## Quick start

```bash
# Create a local config
cp .env.example .env

# Build and start all containers
docker compose up --build

# or via Makefile
make docker-build
```

After startup:
- Web UI and REST API: http://localhost:8080
- Backend gRPC server: `localhost:50051` (plaintext gRPC, reflection enabled)

Without Docker:

```bash
go run ./cmd/main.go &
cd web && go run ./cmd/web.go

# or via Makefile
make run-all
```

Health check:

```bash
grpcurl -plaintext localhost:50051 tax.TaxService/Healthz
```

gRPC call example (the private method requires the API key you set as `API_KEY` in `.env`):

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

REST (JSON) example via the web service:

```bash
curl -X POST http://localhost:8080/api/v1/calc \
  -H 'Content-Type: application/json' \
  -d '{"gross_salary": 20000000, "territorial_multiplier": 120}'
```

OpenAPI documentation (Swagger UI) is available at `/api-docs`.

### Public demo key

This is an open sandbox project, so the private API key is intentionally public — it is the `API_KEY` from `.env.example` and it works in production too:

```bash
curl -X POST https://calculator-ndfl.ru/api/v1/private-calc \
  -H 'Content-Type: application/json' \
  -H 'x-api-key: api_7f28c3a4b1ef49d9a6c1d742e91f35c2' \
  -d '{"gross_salary": 20000000, "employment_type": "TD", "start_date": "2026-01-01"}'
```

Rate limits apply. Feel free to use it and try things out.

Main Make targets:

```bash
make setup          # install dependencies (backend + web)
make codegen        # generate Go gRPC code from .proto
make run-all        # run backend + web locally
make docker-build   # build and start containers
make test-all       # run backend tests
make ci             # run lint + tests before pushing
```

## Architecture

```
Browser ──HTTP──▶ web (BFF, :8080) ──gRPC──▶ backend (:50051)
                   │  HTML pages                 │  business logic
                   │  REST API /api/v1/*         │  auth, rate limiting
                   │  /metrics                   │
                   ▼                             ▼
             Prometheus ◀── scrape         stdout JSON logs ──▶ Promtail ──▶ Loki ──▶ Grafana
```

The project consists of two independent Go applications plus an observability stack.

### Backend
- **Go 1.23**, gRPC server built on `google.golang.org/grpc`
- **Proto-first**: the contract lives in `docs/grpc/tax.proto`, Go code is generated into `gen/`
- Owns all business logic: tax calculation, input validation, API
- Interceptor chain: panic recovery → request logging (request ID) → API-key auth → rate limiting
- Structured logging via `log/slog` (thin wrapper in `pkg/logx`)
- Unit and integration tests with **testify**

### Web (BFF)
- Separate Go module with its own `go.mod`
- Server-side rendered UI with Go `html/template`, vanilla CSS and JavaScript
- Acts as a Backend-for-Frontend: talks to the backend through a gRPC client
- Contains no tax logic — all calculations happen in the backend
- Public JSON REST API (`/api/v1/calc`, `/api/v1/private-calc`) documented with OpenAPI
- Exposes Prometheus metrics (`/metrics`) and writes structured JSON logs
- GeoIP lookup (CSV database) to label metrics by region

### API: public and private methods
- `CalculatePublic` — a simplified calculation, open to everyone
- `CalculatePrivate` — the full calculation (start date, bonuses, deductions, employment type, employer contributions); requires an `x-api-key` header, checked in a backend interceptor. Requests from the web UI are marked as internal traffic
- Separate rate limits for public and private calls (token bucket, `golang.org/x/time/rate`)
- All methods are defined in `.proto` (`tax.TaxService`)

### Observability
- Backend and web log to stdout in JSON
- Promtail collects container logs and ships them to Loki
- Prometheus scrapes the web service metrics
- Grafana is the single place to view logs, metrics and business dashboards

Start the observability stack:

```bash
docker compose -f infra/docker-compose.yaml up -d

# or via Makefile
make docker-infra-up
```

After startup:
- Grafana: http://localhost:3000
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100

### Project structure

```
├── README.md / README.ru.md
├── ABOUT.md                 # Product description (in Russian)
├── contribute.md            # Contribution guide (in Russian)
├── Dockerfile               # Backend image
├── docker-compose.yaml
├── Makefile
│
├── cmd/
│   └── main.go              # Backend entry point (gRPC server)
│
├── docs/
│   └── grpc/
│       └── tax.proto        # Proto contract (source of truth for the API)
│
├── gen/                     # Code generated from .proto (separate Go module)
│
├── internal/                # Backend internals
│   ├── calculate/           # Tax calculation domain logic
│   ├── config/              # Configuration from env
│   ├── middleware/          # gRPC interceptors: auth, logging, recovery, rate limit
│   └── server/              # gRPC server and service implementation
│
├── pkg/
│   └── logx/                # slog-based structured logger
│
├── test/
│   └── server_integ_test.go # Integration tests of the gRPC server
│
├── infra/                   # Observability stack
│   ├── docker-compose.yaml
│   ├── prometheus/
│   ├── loki/
│   ├── promtail/
│   └── grafana/
│
├── web/                     # BFF (separate Go module)
│   ├── Dockerfile
│   ├── Makefile
│   ├── cmd/web.go           # Web entry point
│   ├── handlers/            # HTML page handlers
│   ├── internal/
│   │   ├── api/             # REST API handlers and DTOs
│   │   ├── client/          # gRPC client to the backend
│   │   ├── config/
│   │   ├── geoip/
│   │   ├── metrics/
│   │   ├── middleware/
│   │   └── server/
│   ├── static/              # CSS, JS, OpenAPI spec
│   └── templates/           # Go HTML templates
│
└── project-docs/            # Calculation docs, tax constants, roadmap (in Russian)
```

## Development

See [contribute.md](contribute.md) for contribution guidelines.

## How the calculation works

A detailed description of the algorithm (in Russian): [project-docs/how-calculation-works.md](project-docs/how-calculation-works.md).
