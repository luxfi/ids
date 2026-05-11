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

## NodeIDScheme — wire-discriminated NodeID derivation

`node_id_scheme.go` owns the strict-PQ NodeID surface introduced in v1.2.10:

- `NodeIDScheme` enum: `MLDSA65=0x42` (canonical strict-PQ),
  `MLDSA87=0x43` (high-value), `Secp256k1=0x90`
  (CLASSICAL_COMPAT_UNSAFE only). Bytes mirror the consensus
  `SigSchemeID` enum so transcripts read the same in both packages.
- `NodeIDScheme.DeriveMLDSA(chainID, pubKey) (NodeID, FullDigest, error)`:
  derives a 48-byte SHAKE256-384 commitment under SP 800-185
  left_encode framing of `("LUX_NODE_ID_V1" || chainID || scheme ||
  pubkey)`. NodeID is `FullDigest[:20]` for storage/map-key use;
  `FullDigest` (48 bytes) is what handshake transcripts and validator
  set roots bind.
- `TypedNodeID = {Scheme, NodeID}` is the wire form: one scheme byte
  followed by the 20-byte NodeID. `ParseTypedNodeID` / `Bytes` are the
  wire codec; `ErrTypedNodeIDLen` / `ErrNodeIDSchemeUnknown` /
  `ErrNodeIDSchemeMismatch` are the typed gate errors.

The 20-byte `NodeID` array stays byte-identical for storage and map
keys; the scheme byte travels alongside it on the wire. Strict-PQ
chains derive their 20-byte NodeID from `SHAKE256-384(...)[:20]`;
classical chains keep their RIPEMD160 derivation via `NodeIDFromCert`.

## Integration with Lux Ecosystem

This package is part of the Lux blockchain ecosystem. See the main documentation at:
- GitHub: https://github.com/luxfi
- Docs: https://docs.lux.network

---

*Auto-generated for AI assistants. Last updated: 2025-12-24*
