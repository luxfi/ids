# ids

## Overview

The `ids` package provides strongly-typed identifiers for the Lux Network. It includes implementations for various ID types used throughout the ecosystem, ensuring type safety and preventing ID misuse. - **Type-Safe IDs**: Prevent mixing different ID types at compile time 

## Package Information

- **Type**: go
- **Module**: github.com/luxfi/ids
- **Repository**: github.com/luxfi/ids

## Directory Structure

```
.
cmd
cmd/test_native
idstest
utils
utils/wrappers
```

## Key Files

- aliases_test.go
- aliases.go
- basic_test.go
- bits_test.go
- bits.go
- go.mod
- id_test.go
- id.go
- native_chains.go
- node_id_test.go

## Development

### Prerequisites

- Go 1.21+

### Build

```bash
go build ./...
```

### Test

```bash
go test -v ./...
```

## Integration with Lux Ecosystem

This package is part of the Lux blockchain ecosystem. See the main documentation at:
- GitHub: https://github.com/luxfi
- Docs: https://docs.lux.network

---

*Auto-generated for AI assistants. Last updated: 2025-12-24*
