# go-project-template

> [中文文档](docs/README.zh-CN.md)

A Go HTTP + gRPC service skeleton built with proto + buf + grpc-gateway.

## Stack

- **Go 1.26** — language
- **Proto3** — API definition
- **Buf** — protobuf build tool
- **gRPC** — RPC framework
- **grpc-gateway** — HTTP reverse proxy, one proto for both HTTP and gRPC
- **GORM + GORM Gen** — ORM & code generation
- **PostgreSQL** — database
- **OpenAPI** — auto-generated Swagger docs

## Getting Started

### Prerequisites

```bash
# protobuf toolchain
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# buf
# Windows: scoop install buf
# macOS:   brew install buf
# Linux:   go install github.com/bufbuild/buf/cmd/buf@latest

# gRPC debug tool (optional)
# Windows: scoop install grpcurl
# macOS:   brew install grpcurl
```

### Setup

```bash
cp config/app.template.yaml config/app.yaml   # first time only
```

### Generate & Run

```bash
# 1. Generate proto code + service registration
make pb

# 2. Generate database model & DAO (requires a running PostgreSQL)
make db

# 3. Build
make build

# 4. Run
make run

# gRPC listens on :9090
# HTTP listens on :8080
```

## Test

```bash
# HTTP
curl "http://localhost:8080/v1/hello?name=world"

# gRPC (requires grpcurl)
grpcurl -plaintext -d '{"name": "world"}' localhost:9090 hello.HelloService/SayHello
```

## Directory Layout

```
├── api/v1/               proto definitions
├── cmd/
│   ├── main.go           entrypoint, orchestrates startup
│   ├── router.go         server bootstrap (gRPC + HTTP gateway)
│   ├── router_gen.go     auto-generated service registration (DO NOT EDIT)
│   ├── gormgen/          GORM Gen code generator
│   └── register-gen/     service registration code generator
├── config/               runtime config
├── internal/
│   ├── api/v1/           generated proto code (DO NOT EDIT)
│   ├── config/           config loading
│   ├── consts/           global constants, enums (placeholder)
│   ├── dto/              data transfer objects (placeholder)
│   ├── middleware/        gRPC interceptors / HTTP middleware (placeholder)
│   ├── model/            GORM Gen data models + custom query interfaces
│   ├── mq/               message queue (placeholder)
│   ├── pkg/              shared utilities (i18n, etc.)
│   ├── repository/       data access layer (DB + Cache)
│   └── service/          business logic
├── locales/              i18n / custom error code messages
└── scripts/              SQL, shell scripts (placeholder)
```

## Adding a New Service

1. Create a proto file under `api/v1/`
2. Define the service, messages, and `google.api.http` annotations
3. Run `make pb` — code generation and service registration are automatic
4. Implement the business logic in `internal/service/`

> Service registration is handled by `cmd/register-gen` and runs automatically as part of `make pb`.

## Proto Conventions

- Use `google.api.http` annotations to define HTTP routes
- Route format: `/v1/{service-name}/{action}`, multi-word segments joined with `-`
  (e.g. `/v1/user/model/list`, `/v1/hello/filter-models`)
- All proto files go under `api/v1/`, grouped by domain
- Generated code uses `source_relative` path style

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make pb` | Generate proto code + auto-register services |
| `make db` | Generate GORM model & DAO from database |
| `make build` | Tidy deps, format, lint, then build |
| `make run` | Tidy deps, then run the service |
| `make test` | Run all tests |
| `make fmt` | Format Go code |
| `make lint` | Run `go vet` |
| `make clean` | Remove build artifacts |