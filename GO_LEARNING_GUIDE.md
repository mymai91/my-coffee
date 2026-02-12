# Go Developer Learning Guide - Coffee Shop Microservices 🎓

> **Project**: My Coffee Shop - A production-ready microservices architecture demonstrating Go, gRPC, Connect RPC, GORM, React, and modern API patterns.

---

## 📚 Table of Contents

1. [System Architecture Overview](#system-architecture-overview)
2. [Request Flow Diagrams](#request-flow-diagrams)
3. [Key Go Concepts Demonstrated](#key-go-concepts-demonstrated)
4. [Project Evolution & Learning Path](#project-evolution--learning-path)
5. [Production-Ready Patterns](#production-ready-patterns)
6. [Interview Preparation](#interview-preparation)

---

## 🏗️ System Architecture Overview

### Current Architecture (Connect RPC + REST)

```
┌──────────────────────────────────────────────────────────────────────────┐
│                          CLIENT LAYER                                     │
├────────────────────────────────┬─────────────────────────────────────────┤
│   React Frontend (Vite)        │     CLI Clients                          │
│   localhost:5173               │     (coffeecli)                          │
│   • React Query                │                                          │
│   • TypeScript                 │                                          │
└───────────┬────────────────────┴─────────────────┬────────────────────────┘
            │                                      │
            │ HTTP/REST                            │ gRPC
            │ (JSON)                               │ (Protobuf)
            │                                      │
┌───────────▼──────────────────────────────────────▼────────────────────────┐
│                     API GATEWAY LAYER (apisvc)                            │
│                        localhost:9000                                     │
│  • REST endpoints (/api/*)                                                │
│  • gRPC client connections                                                │
│  • CORS middleware                                                        │
│  • JSON serialization                                                     │
└───────────┬───────────────────────────────────────┬───────────────────────┘
            │                                       │
            │ gRPC                                  │ gRPC
            │ (Protobuf)                            │ (Protobuf)
            │                                       │
┌───────────▼───────────────────┐      ┌───────────▼───────────────────────┐
│   MenuService (menusvc)       │      │   BrewService (brewsvc)           │
│   localhost:50052             │      │   localhost:50051                 │
│                               │      │                                   │
│   Connect RPC Server          │      │   Connect RPC Server              │
│   • HTTP/1.1 + h2c            │      │   • HTTP/1.1 + h2c                │
│   • Supports gRPC protocol    │      │   • Supports gRPC protocol        │
│   • Supports Connect protocol │      │   • Supports Connect protocol     │
│   • Supports gRPC-Web         │      │   • Supports gRPC-Web             │
│                               │      │                                   │
│   Business Logic:             │      │   Business Logic:                 │
│   • In-memory menu data       │      │   • Order management              │
│   • Static coffee items       │      │   • Status tracking               │
│                               │      │   • CRUD operations                │
└───────────────────────────────┘      └───────────┬───────────────────────┘
                                                   │
                                                   │ GORM
                                                   │
                                       ┌───────────▼───────────────────────┐
                                       │   PostgreSQL Database             │
                                       │   localhost:5454                  │
                                       │                                   │
                                       │   Tables:                         │
                                       │   • orders                        │
                                       │     - id (PK)                     │
                                       │     - menu_item_name              │
                                       │     - status                      │
                                       │     - created_at                  │
                                       │     - updated_at                  │
                                       └───────────────────────────────────┘
```

---

## 🔄 Request Flow Diagrams

### Flow 1: Place an Order (POST /api/orders)

```
┌─────────┐                 ┌─────────┐                 ┌─────────┐                 ┌──────────┐
│ React   │                 │ apisvc  │                 │ brewsvc │                 │ Postgres │
│ Client  │                 │ :9000   │                 │ :50051  │                 │ :5454    │
└────┬────┘                 └────┬────┘                 └────┬────┘                 └────┬─────┘
     │                           │                           │                           │
     │ 1. POST /api/orders       │                           │                           │
     │    {"menuItemName":"Latte"}│                          │                           │
     ├──────────────────────────>│                           │                           │
     │                           │                           │                           │
     │                           │ 2. gRPC OrderDrink()      │                           │
     │                           │    OrderRequest{          │                           │
     │                           │      menu_item_name:"Latte"│                          │
     │                           │    }                      │                           │
     │                           ├──────────────────────────>│                           │
     │                           │                           │                           │
     │                           │                           │ 3. Validate Request       │
     │                           │                           │    (Protovalidate)        │
     │                           │                           │                           │
     │                           │                           │ 4. Create Order Model     │
     │                           │                           │    order := &models.Order{│
     │                           │                           │      MenuItemName: "Latte"│
     │                           │                           │      Status: QUEUED       │
     │                           │                           │    }                      │
     │                           │                           │                           │
     │                           │                           │ 5. orderRepo.Create()     │
     │                           │                           ├──────────────────────────>│
     │                           │                           │                           │
     │                           │                           │ 6. INSERT INTO orders ... │
     │                           │                           │<──────────────────────────┤
     │                           │                           │    order.ID = 42          │
     │                           │                           │                           │
     │                           │ 7. OrderResponse{         │                           │
     │                           │      order_id: "order-42" │                           │
     │                           │    }                      │                           │
     │                           │<──────────────────────────┤                           │
     │                           │                           │                           │
     │ 8. 200 OK                 │                           │                           │
     │    {"orderId":"order-42"} │                           │                           │
     │<──────────────────────────┤                           │                           │
     │                           │                           │                           │
```

### Flow 2: Update Order Status (PATCH /api/orders/{id}/status)

```
┌─────────┐                 ┌─────────┐                 ┌─────────┐                 ┌──────────┐
│ Barista │                 │ apisvc  │                 │ brewsvc │                 │ Postgres │
│ CLI     │                 │ :9000   │                 │ :50051  │                 │ :5454    │
└────┬────┘                 └────┬────┘                 └────┬────┘                 └────┬─────┘
     │                           │                           │                           │
     │ 1. PATCH /api/orders/     │                           │                           │
     │    order-42/status        │                           │                           │
     │    {"status":"BREWING"}   │                           │                           │
     ├──────────────────────────>│                           │                           │
     │                           │                           │                           │
     │                           │ 2. Map string to enum     │                           │
     │                           │    "BREWING" ->           │                           │
     │                           │    DrinkStatus_BREWING    │                           │
     │                           │                           │                           │
     │                           │ 3. gRPC UpdateOrderStatus()│                          │
     │                           │    UpdateOrderStatusRequest{│                         │
     │                           │      order_id: "order-42" │                           │
     │                           │      status: BREWING      │                           │
     │                           │    }                      │                           │
     │                           ├──────────────────────────>│                           │
     │                           │                           │                           │
     │                           │                           │ 4. Parse ID               │
     │                           │                           │    "order-42" -> 42       │
     │                           │                           │                           │
     │                           │                           │ 5. orderRepo.FindByID(42) │
     │                           │                           ├──────────────────────────>│
     │                           │                           │                           │
     │                           │                           │ 6. SELECT * FROM orders   │
     │                           │                           │    WHERE id = 42          │
     │                           │                           │<──────────────────────────┤
     │                           │                           │    order found            │
     │                           │                           │                           │
     │                           │                           │ 7. Update status          │
     │                           │                           │    order.Status = BREWING │
     │                           │                           │                           │
     │                           │                           │ 8. orderRepo.Update()     │
     │                           │                           ├──────────────────────────>│
     │                           │                           │                           │
     │                           │                           │ 9. UPDATE orders          │
     │                           │                           │    SET status = 'BREWING' │
     │                           │                           │    WHERE id = 42          │
     │                           │                           │<──────────────────────────┤
     │                           │                           │                           │
     │                           │ 10. UpdateOrderStatusResp │                           │
     │                           │     { order: {...} }      │                           │
     │                           │<──────────────────────────┤                           │
     │                           │                           │                           │
     │ 11. 200 OK                │                           │                           │
     │     {"orderId":"order-42",│                           │                           │
     │      "status":"BREWING"}  │                           │                           │
     │<──────────────────────────┤                           │                           │
     │                           │                           │                           │
```

### Flow 3: Connect RPC Protocol Support

```
┌──────────────────────────────────────────────────────────────────────────┐
│               Connect RPC Server (brewsvc/menusvc)                       │
│                    http.Server with h2c                                  │
└────────────────────────┬─────────────────────────────────────────────────┘
                         │
                         │ Incoming Request
                         │
         ┌───────────────┴──────────────┬──────────────────┐
         │                              │                  │
         ▼                              ▼                  ▼
┌─────────────────┐          ┌──────────────────┐  ┌──────────────────┐
│ gRPC Protocol   │          │ Connect Protocol │  │ gRPC-Web Protocol│
│ (from apisvc)   │          │ (from browsers)  │  │ (from browsers)  │
├─────────────────┤          ├──────────────────┤  ├──────────────────┤
│ • HTTP/2        │          │ • HTTP/1.1 or 2  │  │ • HTTP/1.1       │
│ • Content-Type: │          │ • Content-Type:  │  │ • Content-Type:  │
│   application/  │          │   application/   │  │   application/   │
│   grpc          │          │   connect+proto  │  │   grpc-web+proto │
│ • Protobuf      │          │ • Protobuf/JSON  │  │ • Protobuf       │
└────────┬────────┘          └────────┬─────────┘  └────────┬─────────┘
         │                            │                     │
         └────────────────────────────┴─────────────────────┘
                                      │
                                      ▼
                     ┌────────────────────────────────┐
                     │  brewconnect.BrewServiceHandler│
                     │                                │
                     │  All methods receive:          │
                     │  *connect.Request[T]           │
                     │                                │
                     │  All methods return:           │
                     │  *connect.Response[T]          │
                     └────────────────────────────────┘
```

---

## 💡 Key Go Concepts Demonstrated

### 1. **Interfaces & Type Assertions**

```go
// File: internal/brews/brew.go

// Compile-time check that Server implements the interface
var _ brewconnect.BrewServiceHandler = (*Server)(nil)

type Server struct {
    orderRepo *repository.OrderRepository
}

// Interface implementation - no explicit "implements" keyword in Go
func (s *Server) OrderDrink(
    ctx context.Context,
    req *connect.Request[brewpb.OrderRequest],
) (*connect.Response[brewpb.OrderResponse], error) {
    // Implementation
}
```

**Learning Points:**
- Go uses **implicit interfaces** (duck typing)
- No `implements` keyword needed
- Use `var _ Interface = (*Type)(nil)` for compile-time checks
- Methods with receiver `(s *Server)` implement the interface

---

### 2. **Error Handling & Custom Errors**

```go
// Before (plain gRPC):
if err != nil {
    return nil, fmt.Errorf("failed to create order: %w", err)
}

// After (Connect RPC with proper error codes):
if err != nil {
    log.Printf("Failed to create order: %v", err)
    return nil, connect.NewError(connect.CodeInternal, 
        fmt.Errorf("failed to create order: %w", err))
}

// Different error codes for different scenarios:
connect.CodeInternal       // Server error (500)
connect.CodeNotFound       // Resource not found (404)
connect.CodeInvalidArgument // Bad request (400)
```

**Learning Points:**
- Always wrap errors with `%w` for error chains
- Use appropriate error codes (not just generic errors)
- Log before returning errors for debugging
- Connect RPC maps to HTTP status codes automatically

---

### 3. **Repository Pattern (Data Access Layer)**

```go
// File: internal/repository/order_repository.go

type OrderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
    return &OrderRepository{db: db}
}

// Single responsibility - each method does ONE thing
func (r *OrderRepository) Create(order *models.Order) error {
    return r.db.Create(order).Error
}

func (r *OrderRepository) FindByID(id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.First(&order, id).Error
    if err != nil {
        return nil, err
    }
    return &order, nil
}
```

**Learning Points:**
- **Repository Pattern** abstracts database operations
- Each repository method has a single responsibility
- Makes testing easier (can mock repository)
- Clear separation between data access and business logic

---

### 4. **GORM ORM Patterns**

```go
// File: internal/models/order.go

type Order struct {
    ID           uint `gorm:"primaryKey"`
    MenuItemName string `gorm:"not null"`
    Status       OrderStatus `gorm:"default:QUEUED"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Custom table name
func (Order) TableName() string {
    return "orders"
}

// Custom type for enum-like behavior
type OrderStatus string

const (
    StatusQueued   OrderStatus = "QUEUED"
    StatusGrinding OrderStatus = "GRINDING"
    StatusBrewing  OrderStatus = "BREWING"
    StatusFrothing OrderStatus = "FROTHING"
    StatusReady    OrderStatus = "READY"
)
```

**Learning Points:**
- GORM tags define database constraints
- `CreatedAt`/`UpdatedAt` auto-managed by GORM
- Custom types for type safety (OrderStatus vs string)
- Constants for enum-like behavior

---

### 5. **Dependency Injection**

```go
// File: cmd/brewsvc/main.go

func main() {
    // Load config
    config.Load()

    // Create dependencies
    db := database.Connect()
    defer database.Close()

    // Inject dependencies into service
    brewService := brews.New(db)  // <- Dependency injection

    // Create handler with injected service
    path, handler := brewconnect.NewBrewServiceHandler(brewService)
    mux.Handle(path, handler)
}
```

**Learning Points:**
- Dependencies passed explicitly (not hidden globals)
- Easy to test (inject mock database)
- Clear initialization order
- Constructor functions (`New()`) for clean setup

---

### 6. **Middleware & Interceptors**

```go
// File: cmd/apisvc/main.go

// CORS middleware
func cors(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Usage:
http.ListenAndServe(":9000", cors(mux))
```

```go
// File: cmd/brewsvc/main.go

// Connect RPC interceptor
validateInterceptor := validate.NewInterceptor()
path, handler := brewconnect.NewBrewServiceHandler(
    brewService,
    connect.WithInterceptors(validateInterceptor), // <- Interceptor
)
```

**Learning Points:**
- Middleware wraps handlers for cross-cutting concerns
- Chain multiple middlewares
- Interceptors for gRPC/Connect RPC validation
- Separation of concerns (auth, logging, validation separate)

---

### 7. **Context Usage**

```go
func (s *Server) OrderDrink(
    ctx context.Context,  // <- Context always first parameter
    req *connect.Request[brewpb.OrderRequest],
) (*connect.Response[brewpb.OrderResponse], error) {
    // Context can carry:
    // - Cancellation signals
    // - Deadlines
    // - Request-scoped values (user ID, trace ID)
    
    // Example: pass context to repository
    order, err := s.orderRepo.CreateWithContext(ctx, order)
}
```

**Learning Points:**
- `context.Context` is Go's way to handle cancellation, timeouts, deadlines
- Always first parameter by convention
- Pass context down the call stack
- Used in production for request tracing, timeouts

---

### 8. **Protocol Buffers & Code Generation**

```protobuf
// File: proto/brew/brew.proto

syntax = "proto3";
package brew;

import "buf/validate/validate.proto";

service BrewService {
  rpc OrderDrink (OrderRequest) returns (OrderResponse);
  rpc GetOrder (GetOrderRequest) returns (GetOrderResponse);
  rpc UpdateOrderStatus (UpdateOrderStatusRequest) returns (UpdateOrderStatusResponse);
  rpc DeleteOrder (DeleteOrderRequest) returns (DeleteOrderResponse);
  rpc ListOrders (ListOrdersRequest) returns(ListOrdersResponse);
}

message OrderRequest {
  string menu_item_name = 1 [(buf.validate.field).string.min_len = 1];
}
```

**Generated code creates:**
- Go structs for messages
- Client interfaces
- Server interfaces
- Serialization/deserialization
- Validation logic

**Learning Points:**
- `.proto` files are single source of truth
- Code generated automatically
- Type safety across languages
- Built-in validation with protovalidate

---

### 9. **HTTP/2 Server with h2c**

```go
// File: cmd/brewsvc/main.go

// Setup h2c (HTTP/2 Cleartext) for gRPC compatibility
p := new(http.Protocols)
p.SetHTTP1(true)              // Support HTTP/1.1
p.SetUnencryptedHTTP2(true)   // Support HTTP/2 without TLS (h2c)

server := http.Server{
    Addr:      ":50051",
    Handler:   mux,
    Protocols: p,  // <- Enables multi-protocol support
}
```

**Learning Points:**
- h2c = HTTP/2 without TLS (for local development)
- Enables gRPC compatibility
- Multi-protocol support (HTTP/1.1 + HTTP/2)
- Production would use TLS for security

---

### 10. **Database Migrations**

```sql
-- File: migrations/000001_create_orders_table.up.sql

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    menu_item_name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'QUEUED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```go
// File: cmd/migrate/main.go

m, err := migrate.New(
    "file://migrations",
    databaseURL,
)
if err != nil {
    log.Fatal(err)
}

if err := m.Up(); err != nil && err != migrate.ErrNoChange {
    log.Fatal(err)
}
```

**Learning Points:**
- Version control for database schema
- Up/down migrations for rollback
- `golang-migrate` CLI for automation
- SQL-based migrations for clarity

---

## 🚀 Project Evolution & Learning Path

### Phase 1: Basic gRPC Service (Feature 1-3)
**Commit:** `Merge branch 'feature-3-brewsvc'`

```
✅ Learn:
- Go modules (go.mod)
- Protocol buffers
- gRPC server setup
- Basic CRUD with GORM
- Repository pattern
```

---

### Phase 2: Validation (Feature 4-5)
**Commit:** `Add build target for brew service in Makefile`

```
✅ Learn:
- Protovalidate for input validation
- Buf tool for protobuf management
- buf.gen.yaml configuration
- Validation in .proto files
```

---

### Phase 3: CLI Client (Feature 6)
**Commit:** `feat: initialize web application with React, Vite, and TypeScript`

```
✅ Learn:
- gRPC client implementation
- Context and cancellation
- CLI tool development
- gRPC service consumption
```

---

### Phase 4: REST Gateway (Feature 7)
**Commit:** `feat: add order management endpoints`

```
✅ Learn:
- REST API design
- JSON marshaling/unmarshaling
- HTTP routing with ServeMux
- CORS middleware
- gRPC to REST translation
```

---

### Phase 5: Frontend Integration
**Commit:** `feat: integrate React Query for data fetching`

```
✅ Learn:
- React Query for state management
- API client patterns
- Error handling in frontend
- TypeScript integration
```

---

### Phase 6: Connect RPC Migration (Current)
**Commit:** `feat: migrate to Connect RPC for brew and menu services`

```
✅ Learn:
- Connect RPC framework
- Multi-protocol support (gRPC, Connect, gRPC-Web)
- HTTP/2 and h2c
- Modern RPC patterns
- Interceptors for validation
- Error code mapping
```

---

## 🎯 Production-Ready Patterns

### 1. **Configuration Management**

```go
// File: config/config.go

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
}

func Load() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}
```

**Production Best Practices:**
- Environment variables for configuration
- `.env` for local development
- Never commit secrets to git
- Use secret management in production (AWS Secrets Manager, Vault)

---

### 2. **Structured Logging**

```go
// Current: Standard log
log.Printf("Failed to create order: %v", err)

// Production: Structured logging (recommendation)
import "go.uber.org/zap"

logger.Error("failed to create order",
    zap.String("menu_item", req.MenuItemName),
    zap.Error(err),
)
```

---

### 3. **Graceful Shutdown**

```go
// Production pattern
func main() {
    server := &http.Server{Addr: ":50051", Handler: mux}

    // Run server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
}
```

---

### 4. **Health Checks**

```go
// Add to apisvc
mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
    })
})
```

---

### 5. **Observability (Metrics & Tracing)**

```go
// Production: Add Prometheus metrics
import "github.com/prometheus/client_golang/prometheus/promhttp"

mux.Handle("/metrics", promhttp.Handler())

// Add OpenTelemetry tracing
import "go.opentelemetry.io/otel"
```

---

## 📝 Interview Preparation

### Common Go Interview Questions from This Project

#### **Q1: Explain the difference between gRPC and REST**

**Answer:**
- **gRPC**: Uses HTTP/2, binary protocol (Protobuf), strongly typed, supports bidirectional streaming
- **REST**: Uses HTTP/1.1, text protocol (JSON/XML), loosely typed, request-response only
- **When to use gRPC**: Microservice-to-microservice communication, performance critical
- **When to use REST**: Browser clients, simplicity, caching, third-party integrations

---

#### **Q2: How does Connect RPC improve on plain gRPC?**

**Answer:**
- Connect RPC supports **3 protocols**: gRPC, Connect, gRPC-Web
- Works over HTTP/1.1 (doesn't require HTTP/2)
- Browser-friendly (no need for proxy like Envoy)
- Same `.proto` definition generates multiple protocols
- Easier testing with curl (supports JSON)

---

#### **Q3: Explain the Repository Pattern**

**Answer:**
```go
// Repository abstracts data access
type OrderRepository struct {
    db *gorm.DB
}

// Business logic doesn't know about database
type Server struct {
    orderRepo *repository.OrderRepository
}

// Benefits:
// - Testability (mock repository)
// - Separation of concerns
// - Easy to swap database (Postgres -> MySQL)
// - Single source for queries
```

---

#### **Q4: How do you handle errors in Go?**

**Answer:**
```go
// 1. Always check errors
if err != nil {
    // handle it
}

// 2. Wrap errors for context
return fmt.Errorf("failed to create order: %w", err)

// 3. Use proper error codes (Connect RPC)
return connect.NewError(connect.CodeNotFound, err)

// 4. Log before returning
log.Printf("Error: %v", err)
return err
```

---

#### **Q5: Explain dependency injection in this project**

**Answer:**
```go
// 1. Create dependency
db := database.Connect()

// 2. Inject into service
brewService := brews.New(db)  // <- DI here

// 3. Inject into handler
handler := brewconnect.NewBrewServiceHandler(brewService)

// Benefits:
// - Testable (inject mock DB)
// - No global state
// - Explicit dependencies
```

---

#### **Q6: What is context.Context and why is it important?**

**Answer:**
- `context.Context` carries cancellation signals, deadlines, request-scoped values
- Always first parameter in functions
- Prevents goroutine leaks (cancellation propagates)
- Used for request tracing in production
- Example: HTTP request timeout propagates to DB query

```go
func (s *Server) OrderDrink(ctx context.Context, req *Request) (*Response, error) {
    // If client cancels request, ctx is cancelled
    // DB operations respect ctx and stop
    err := s.orderRepo.CreateWithContext(ctx, order)
}
```

---

#### **Q7: How does GORM handle timestamps?**

**Answer:**
```go
type Order struct {
    CreatedAt time.Time  // Auto-set on creation
    UpdatedAt time.Time  // Auto-updated on save
}

// GORM hooks:
// - BeforeCreate sets CreatedAt
// - BeforeUpdate sets UpdatedAt
// No manual timestamp management needed
```

---

#### **Q8: Explain the three-layer architecture in this project**

**Answer:**
```
1. Handler Layer (brews.Server)
   - Receives requests
   - Validates input
   - Returns responses

2. Business Logic Layer (in handler methods)
   - Order creation logic
   - Status validation
   - Business rules

3. Data Access Layer (OrderRepository)
   - Database operations
   - CRUD methods
   - No business logic
```

---

### Code Challenge: Add a New Feature

**Task:** Add a "Cancel Order" endpoint

```go
// 1. Update proto/brew/brew.proto
message CancelOrderRequest {
  string order_id = 1 [(buf.validate.field).string.min_len = 1];
}

message CancelOrderResponse {
  bool success = 1;
  string message = 2;
}

service BrewService {
  // ... existing methods
  rpc CancelOrder (CancelOrderRequest) returns (CancelOrderResponse);
}

// 2. Regenerate proto code
// make proto

// 3. Implement in internal/brews/brew.go
func (s *Server) CancelOrder(
    ctx context.Context,
    req *connect.Request[brewpb.CancelOrderRequest],
) (*connect.Response[brewpb.CancelOrderResponse], error) {
    var orderID uint
    if _, err := fmt.Sscanf(req.Msg.OrderId, "order-%d", &orderID); err != nil {
        return nil, connect.NewError(connect.CodeInvalidArgument, err)
    }

    order, err := s.orderRepo.FindByID(orderID)
    if err != nil {
        return nil, connect.NewError(connect.CodeNotFound, err)
    }

    // Business rule: can only cancel if not READY
    if order.Status == models.StatusReady {
        return connect.NewResponse(&brewpb.CancelOrderResponse{
            Success: false,
            Message: "Cannot cancel order that is ready",
        }), nil
    }

    if err := s.orderRepo.Delete(orderID); err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    return connect.NewResponse(&brewpb.CancelOrderResponse{
        Success: true,
        Message: "Order cancelled successfully",
    }), nil
}

// 4. Add REST endpoint in cmd/apisvc/main.go
mux.HandleFunc("POST /api/orders/{orderId}/cancel", func(w http.ResponseWriter, r *http.Request) {
    orderId := r.PathValue("orderId")
    
    resp, err := brewClient.CancelOrder(context.Background(), 
        connect.NewRequest(&brewpb.CancelOrderRequest{
            OrderId: orderId,
        }))
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "success": resp.Msg.Success,
        "message": resp.Msg.Message,
    })
})
```

---

## 🎓 Next Steps for Learning

1. **Add Authentication**
   - JWT tokens
   - Middleware for auth
   - User context

2. **Add Testing**
   - Unit tests with mocks
   - Integration tests
   - Table-driven tests

3. **Add Observability**
   - Prometheus metrics
   - OpenTelemetry tracing
   - Structured logging (zap/logrus)

4. **Add CI/CD**
   - GitHub Actions
   - Docker containerization
   - Kubernetes deployment

5. **Add Advanced Features**
   - Rate limiting
   - Caching (Redis)
   - Message queue (RabbitMQ/Kafka)
   - Websockets for real-time updates

---

## 📚 Recommended Resources

### Books
- **"The Go Programming Language"** by Donovan & Kernighan
- **"Go in Action"** by William Kennedy
- **"Building Microservices"** by Sam Newman

### Online
- Official Go Tour: https://go.dev/tour/
- Connect RPC Docs: https://connectrpc.com/docs/
- gRPC Go: https://grpc.io/docs/languages/go/
- GORM Docs: https://gorm.io/docs/

### Practice
- LeetCode Go problems
- Build another microservice
- Contribute to open source Go projects

---

## ✅ Job-Ready Checklist

- [ ] Understand Go basics (structs, interfaces, goroutines, channels)
- [ ] Can explain this project architecture
- [ ] Understand gRPC vs REST trade-offs
- [ ] Know repository pattern and DI
- [ ] Can write unit tests
- [ ] Understand database migrations
- [ ] Know error handling best practices
- [ ] Can explain context.Context usage
- [ ] Understand middleware/interceptors
- [ ] Know production deployment basics (Docker, env vars)

---

**Good luck with your Go developer interviews! 🚀**
