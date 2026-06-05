// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"encoding/binary"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/luxfi/crypto/hash"
)

// referencePacker reproduces the exact byte layout the historical
// github.com/luxfi/codec/wrappers.Packer produced for Prefix / Append
// preimages. It exists only inside this test file to lock down byte-equality
// after dropping the codec dependency.
//
// Reference behavior (codec v1.1.4 wrappers/packer.go):
//   - PackLong  → binary.BigEndian.PutUint64
//   - PackInt   → binary.BigEndian.PutUint32
//   - PackFixedBytes → raw byte copy, no length prefix
type referencePacker struct {
	bytes  []byte
	offset int
}

func (p *referencePacker) packLong(v uint64) {
	binary.BigEndian.PutUint64(p.bytes[p.offset:], v)
	p.offset += 8
}

func (p *referencePacker) packInt(v uint32) {
	binary.BigEndian.PutUint32(p.bytes[p.offset:], v)
	p.offset += 4
}

func (p *referencePacker) packFixedBytes(b []byte) {
	copy(p.bytes[p.offset:], b)
	p.offset += len(b)
}

// referencePrefixPreimage mirrors the pre-refactor ID.Prefix preimage.
func referencePrefixPreimage(id ID, prefixes ...uint64) []byte {
	p := &referencePacker{bytes: make([]byte, len(prefixes)*8+IDLen)}
	for _, prefix := range prefixes {
		p.packLong(prefix)
	}
	p.packFixedBytes(id[:])
	return p.bytes
}

// referenceAppendPreimage mirrors the pre-refactor ID.Append preimage.
func referenceAppendPreimage(id ID, suffixes ...uint32) []byte {
	p := &referencePacker{bytes: make([]byte, IDLen+len(suffixes)*4)}
	p.packFixedBytes(id[:])
	for _, suffix := range suffixes {
		p.packInt(suffix)
	}
	return p.bytes
}

func TestIDPrefix_ByteEqualWithReferencePacker(t *testing.T) {
	require := require.New(t)
	id := GenerateTestID()

	cases := [][]uint64{
		nil,
		{0},
		{1},
		{0xFFFFFFFFFFFFFFFF},
		{1, 256},
		{0, 1, 2, 3, 4, 5, 6, 7},
	}
	for _, prefixes := range cases {
		gotPreimage := make([]byte, len(prefixes)*8+IDLen)
		off := 0
		for _, p := range prefixes {
			binary.BigEndian.PutUint64(gotPreimage[off:], p)
			off += 8
		}
		copy(gotPreimage[off:], id[:])

		refPreimage := referencePrefixPreimage(id, prefixes...)
		require.Equal(refPreimage, gotPreimage,
			"preimage byte mismatch for prefixes=%v", prefixes)

		expected := ID(hash.ComputeHash256Array(refPreimage))
		require.Equal(expected, id.Prefix(prefixes...),
			"Prefix output diverged from reference for prefixes=%v", prefixes)
	}
}

func TestIDAppend_ByteEqualWithReferencePacker(t *testing.T) {
	require := require.New(t)
	id := GenerateTestID()

	cases := [][]uint32{
		nil,
		{0},
		{1},
		{0xFFFFFFFF},
		{1, 256},
		{1, 2, 3, 4, 5, 6, 7, 8},
	}
	for _, suffixes := range cases {
		gotPreimage := make([]byte, IDLen+len(suffixes)*4)
		copy(gotPreimage, id[:])
		off := IDLen
		for _, s := range suffixes {
			binary.BigEndian.PutUint32(gotPreimage[off:], s)
			off += 4
		}

		refPreimage := referenceAppendPreimage(id, suffixes...)
		require.Equal(refPreimage, gotPreimage,
			"preimage byte mismatch for suffixes=%v", suffixes)

		expected := ID(hash.ComputeHash256Array(refPreimage))
		require.Equal(expected, id.Append(suffixes...),
			"Append output diverged from reference for suffixes=%v", suffixes)
	}
}

// TestIDPrefix_GoldenLP77 locks the validationID byte layout described
// in LP-77 (ConvertSubnetToL1Tx). The golden vector is the SHA-256 of:
//
//	id || suffix0 || suffix1 || ...
//
// where id is a known fixed 32-byte input and suffixes are big-endian uint32.
// Any future change that breaks this test breaks LP-77 validationID
// derivation on-chain.
func TestIDAppend_GoldenLP77(t *testing.T) {
	require := require.New(t)

	// Deterministic test vector: id = 0x00..0x1f, suffixes = {0, 1}
	var id ID
	for i := range id {
		id[i] = byte(i)
	}

	preimage := slices.Concat(
		id[:],
		[]byte{0x00, 0x00, 0x00, 0x00},
		[]byte{0x00, 0x00, 0x00, 0x01},
	)
	expected := ID(hash.ComputeHash256Array(preimage))
	require.Equal(expected, id.Append(0, 1))
}
