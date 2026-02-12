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