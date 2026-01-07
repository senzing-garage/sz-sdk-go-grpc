# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

sz-sdk-go-grpc is a Go SDK that provides gRPC client implementations for Senzing entity resolution services. It implements the interfaces defined in [sz-sdk-go](https://github.com/senzing-garage/sz-sdk-go) to communicate with a remote [Senzing gRPC server](https://github.com/senzing-garage/servegrpc).

## Build and Development Commands

```bash
# Install development tools (one-time)
make dependencies-for-development

# Update Go dependencies
make dependencies

# Lint (runs golangci-lint, govulncheck, cspell)
make lint

# Run tests (requires Docker - starts gRPC server container)
make clean setup test

# Run a single test
go test -v -run TestFunctionName ./package/...

# Test with coverage
make clean setup coverage

# Auto-fix lint issues
make fix
```

### TLS Testing

```bash
# Server-side TLS
make clean setup-server-side-tls test-server-side-tls

# Mutual TLS
make clean setup-mutual-tls test-mutual-tls
```

## Architecture

### Package Structure

- **szabstractfactory/** - Abstract factory for creating all service clients from a single gRPC connection
- **szengine/** - Entity resolution engine (AddRecord, GetEntity, SearchByAttributes, WhyRecords, etc.)
- **szconfig/** - In-memory configuration management
- **szconfigmanager/** - Persistent configuration lifecycle management
- **szproduct/** - Product/license information
- **szdiagnostic/** - System diagnostics
- **helper/** - Shared utilities: error conversion, TLS credentials, logging

### Client Implementation Pattern

Each client follows a two-layer delegation pattern:

1. **Public method** - Handles tracing, error wrapping, and async observer notifications
2. **Private method** - Builds protobuf request, calls gRPC, converts errors

```go
// Public: Szengine.AddRecord() -> logging + delegation + observer notify
// Private: Szengine.addRecord() -> protobuf request -> gRPC call -> error conversion
```

All public methods implement interfaces from `github.com/senzing-garage/sz-sdk-go/senzing`.

### Error Handling

gRPC errors are converted to Senzing nested errors via `helper.ConvertGrpcError()`. Always wrap errors with `wraperror.Errorf()`.

### TLS Configuration

Transport credentials are configured via environment variables:

- `SENZING_TOOLS_SERVER_CA_CERTIFICATE_FILE` - Server CA cert (enables TLS)
- `SENZING_TOOLS_CLIENT_CERTIFICATE_FILE` - Client cert (enables mutual TLS)
- `SENZING_TOOLS_CLIENT_KEY_FILE` - Client private key
- `SENZING_TOOLS_CLIENT_KEY_PASSPHRASE` - Optional key passphrase

### Component IDs

Each package has a unique ID for error messages and logging:

- szengine: 6020
- szconfig: 6021
- szconfigmanager: 6022
- szdiagnostic: 6023
- szproduct: 6026

## Code Style

- Go version: 1.24.4
- Line length limit: 120 characters
- Uses extensive golangci-lint configuration (see `.github/linters/.golangci.yaml`)
- JSON field tags use upperSnake case
- Tests run with `-p 1` (sequential) due to shared gRPC server state
