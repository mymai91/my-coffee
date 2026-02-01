# Go Tutorial: Building a Coffee Shop System
## 📋 Project Overview
The coffee project is a microservices system with:

2 gRPC Services: menusvc (menu) and brewsvc (orders)
2 CLI Apps: coffeecli (customer) and baristacli (barista)

┌─────────────┐      gRPC        ┌─────────────┐
│  coffeecli  │ ◄──────────────► │  menusvc    │ (port 50052)
│  (Customer) │                  └─────────────┘
│             │      gRPC        ┌─────────────┐
│             │ ◄──────────────► │  brewsvc    │ (port 50051)
└─────────────┘                  └─────────────┘

## 🚀 Step 1: Setup Your Project

Created `go.mod`

```
module github.com/jany/my-coffee

go 1.24.0

require (
	connectrpc.com/connect v1.19.1
	google.golang.org/grpc v1.78.0
	google.golang.org/protobuf v1.36.11
)
```

## 🚀 Step 2: Create Proto Files (API Contracts)

What is Protocol Buffers (Protobuf)?

A language-neutral way to define your API
Generates Go code automatically
Smaller and faster than JSON
Key concepts in proto files:

```
syntax = "proto3";           // Version of protobuf
package menu;                // Namespace
option go_package = "...";   // Where Go code goes

service MenuService {        // Defines the gRPC service
  rpc GetMenu(...) returns (...);  // Method definition
}

message MenuItem {           // Data structure (like a struct)
  string name = 1;           // Field with tag number
  double price = 2;
  repeated Item items = 3;   // "repeated" = array/slice
}

```

## 🔧 Step 3: Setup Proto Generation (Buf)

**Buf** is a modern tool for working with Protocol Buffers. It:
- Validates your proto files
- Generates Go code from proto files
- Manages dependencies ---

## 🎯 Step 4: Implement Menu Service

Now let's write the **business logic**. This is where the real Go learning begins!

### Key Go Concepts:

1. **Package** - Groups related code together
2. **Struct** - Custom data type with fields
3. **Method** - Function attached to a struct
4. **Interface** - Contract that types can implement
5. **Pointer receiver**  - Method can modify the struct ---

## 🎯 Step 5: Implement Brew Service

This service is more complex - it manages **state** (orders) and needs **concurrency safety**.

### Key Go Concepts:

1. **`sync.RWMutex`** - Protects shared data from race conditions
   - `RLock()` - Read lock (multiple readers allowed)
   - `Lock()` - Write lock (exclusive access)
2. **`defer`** - Ensures code runs when function exits (like `finally`)
3. **`map`** - Go's hash table / dictionary ---

## 🖥️ Step 6: Create gRPC Servers

Now we create the **entry points** - the  functions that start our services.

### Key Go Concepts:

1. ** package** - Entry point of executable programs
2. **** - Opens a network port
3. **** - Creates a gRPC server
4. **Error handling** - Always check errors in Go! ---

## 🧪 Step 7: Add Unit Tests

Testing is crucial in Go! Let's write tests for our services.

### Key Go Testing Concepts:

1. **`_test.go`** - Test files must end with `_test.go`
2. **`testing.T`** - Provides test methods like `Error()`, `Fatal()`
3. **Table-driven tests** - Testing multiple cases in a loop
4. **** - Creates subtests for better organization ---

## 🔨 Step 8: Create Makefile

A Makefile makes it easy to run common commands. ---

## 🎉 Summary: Your Project Structure

Your  project is now complete:
