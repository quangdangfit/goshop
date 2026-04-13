---
description: Proto/gRPC definitions, code generation, and plugin requirements
globs:
  - "proto/**"
alwaysApply: false
---

# Proto / gRPC

Proto definitions live in `proto/`. Generated Go code is in `proto/gen/go/`. To regenerate:

```bash
cd proto && make build
# Equivalent: buf generate (uses buf.gen.yaml with local protoc-gen-go and protoc-gen-go-grpc plugins)
```

Local plugins required: `protoc-gen-go@v1.30.0` and `protoc-gen-go-grpc@v1.3.0` (must match `google.golang.org/grpc v1.57.0` and `google.golang.org/protobuf v1.30.0` in go.mod). Install via:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

The gRPC server registers `UserService`, `CartService`, `ProductService`, and `OrderService`.
