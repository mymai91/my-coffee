# Coffee Shop - System Architecture

## 🏗️ System Overview

```
┌────────────────────────────────────────────────────────────────────┐
│                        FRONTEND LAYER                              │
│  ┌──────────────────────────┐     ┌─────────────────────────┐     │
│  │   React + Vite + TS      │     │   CLI Clients           │     │
│  │   localhost:5173         │     │   (coffeecli)           │     │
│  │                          │     │                         │     │
│  │  • React Query           │     │  • gRPC Client          │     │
│  │  • TypeScript            │     │  • Menu browsing        │     │
│  │  • Order management UI   │     │  • Order status check   │     │
│  └────────┬─────────────────┘     └───────────┬─────────────┘     │
│           │                                   │                   │
└───────────┼───────────────────────────────────┼───────────────────┘
            │                                   │
            │ HTTP/REST (JSON)                  │ gRPC (Protobuf)
            │                                   │
┌───────────▼───────────────────────────────────▼───────────────────┐
│                     API GATEWAY LAYER                             │
│  ┌────────────────────────────────────────────────────────────┐   │
│  │              apisvc (localhost:9000)                       │   │
│  │                                                            │   │
│  │  REST Endpoints:                                           │   │
│  │  • GET    /api/menu              → MenuService.GetMenu    │   │
│  │  • GET    /api/orders            → BrewService.ListOrders │   │
│  │  • POST   /api/orders            → BrewService.OrderDrink │   │
│  │  • GET    /api/orders/:id        → BrewService.GetOrder   │   │
│  │  • PATCH  /api/orders/:id/status → UpdateOrderStatus      │   │
│  │  • DELETE /api/orders/:id        → BrewService.DeleteOrder│   │
│  │                                                            │   │
│  │  Features:                                                 │   │
│  │  • CORS middleware                                         │   │
│  │  • JSON ↔ Protobuf conversion                             │   │
│  │  • gRPC client connections                                │   │
│  └───────────┬────────────────────────────┬───────────────────┘   │
└─────────────┼────────────────────────────┼───────────────────────┘
              │                            │
              │ gRPC                       │ gRPC
              │ (Protobuf)                 │ (Protobuf)
              │                            │
┌─────────────▼──────────────┐  ┌──────────▼────────────────────────┐
│  MENU SERVICE              │  │  BREW SERVICE                     │
│  menusvc (port 50052)      │  │  brewsvc (port 50051)             │
│                            │  │                                   │
│  Connect RPC Server        │  │  Connect RPC Server               │
│  ┌──────────────────────┐  │  │  ┌────────────────────────────┐   │
│  │ Protocols:           │  │  │  │ Protocols:                 │   │
│  │ • gRPC               │  │  │  │ • gRPC                     │   │
│  │ • Connect            │  │  │  │ • Connect                  │   │
│  │ • gRPC-Web           │  │  │  │ • gRPC-Web                 │   │
│  └──────────────────────┘  │  │  └────────────────────────────┘   │
│                            │  │                                   │
│  Handler:                  │  │  Handler:                         │
│  • GetMenu()               │  │  • OrderDrink()                   │
│    Returns static menu     │  │  • ListOrders()                   │
│    items (Espresso, Latte) │  │  • GetOrder()                     │
│                            │  │  • UpdateOrderStatus()            │
│                            │  │  • DeleteOrder()                  │
│                            │  │                                   │
│  No database needed        │  │  Repository Layer:                │
│                            │  │  ┌────────────────────────────┐   │
│                            │  │  │ OrderRepository            │   │
│                            │  │  │ • Create()                 │   │
│                            │  │  │ • FindAll()                │   │
│                            │  │  │ • FindByID()               │   │
│                            │  │  │ • Update()                 │   │
│                            │  │  │ • Delete()                 │   │
│                            │  │  └───────────┬────────────────┘   │
│                            │  │              │                    │
└────────────────────────────┘  └──────────────┼────────────────────┘
                                               │
                                               │ GORM
                                               │
                                   ┌───────────▼───────────────────┐
                                   │  PERSISTENCE LAYER            │
                                   │  PostgreSQL (port 5454)       │
                                   │                               │
                                   │  Table: orders                │
                                   │  ┌─────────────────────────┐  │
                                   │  │ id (PK, auto)           │  │
                                   │  │ menu_item_name          │  │
                                   │  │ status (default QUEUED) │  │
                                   │  │ created_at              │  │
                                   │  │ updated_at              │  │
                                   │  └─────────────────────────┘  │
                                   └───────────────────────────────┘
```

## 📊 Data Flow Examples

### Example 1: Customer Orders a Latte

```
┌──────────┐                                                    ┌──────────┐
│ Customer │                                                    │ Database │
│ (React)  │                                                    │ (Postgres│
└─────┬────┘                                                    └────┬─────┘
      │                                                              │
      │ 1. Click "Order Latte"                                      │
      │    POST /api/orders                                         │
      │    {"menuItemName": "Latte"}                                │
      ├─────────────────────────────────────────────┐               │
      │                                             │               │
      │                                    ┌────────▼─────────┐     │
      │                                    │    apisvc        │     │
      │                                    │                  │     │
      │                                    │ 2. Parse JSON    │     │
      │                                    │ 3. Create gRPC   │     │
      │                                    │    Request       │     │
      │                                    └────────┬─────────┘     │
      │                                             │               │
      │                                             │ gRPC          │
      │                                    ┌────────▼─────────┐     │
      │                                    │   brewsvc        │     │
      │                                    │                  │     │
      │                                    │ 4. Validate      │     │
      │                                    │    (protovalidate)    │
      │                                    │ 5. Create Order  │     │
      │                                    │    model         │     │
      │                                    └────────┬─────────┘     │
      │                                             │               │
      │                                             │ GORM          │
      │                                             ├───────────────▶
      │                                             │ 6. INSERT INTO │
      │                                             │    orders      │
      │                                             │ 7. Return ID=42│
      │                                             │◀───────────────┤
      │                                    ┌────────▼─────────┐     │
      │                                    │   brewsvc        │     │
      │                                    │ 8. Return        │     │
      │                                    │    order-42      │     │
      │                                    └────────┬─────────┘     │
      │                                             │               │
      │                                             │ gRPC          │
      │                                    ┌────────▼─────────┐     │
      │                                    │    apisvc        │     │
      │                                    │ 9. Convert to    │     │
      │                                    │    JSON          │     │
      │                                    └────────┬─────────┘     │
      │                                             │               │
      │ 10. 200 OK                                  │               │
      │     {"orderId": "order-42"}                 │               │
      │◀────────────────────────────────────────────┘               │
      │                                                              │
      │ 11. React Query updates cache                               │
      │     & triggers re-render                                    │
      │                                                              │
```

### Example 2: Status Update Flow (Barista makes coffee)

```
Barista                 brewsvc                  Database
  │                        │                         │
  │ PATCH /api/orders/     │                         │
  │ order-42/status        │                         │
  │ {"status":"BREWING"}   │                         │
  ├───────────────────────▶│                         │
  │                        │                         │
  │                        │ 1. Parse "order-42"     │
  │                        │    Extract ID: 42       │
  │                        │                         │
  │                        │ 2. FindByID(42)         │
  │                        ├────────────────────────▶│
  │                        │                         │
  │                        │ 3. SELECT * FROM orders │
  │                        │    WHERE id = 42        │
  │                        │◀────────────────────────┤
  │                        │                         │
  │                        │ 4. Update Status        │
  │                        │    QUEUED -> BREWING    │
  │                        │                         │
  │                        │ 5. Save()               │
  │                        ├────────────────────────▶│
  │                        │                         │
  │                        │ 6. UPDATE orders        │
  │                        │    SET status='BREWING',│
  │                        │    updated_at=NOW()     │
  │                        │◀────────────────────────┤
  │                        │                         │
  │ 200 OK                 │                         │
  │ {order with BREWING}   │                         │
  │◀───────────────────────┤                         │
  │                        │                         │
```

## 🔐 Protocol Support Matrix

### Connect RPC Multi-Protocol Support

```
┌────────────────────────────────────────────────────────────────┐
│           brewsvc / menusvc (Connect RPC Server)               │
│                                                                │
│  Single codebase, handles 3 protocols automatically:          │
└────────────┬───────────────────┬──────────────────────────────┘
             │                   │                   │
    ┌────────▼────────┐  ┌───────▼────────┐  ┌──────▼───────────┐
    │ gRPC Protocol   │  │ Connect Proto  │  │ gRPC-Web         │
    └────────┬────────┘  └───────┬────────┘  └──────┬───────────┘
             │                   │                   │
┌────────────▼────────┐  ┌───────▼──────────┐ ┌─────▼────────────┐
│ FROM: apisvc        │  │ FROM: curl       │ │ FROM: Browser JS │
│ • HTTP/2            │  │ • HTTP/1.1 or 2  │ │ • HTTP/1.1       │
│ • Binary Protobuf   │  │ • JSON or Proto  │ │ • Binary Protobuf│
│ • Content-Type:     │  │ • Content-Type:  │ │ • Content-Type:  │
│   application/grpc  │  │   connect+proto  │ │   grpc-web+proto │
└─────────────────────┘  └──────────────────┘ └──────────────────┘

Example requests:

1. gRPC (from apisvc):
   grpcClient.OrderDrink(ctx, &OrderRequest{...})

2. Connect (from curl):
   curl -X POST http://localhost:50051/brew.BrewService/OrderDrink \
     -H "Content-Type: application/json" \
     -d '{"menuItemName":"Latte"}'

3. gRPC-Web (from browser):
   const client = createPromiseClient(BrewService, transport);
   await client.orderDrink({menuItemName: "Latte"});
```

## 🗂️ Code Organization

### Project Structure

```
my-coffee/
│
├── cmd/                          # Executables (main packages)
│   ├── apisvc/main.go           # REST API Gateway (port 9000)
│   ├── brewsvc/main.go          # Brew gRPC service (port 50051)
│   ├── menusvc/main.go          # Menu gRPC service (port 50052)
│   ├── coffeecli/main.go        # Customer CLI client
│   └── migrate/main.go          # Database migration tool
│
├── internal/                     # Private application code
│   ├── brews/                   # Brew service logic
│   │   └── brew.go              # Handler implementation
│   ├── menus/                   # Menu service logic
│   │   └── menu.go              # Handler implementation
│   ├── models/                  # Database models
│   │   └── order.go             # Order struct + GORM config
│   ├── repository/              # Data access layer
│   │   └── order_repository.go  # CRUD operations
│   └── datbase/                 # Database connection
│       └── connection.go        # GORM setup
│
├── proto/                        # Protocol Buffer definitions
│   ├── brew/
│   │   └── brew.proto           # Brew service API
│   └── menu/
│       └── menu.proto           # Menu service API
│
├── gen/                          # Generated code (DO NOT EDIT)
│   └── proto/
│       ├── brew/
│       │   ├── brew.pb.go       # Protobuf types
│       │   └── brewconnect/
│       │       └── brew.connect.go  # Connect RPC handlers
│       └── menu/
│           ├── menu.pb.go
│           └── menuconnect/
│               └── menu.connect.go
│
├── migrations/                   # Database migrations
│   ├── 000001_create_orders_table.up.sql
│   └── 000001_create_orders_table.down.sql
│
├── web/                          # Frontend React app
│   ├── src/
│   │   ├── api.ts               # API client functions
│   │   ├── hooks.ts             # React Query hooks
│   │   ├── App.tsx              # Main component
│   │   └── components/
│   │       ├── Menu.tsx
│   │       ├── Orders.tsx
│   │       └── OrderForm.tsx
│   └── package.json
│
├── config/                       # Configuration management
│   └── config.go                # Env var loading
│
├── .env                          # Environment variables (git ignored)
├── go.mod                        # Go dependencies
├── buf.yaml                      # Buf configuration
├── buf.gen.yaml                  # Code generation config
├── Makefile                      # Build automation
└── README.md                     # Quick start guide
```

## 🧩 Component Responsibilities

### 1. **apisvc** - API Gateway
**Responsibility:** Translate between HTTP/REST and gRPC
```
Input:  HTTP requests from frontend (JSON)
Output: HTTP responses (JSON)
Does:   
  • Route HTTP requests to gRPC services
  • Handle CORS
  • Convert JSON ↔ Protobuf
  • Map REST endpoints to gRPC methods
Doesn't:
  • Business logic
  • Database access
  • Data validation (handled by gRPC services)
```

### 2. **brewsvc** - Brew Service
**Responsibility:** Order management and business logic
```
Input:  gRPC requests (Protobuf)
Output: gRPC responses (Protobuf)
Does:
  • Validate requests (protovalidate)
  • Create/read/update/delete orders
  • Enforce business rules
  • Call repository for data access
Doesn't:
  • Direct SQL queries (uses repository)
  • HTTP handling
  • JSON parsing
```

### 3. **menusvc** - Menu Service
**Responsibility:** Provide coffee menu
```
Input:  gRPC requests (Protobuf)
Output: gRPC responses (Protobuf)
Does:
  • Return static menu items
  • Could be extended to database-backed menu
Doesn't:
  • Order management
  • Database access (currently)
```

### 4. **OrderRepository** - Data Access Layer
**Responsibility:** Database operations
```
Input:  Model structs (models.Order)
Output: Model structs or errors
Does:
  • CRUD operations
  • SQL query execution (via GORM)
  • Data persistence
Doesn't:
  • Business logic
  • Validation
  • Error code mapping
```

## 🔄 Technologies & Their Roles

| Technology | Purpose | Why Used |
|-----------|---------|----------|
| **Go** | Backend language | Fast, simple, great for microservices |
| **gRPC** | RPC framework | Efficient binary protocol, type-safe |
| **Connect RPC** | Modern RPC | Multi-protocol support (gRPC + REST) |
| **Protocol Buffers** | Serialization | Strongly typed, language agnostic |
| **GORM** | ORM | Type-safe database access, migrations |
| **PostgreSQL** | Database | Reliable, feature-rich relational DB |
| **Buf** | Protobuf tooling | Better than protoc, linting, breaking change detection |
| **golang-migrate** | Migrations | Version control for database schema |
| **React** | Frontend framework | Component-based UI |
| **React Query** | State management | Caching, auto-refetch, loading states |
| **Vite** | Build tool | Fast dev server, optimized builds |
| **TypeScript** | Type safety | Catch errors at compile time |

## 🚀 Deployment View

### Development
```
┌─────────────────────────────────────────────────┐
│              Developer Laptop                   │
│                                                 │
│  Terminal 1: make run-brewsvc   (port 50051)  │
│  Terminal 2: make run-menusvc   (port 50052)  │
│  Terminal 3: make run-apisvc    (port 9000)   │
│  Terminal 4: cd web && npm run dev (port 5173)│
│                                                 │
│  Docker: PostgreSQL (port 5454)                │
└─────────────────────────────────────────────────┘
```

### Production (Future)
```
┌──────────────────────────────────────────────────────┐
│                  Cloud Provider                      │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Kubernetes Cluster                            │ │
│  │                                                │ │
│  │  Pods:                                         │ │
│  │  • apisvc (replicas: 3)                       │ │
│  │  • brewsvc (replicas: 2)                      │ │
│  │  • menusvc (replicas: 2)                      │ │
│  │                                                │ │
│  │  Services:                                     │ │
│  │  • Load Balancer → apisvc                     │ │
│  │  • ClusterIP → brewsvc                        │ │
│  │  • ClusterIP → menusvc                        │ │
│  └────────────────────────────────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Managed PostgreSQL                            │ │
│  │  (e.g., AWS RDS, Cloud SQL)                    │ │
│  └────────────────────────────────────────────────┘ │
│                                                      │
│  ┌────────────────────────────────────────────────┐ │
│  │  Frontend (Static hosting)                     │ │
│  │  (e.g., Netlify, Vercel, S3 + CloudFront)     │ │
│  └────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

## 📈 Scalability Considerations

### Current State (Single Instance)
```
1 request → 1 service instance → 1 database
```

### Future State (Horizontal Scaling)
```
Load Balancer
     │
     ├─→ apisvc-1 ─┐
     ├─→ apisvc-2 ─┼─→ brewsvc-1 ─┐
     └─→ apisvc-3 ─┘   brewsvc-2 ─┼─→ Database (with connection pooling)
                                  ─┘
```

### What Makes This Scalable:
- ✅ Stateless services (can add more instances)
- ✅ Database connection pooling (GORM handles this)
- ✅ gRPC for efficient inter-service communication
- ✅ Proper error handling
- ✅ Repository pattern (easy to add caching layer)

---

## 🔍 Key Design Decisions

### 1. Why Connect RPC instead of plain gRPC?
- **Browser compatibility**: Can call from React without gRPC-Web proxy
- **Debuggability**: Can test with curl using JSON
- **Flexibility**: Supports 3 protocols from same code
- **Future-proof**: Modern RPC framework

### 2. Why separate apisvc from brewsvc/menusvc?
- **Separation of concerns**: REST translation separate from business logic
- **Protocol flexibility**: Can swap apisvc for GraphQL gateway later
- **Security**: Can add auth at gateway level
- **Scalability**: Scale REST and gRPC services independently

### 3. Why Repository pattern?
- **Testability**: Can mock data access in unit tests
- **Single source of truth**: All queries in one place
- **Easy to switch databases**: Just change repository implementation
- **Follows SOLID principles**

### 4. Why GORM over raw SQL?
- **Type safety**: Compile-time checks
- **Productivity**: Less boilerplate
- **Migrations**: Built-in schema versioning
- **Relations**: Easy to add foreign keys later

---

**For detailed learning guide, see [GO_LEARNING_GUIDE.md](./GO_LEARNING_GUIDE.md)**
