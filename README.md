# My Coffee Shop ☕

A Go microservices project for learning gRPC.

## Architecture

```
┌─────────────┐      gRPC       ┌─────────────┐
│  coffeecli  │ ◄──────────────► │  menusvc    │ (port 50052)
│  (Customer) │                  └─────────────┘
│             │      gRPC       ┌─────────────┐
│             │ ◄──────────────► │  brewsvc    │ (port 50051)
└─────────────┘                  └─────────────┘
```

## Quick Start

### 1. Install tools

```bash
make install-tools
```

### 2. Generate proto code

```bash
make proto
```

### 3. Build & Test

```bash
make build
make test
```

### 4. Run services

```bash
make start-services
```

## Project Structure

Option 1: Monorepo (What You're Using)

```
my-coffee/
├── go.mod                    # One go.mod for all services
├── proto/
│   └── menu.proto
├── cmd/
│   ├── menusvc/
│   │   └── main.go          # Menu service entry point
│   ├── brewsvc/
│   │   └── main.go          # Brew service entry point
│   ├── coffeecli/
│   │   └── main.go          # Customer CLI entry point
│   └── baristacli/
│       └── main.go          # Barista CLI entry point
└── Makefile
```


```
my-coffee/
├── cmd/                    # Application entry points
│   ├── menusvc/           # Menu gRPC service
│   └── brewsvc/           # Brew gRPC service
├── internal/              # Private application code
│   ├── menus/             # Menu business logic
│   └── brews/             # Brew business logic
├── proto/                 # Protocol buffer definitions
│   ├── menu/
│   └── brew/
├── gen/                   # Generated code (auto-created)
└── scripts/               # Helper scripts
```

### How to run

1. Start the gRPC services as usual (brew on :50051, menu on :50052)

```

make run-menusvc

make run-brewsvc

```

2. Start the API gateway:

```
go run ./cmd/apisvc
```

3. Start the Vite dev server:

```
cd web && npm run dev
```

4. Open http://localhost:5173 in your browser
The Vite dev server is configured to proxy /api requests to the Go API gateway on :9000.


### Building REST Gateway for Your Coffee Shop

Architecture Overview

```
┌─────────────────────────────────────────────────┐
│              Vite Frontend                      │
│            (localhost:5173)                     │
└───────────────────┬─────────────────────────────┘
                    │
                    │ HTTP/REST + JSON
                    │
┌───────────────────▼─────────────────────────────┐
│           REST API Gateway                      │
│            (localhost:8080)                     │
│                                                 │
│  Routes:                                        │
│  GET    /api/menu          → ListMenuItems     │
│  GET    /api/menu/:id      → GetMenuItem       │
│  POST   /api/menu          → CreateMenuItem    │
│  PUT    /api/menu/:id      → UpdateMenuItem    │
│  DELETE /api/menu/:id      → DeleteMenuItem    │
│                                                 │
│  GET    /api/orders        → ListOrders        │
│  POST   /api/orders        → CreateOrder       │
└─────────────┬───────────────┬───────────────────┘
              │               │
              │ gRPC          │ gRPC
              │               │
    ┌─────────▼────────┐  ┌──▼──────────────┐
    │   MenuService    │  │   BrewService   │
    │  (localhost:50051)│  │ (localhost:50052)│
    │                  │  │                 │
    │  Internal gRPC   │  │  Internal gRPC  │
    └──────────────────┘  └─────────────────┘

```

