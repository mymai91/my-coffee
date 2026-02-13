# My Coffee Shop ☕

A production-ready Go microservices project demonstrating gRPC, Connect RPC, GORM, and modern API patterns.

> 📚 **New to this project?** Check out the comprehensive [Go Learning Guide](./GO_LEARNING_GUIDE.md) with flow diagrams, architecture explanations, and interview prep!

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

### Create REST with Connectrpc

Key details
apisvc is untouched — its gRPC clients still work because Connect RPC servers with h2c support the gRPC protocol natively
Validation — uses connectrpc.com/validate interceptor instead of the gRPC middleware
Protocol support — all three protocols (Connect, gRPC, gRPC-Web) are supported out of the box


#### 1️⃣ REST vs gRPC + Connect RPC

| Point                | REST               | gRPC (Connect RPC)          |
| -------------------- | ------------------ | --------------------------- |
| Protocol             | HTTP/1.1           | HTTP/2                      |
| Data format          | JSON               | Protobuf                    |
| Typing               | Weakly typed       | Strongly typed              |
| Frontend browser     | Easy               | Need gRPC-Web / Connect RPC |
| Streaming            | No                 | Yes (bi-directional)        |
| Performance          | Slower             | Faster, more efficient      |
| Generate client code | Depends on library | Auto from .proto            |

#### 2️⃣ Connect RPC

* Connect RPC is a Go library built on gRPC.

* When you write a .proto file once, Connect RPC can automatically generate:

    - gRPC server & client → for backend-to-backend communication.

    - REST endpoints → for frontend to call using JSON.

* This means the frontend does not need to know gRPC. It can just use REST JSON.

* Backend services still talk to each other with gRPC, so internal performance is fast.

### Backend — Added CORS to both services:

main.go — CORS middleware with Connect-Protocol-Version header
main.go — Same CORS middleware

Frontend — Connect RPC via plain fetch():

api.ts — Calls Connect RPC endpoints directly:
fetchMenu() → POST /menu.MenuService/GetMenu
fetchOrders() → POST /brew.BrewService/ListOrders
createOrder() → POST /brew.BrewService/OrderDrink
getOrder() → POST /brew.BrewService/GetOrder
updateOrderStatus() → POST /brew.BrewService/UpdateOrderStatus
deleteOrder() → POST /brew.BrewService/DeleteOrder

vite.config.ts — Removed proxy (no longer needed)

How to run now (only 2 backend services!)

```
make run-brewsvc     # Terminal 1 - port 50051
make run-menusvc     # Terminal 2 - port 50052
cd web && npm run dev # Terminal 3 - port 5173
```

## Hexagon architecture

branch feature-9-update-with-hexagonal-architecture

The migration to hexagonal architecture is complete and compiles successfully. Here's a summary of the new structure:

New hexagonal layout

```
internal/
  core/                          ← THE HEXAGON (no external dependencies)
    domain/
      order.go                   ← Order entity + OrderStatus (pure Go)
      menu_item.go               ← MenuItem entity (pure Go)
    ports/
      repositories.go            ← OrderRepository interface (driven port)
      services.go                ← OrderService + MenuService interfaces (driving ports)
    services/
      order_service.go           ← Use-case logic, depends only on ports
      menu_service.go            ← Menu use-case logic
  adapters/                      ← OUTSIDE THE HEXAGON
    repository/
      order_model.go             ← GORM model + domain ↔ model mappers
      gorm_order_repository.go   ← Implements ports.OrderRepository with GORM
    handler/
      brew_handler.go            ← Implements brewconnect.BrewServiceHandler, calls ports.OrderService
      menu_handler.go            ← Implements menuconnect.MenuServiceHandler, calls ports.MenuService
```

#### Key design decisions

## Architecture Dependency Rules

| Layer | Depends On | Must NOT Depend On |
|-------|------------|-------------------|
| **Domain** | Nothing | Any other layer |
| **Ports (Interfaces)** | Domain | Adapters, Infrastructure |
| **Services (Use Cases)** | Domain, Ports | Adapters, GORM, Protobuf |
| **Handler Adapter** | Ports, Protobuf | GORM, Repository Implementation |
| **Repository Adapter** | Ports, GORM | Handlers, Protobuf |

Wiring (in cmd)

Each main.go wires the layers with plain constructor injection:

```
Repository (driven adapter) → Service (core) → Handler (driving adapter)
```

The old packages under brews, menus, models, and repository are no longer imported by any cmd entrypoint and can be deleted when you're ready.

### (old version) Building REST Gateway for Your Coffee Shop

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


## Frontend Vite & React Query


main.tsx

App.tsx

Menu.tsx

Orders.tsx

OrderForm.tsx

What React Query gives you now

| Before (manual)                        | After (React Query)                                        |
|----------------------------------------|------------------------------------------------------------|
| `useState` + `useEffect` + `setLoading` | `useQuery` returns `{ data, isLoading, error }` automatically |
| `ordersKey` ref hack to force refetch   | `invalidateQueries` on mutation success — the right way   |
| Manual refresh button re-runs full fetch | `refetch()` from the query + shows `isFetching` spinner  |
| No background updates                   | Orders auto-refetch every 10s so you see status changes live |
| No window focus refetch                 | Refetches on tab focus — switch back and see latest orders |
| Menu fetched every mount                | Menu cached for 5 min (`staleTime`) — instant on tab switch |


### File structure

* hooks.ts — useMenu, useOrders, useCreateOrder — all query logic in one place
* api.ts — unchanged, pure fetch functions (clean separation)
* main.tsx — QueryClientProvider wraps the app
* Components are now much simpler — no manual loading/error state management