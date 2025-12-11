// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

// Native Chain IDs for Lux Network
//
// These are well-known, hardcoded chain IDs for the native chains in Lux.
// Each ID follows a pattern of leading zeros with a distinguishing letter
// for easy visual identification:
//
//   P-Chain: 11111111111111111111111111111111P (Primary/Platform)
//   C-Chain: 11111111111111111111111111111111C (Contract/EVM)
//   X-Chain: 11111111111111111111111111111111X (Exchange/DAG)
//   Q-Chain: 11111111111111111111111111111111Q (Quantum)
//   A-Chain: 11111111111111111111111111111111A (AI)
//   B-Chain: 11111111111111111111111111111111B (Bridge)
//   T-Chain: 11111111111111111111111111111111T (Threshold)
//
// The string representation is for display only - internally these use
// standard 32-byte IDs with the distinguishing byte at position 31.
//
// PERFORMANCE: These IDs provide a fast-path that bypasses base58 encoding/decoding.
// Native chain lookup is O(1) via direct byte comparison.

// Native chain constants - precomputed for maximum speed
const (
	// nativeChainLetterPos is the byte position for the chain letter (last byte)
	nativeChainLetterPos = 31

	// Precomputed string prefix (32 '1's)
	nativeChainPrefix = "11111111111111111111111111111111"

	// Precomputed full strings for each chain
	PChainIDStr = nativeChainPrefix + "P"
	CChainIDStr = nativeChainPrefix + "C"
	XChainIDStr = nativeChainPrefix + "X"
	QChainIDStr = nativeChainPrefix + "Q"
	AChainIDStr = nativeChainPrefix + "A"
	BChainIDStr = nativeChainPrefix + "B"
	TChainIDStr = nativeChainPrefix + "T"
)

var (
	// PChainID is the well-known P-Chain (Platform) ID
	PChainID ID

	// CChainID is the well-known C-Chain (Contract/EVM) ID
	CChainID ID

	// XChainID is the well-known X-Chain (Exchange/DAG) ID
	XChainID ID

	// QChainID is the well-known Q-Chain (Quantum) ID
	QChainID ID

	// AChainID is the well-known A-Chain (AI) ID
	AChainID ID

	// BChainID is the well-known B-Chain (Bridge) ID
	BChainID ID

	// TChainID is the well-known T-Chain (Threshold) ID
	TChainID ID
)

func init() {
	// Initialize the chain IDs - all zeros except last byte
	PChainID[nativeChainLetterPos] = 'P'
	CChainID[nativeChainLetterPos] = 'C'
	XChainID[nativeChainLetterPos] = 'X'
	QChainID[nativeChainLetterPos] = 'Q'
	AChainID[nativeChainLetterPos] = 'A'
	BChainID[nativeChainLetterPos] = 'B'
	TChainID[nativeChainLetterPos] = 'T'
}

// NativeChainString returns the human-friendly string for a native chain ID.
// Returns empty string if not a native chain. This is the fast-path for String().
//
//go:inline
func NativeChainString(id ID) string {
	// Fast check: first 31 bytes must be zero
	// Unrolled loop for speed - check 8 bytes at a time
	if id[0]|id[1]|id[2]|id[3]|id[4]|id[5]|id[6]|id[7] != 0 {
		return ""
	}
	if id[8]|id[9]|id[10]|id[11]|id[12]|id[13]|id[14]|id[15] != 0 {
		return ""
	}
	if id[16]|id[17]|id[18]|id[19]|id[20]|id[21]|id[22]|id[23] != 0 {
		return ""
	}
	if id[24]|id[25]|id[26]|id[27]|id[28]|id[29]|id[30] != 0 {
		return ""
	}

	// Last byte determines the chain
	switch id[nativeChainLetterPos] {
	case 'P':
		return PChainIDStr
	case 'C':
		return CChainIDStr
	case 'X':
		return XChainIDStr
	case 'Q':
		return QChainIDStr
	case 'A':
		return AChainIDStr
	case 'B':
		return BChainIDStr
	case 'T':
		return TChainIDStr
	case 0:
		// All zeros = Empty ID, not a native chain (handled separately)
		return ""
	default:
		return ""
	}
}

// IsNativeChain returns true if the ID is a well-known native chain ID.
// This is the fastest possible check - just verify all zeros except valid last byte.
//
//go:inline
func IsNativeChain(id ID) bool {
	return NativeChainString(id) != ""
}

// NativeChainFromString parses a native chain string and returns the ID.
// Supports full strings (11111111111111111111111111111111P) and aliases (P, p).
// Returns Empty and false if not a native chain string.
//
//go:inline
func NativeChainFromString(s string) (ID, bool) {
	// Fast path: single character aliases
	if len(s) == 1 {
		switch s[0] {
		case 'P', 'p':
			return PChainID, true
		case 'C', 'c':
			return CChainID, true
		case 'X', 'x':
			return XChainID, true
		case 'Q', 'q':
			return QChainID, true
		case 'A', 'a':
			return AChainID, true
		case 'B', 'b':
			return BChainID, true
		case 'T', 't':
			return TChainID, true
		}
		return Empty, false
	}

	// Full string format: 33 chars (32 ones + letter)
	if len(s) != 33 {
		return Empty, false
	}

	// Verify prefix (32 ones)
	if s[:32] != nativeChainPrefix {
		return Empty, false
	}

	// Check last character
	switch s[32] {
	case 'P':
		return PChainID, true
	case 'C':
		return CChainID, true
	case 'X':
		return XChainID, true
	case 'Q':
		return QChainID, true
	case 'A':
		return AChainID, true
	case 'B':
		return BChainID, true
	case 'T':
		return TChainID, true
	}
	return Empty, false
}

// NativeChainAlias returns the single-letter alias for a native chain (P, C, X, Q, A, B, T).
// Returns empty string if not a native chain.
//
//go:inline
func NativeChainAlias(id ID) string {
	if str := NativeChainString(id); str != "" {
		return string(str[32])
	}
	return ""
}

// AllNativeChainIDs returns all well-known native chain IDs.
func AllNativeChainIDs() []ID {
	return []ID{PChainID, CChainID, XChainID, QChainID, AChainID, BChainID, TChainID}
}

// NativeChainIDFromLetter returns the chain ID for a given letter.
// Returns Empty and false if the letter is not a valid chain identifier.
//
//go:inline
func NativeChainIDFromLetter(letter byte) (ID, bool) {
	switch letter {
	case 'P', 'p':
		return PChainID, true
	case 'C', 'c':
		return CChainID, true
	case 'X', 'x':
		return XChainID, true
	case 'Q', 'q':
		return QChainID, true
	case 'A', 'a':
		return AChainID, true
	case 'B', 'b':
		return BChainID, true
	case 'T', 't':
		return TChainID, true
	}
	return Empty, false
}
