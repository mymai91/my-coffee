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
