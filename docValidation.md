# Adding Protovalidate to a gRPC Go Service

[Protovalidate](https://github.com/bufbuild/protovalidate) lets you declare validation rules directly in `.proto` files and enforce them automatically via a gRPC interceptor — no manual validation code needed in handlers.

---

## Step 1: Add Protovalidate Dependency to `buf.yaml`

```yaml
deps:
  - buf.build/bufbuild/protovalidate
```

Then run:

```bash
buf dep update
```

This fetches the protovalidate proto files and creates/updates `buf.lock`.

---

## Step 2: Update `buf.gen.yaml`

Disable managed mode's `go_package` override for the protovalidate module, so it keeps its own import path:

```yaml
managed:
  enabled: true
  override:
    - file_option: go_package_prefix
      value: github.com/jany/my-coffee/gen/proto
  disable:
    - file_option: go_package
      module: buf.build/bufbuild/protovalidate
```

---

## Step 3: Add Validation Rules in `.proto` Files

Import protovalidate and annotate fields:

```proto
import "buf/validate/validate.proto";

message OrderRequest {
  string menu_item_name = 1 [(buf.validate.field).string.min_len = 1];
}
```

### Common Rule Examples

| Rule | Usage |
|------|-------|
| Required string (non-empty) | `[(buf.validate.field).string.min_len = 1]` |
| Max length | `[(buf.validate.field).string.max_len = 64]` |
| UUID format | `[(buf.validate.field).string.uuid = true]` |
| Email format | `[(buf.validate.field).string.email = true]` |
| Number range | `[(buf.validate.field).uint32.lte = 150]` |
| Repeated min items | `[(buf.validate.field).repeated.min_items = 1]` |

For the full list, see the [rule reference](https://protovalidate.com/reference/rules/).

---

## Step 4: Regenerate Code

```bash
buf generate
```

This re-generates the `.pb.go` files with the validation annotations embedded in the proto descriptors.

---

## Step 5: Install Go Dependencies

```bash
go get buf.build/go/protovalidate
go get github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate
go mod tidy
```

---

## Step 6: Add the Interceptor to the gRPC Server

```go
import (
    "buf.build/go/protovalidate"
    protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
)

// Create a validator
validator, err := protovalidate.New()
if err != nil {
    log.Fatalf("failed to create validator: %v", err)
}

// Pass it as a gRPC unary interceptor
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
)
```

The interceptor automatically validates every incoming unary request **before** it reaches your handler. Invalid requests get rejected with an `InvalidArgument` gRPC status containing the violation details.

---

## How It Works

```
Client Request
    ↓
gRPC Server receives request
    ↓
Protovalidate Interceptor checks request against proto annotations
    ↓
  PASS → Handler (OrderDrink, etc.)
  FAIL → Returns InvalidArgument error with violation details
```

No validation logic in your handler code — it's all driven by the proto annotations.

---

## References

- [Protovalidate docs](https://protovalidate.com/)
- [gRPC + Go quickstart](https://protovalidate.com/quickstart/grpc-go/)
- [Standard rules](https://protovalidate.com/schemas/standard-rules/)
- [Custom CEL rules](https://protovalidate.com/schemas/custom-rules/)
- [protovalidate-go](https://github.com/bufbuild/protovalidate-go)
