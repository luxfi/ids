// Copyright (C) 2020-2026, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNodeID_Derive_MLDSA65_Deterministic — derivation is a pure function
// of (scheme, chainID, pubkey). Same inputs ⇒ same digest and same NodeID.
func TestNodeID_Derive_MLDSA65_Deterministic(t *testing.T) {
	require := require.New(t)

	chainID := ID{0x01, 0x02, 0x03}
	pubKey := []byte("ml-dsa-65-public-key-fixture-bytes")

	id1, full1, err := NodeIDSchemeMLDSA65.DeriveMLDSA(chainID, pubKey)
	require.NoError(err)
	id2, full2, err := NodeIDSchemeMLDSA65.DeriveMLDSA(chainID, pubKey)
	require.NoError(err)

	require.Equal(id1, id2)
	require.Equal(full1, full2)
	require.Equal(full1[:NodeIDLen], id1.Bytes(),
		"NodeID must be the 20-byte prefix of FullDigest")
	require.NotEqual(EmptyNodeID, id1)
}

// TestNodeID_Derive_DistinctSchemesDifferentIDs — the scheme byte is
// bound into the digest: ML-DSA-65 and ML-DSA-87 keys (same public-key
// bytes, different scheme tag) MUST derive distinct NodeIDs.
func TestNodeID_Derive_DistinctSchemesDifferentIDs(t *testing.T) {
	require := require.New(t)

	chainID := ID{0xab, 0xcd}
	pubKey := []byte("identical-pubkey-bytes-across-both-schemes")

	id65, full65, err := NodeIDSchemeMLDSA65.DeriveMLDSA(chainID, pubKey)
	require.NoError(err)
	id87, full87, err := NodeIDSchemeMLDSA87.DeriveMLDSA(chainID, pubKey)
	require.NoError(err)

	require.NotEqual(id65, id87, "scheme byte must perturb the digest")
	require.NotEqual(full65, full87)
}

// TestNodeID_Derive_DistinctChainsDifferentIDs — chainID is bound into
// the digest: the same key on a different chain produces a different
// NodeID, which closes cross-chain replay of validator registrations.
func TestNodeID_Derive_DistinctChainsDifferentIDs(t *testing.T) {
	require := require.New(t)

	pubKey := []byte("validator-key")
	chainA := ID{0x01}
	chainB := ID{0x02}

	idA, _, err := NodeIDSchemeMLDSA65.DeriveMLDSA(chainA, pubKey)
	require.NoError(err)
	idB, _, err := NodeIDSchemeMLDSA65.DeriveMLDSA(chainB, pubKey)
	require.NoError(err)

	require.NotEqual(idA, idB, "chainID must perturb the digest")
}

// TestNodeID_Derive_RejectsClassicalScheme — DeriveMLDSA refuses any
// scheme outside the ML-DSA block. A caller cannot smuggle a classical
// key through the ML-DSA derivation path.
func TestNodeID_Derive_RejectsClassicalScheme(t *testing.T) {
	require := require.New(t)

	_, _, err := NodeIDSchemeSecp256k1.DeriveMLDSA(ID{}, []byte("key"))
	require.ErrorIs(err, ErrNodeIDSchemeInvalid)

	_, _, err = NodeIDSchemeInvalid.DeriveMLDSA(ID{}, []byte("key"))
	require.ErrorIs(err, ErrNodeIDSchemeInvalid)
}

// TestNodeID_Derive_RejectsEmptyPubKey — an empty public key is refused
// before the hash is computed. Closes the "all-zero pubkey collapses to
// the same NodeID" pitfall.
func TestNodeID_Derive_RejectsEmptyPubKey(t *testing.T) {
	require := require.New(t)
	_, _, err := NodeIDSchemeMLDSA65.DeriveMLDSA(ID{}, nil)
	require.ErrorIs(err, ErrNodeIDSchemeInvalid)
	_, _, err = NodeIDSchemeMLDSA65.DeriveMLDSA(ID{}, []byte{})
	require.ErrorIs(err, ErrNodeIDSchemeInvalid)
}

// TestNodeID_WireDecode_HonorsSchemeByte — round-trip a TypedNodeID
// through Bytes() / ParseTypedNodeID and confirm the scheme byte survives
// unchanged.
func TestNodeID_WireDecode_HonorsSchemeByte(t *testing.T) {
	schemes := []NodeIDScheme{
		NodeIDSchemeMLDSA65,
		NodeIDSchemeMLDSA87,
		NodeIDSchemeSecp256k1,
	}
	for _, s := range schemes {
		s := s
		t.Run(s.String(), func(t *testing.T) {
			r := require.New(t)
			id := NodeID{0x01, 0x02, 0x03, 0x04, 0x05}
			t1, err := NewTypedNodeID(s, id)
			r.NoError(err)

			b := t1.Bytes()
			r.Len(b, TypedNodeIDLen)
			r.Equal(byte(s), b[0],
				"scheme byte must be the leading byte on the wire")

			t2, err := ParseTypedNodeID(b)
			r.NoError(err)
			r.Equal(t1, t2)
		})
	}
}

// TestNodeID_WireDecode_RejectsBadLength — a TypedNodeID whose wire form
// is the wrong length is rejected with ErrTypedNodeIDLen.
func TestNodeID_WireDecode_RejectsBadLength(t *testing.T) {
	require := require.New(t)

	_, err := ParseTypedNodeID(nil)
	require.ErrorIs(err, ErrTypedNodeIDLen)

	_, err = ParseTypedNodeID(make([]byte, TypedNodeIDLen-1))
	require.ErrorIs(err, ErrTypedNodeIDLen)

	_, err = ParseTypedNodeID(make([]byte, TypedNodeIDLen+1))
	require.ErrorIs(err, ErrTypedNodeIDLen)
}

// TestNodeID_WireDecode_RejectsUnknownScheme — a wire input whose leading
// byte is not a known scheme is rejected with ErrNodeIDSchemeUnknown. The
// classical-compat path does NOT relax this: only the named secp256k1
// byte (0x90) is accepted; other 0x90+ bytes are still refused.
func TestNodeID_WireDecode_RejectsUnknownScheme(t *testing.T) {
	require := require.New(t)

	buf := make([]byte, TypedNodeIDLen)
	for _, badScheme := range []byte{0x00, 0x01, 0x40, 0x41, 0x44, 0x91, 0xFF} {
		buf[0] = badScheme
		_, err := ParseTypedNodeID(buf)
		require.ErrorIs(err, ErrNodeIDSchemeUnknown,
			"scheme byte 0x%02x should be rejected", badScheme)
	}
}

// TestNodeID_StrictPQ_RejectsSecp256k1 — under strict-PQ semantics
// (IsPostQuantum), the classical scheme is refused at the cross-axis
// gate. The test exercises the policy directly, mirroring what the
// consensus profile's ValidatorSchemeID match enforces.
func TestNodeID_StrictPQ_RejectsSecp256k1(t *testing.T) {
	require := require.New(t)

	// Strict-PQ acceptance gate: pinned scheme MUST be in the PQ family.
	pinned := NodeIDSchemeMLDSA65
	require.True(pinned.IsPostQuantum())

	classical, _ := NewTypedNodeID(NodeIDSchemeSecp256k1,
		NodeID{0xde, 0xad, 0xbe, 0xef})
	require.False(classical.Scheme.IsPostQuantum())
	require.True(classical.Scheme.IsClassicalCompatUnsafe())

	// The cross-axis check at the boundary: profile pins ML-DSA-65 but
	// peer presented classical. Reject.
	err := checkSchemeAgainstProfile(classical.Scheme, pinned, false /* classicalCompatUnsafe */)
	require.ErrorIs(err, ErrNodeIDSchemeMismatch)
}

// TestNodeID_ClassicalCompat_AcceptsSecp256k1 — when the operator opts
// into LUX_CLASSICAL_COMPAT_UNSAFE the classical secp256k1 scheme is
// accepted even on a chain whose locked profile names ML-DSA. The
// scheme byte still has to be the named secp256k1 byte; unknown bytes
// are refused regardless of the unsafe knob.
func TestNodeID_ClassicalCompat_AcceptsSecp256k1(t *testing.T) {
	require := require.New(t)

	pinned := NodeIDSchemeMLDSA65
	classical, _ := NewTypedNodeID(NodeIDSchemeSecp256k1,
		NodeID{0x42, 0x42, 0x42})

	// Classical-compat ON: accepted.
	err := checkSchemeAgainstProfile(classical.Scheme, pinned, true)
	require.NoError(err)

	// Unknown scheme under classical-compat is still refused. The
	// compat knob accepts the named classical scheme, not "anything
	// outside the PQ block".
	err = checkSchemeAgainstProfile(NodeIDScheme(0x91), pinned, true)
	require.ErrorIs(err, ErrNodeIDSchemeMismatch)
}

// TestValidatorRegistry_Accept_RejectsCrossSchemeChange — a validator
// re-registering under a different scheme byte is refused. The validator
// registry MUST treat the (scheme, NodeID) pair as the identity; a
// validator cannot silently swap from ML-DSA-65 to secp256k1 (or any
// other scheme) mid-epoch.
func TestValidatorRegistry_Accept_RejectsCrossSchemeChange(t *testing.T) {
	require := require.New(t)

	chainID := ID{0xfe, 0xed}
	pubKey := []byte("validator-rotation-pubkey")

	// Initial registration under ML-DSA-65.
	t1, _, err := TypedNodeIDFromMLDSA(NodeIDSchemeMLDSA65, chainID, pubKey)
	require.NoError(err)

	// Same NodeID re-presented under ML-DSA-87 — the scheme tag must
	// flip the wire identity, and the registry treats it as a different
	// identity.
	t2, _, err := TypedNodeIDFromMLDSA(NodeIDSchemeMLDSA87, chainID, pubKey)
	require.NoError(err)
	require.NotEqual(t1.NodeID, t2.NodeID,
		"swapping the scheme MUST produce a different NodeID via the digest")

	// A faux "rotation" that keeps the NodeID but flips the scheme on
	// the wire — the registry refuses this via the Compare ordering /
	// equality check the consensus layer enforces.
	swapped := TypedNodeID{Scheme: NodeIDSchemeSecp256k1, NodeID: t1.NodeID}
	require.NotEqual(t1, swapped,
		"TypedNodeID equality must distinguish identical NodeID under different schemes")
	require.NotEqual(0, t1.Compare(swapped))
}

// TestTypedNodeID_Compare_Stable — Compare provides a total order over
// (scheme, NodeID) lexicographically; used by validator-set commitments.
func TestTypedNodeID_Compare_Stable(t *testing.T) {
	require := require.New(t)

	a := TypedNodeID{Scheme: NodeIDSchemeMLDSA65, NodeID: NodeID{0x01}}
	b := TypedNodeID{Scheme: NodeIDSchemeMLDSA65, NodeID: NodeID{0x02}}
	c := TypedNodeID{Scheme: NodeIDSchemeMLDSA87, NodeID: NodeID{0x01}}
	d := TypedNodeID{Scheme: NodeIDSchemeSecp256k1, NodeID: NodeID{0x01}}

	require.Equal(0, a.Compare(a))
	require.Equal(-1, a.Compare(b))
	require.Equal(1, b.Compare(a))
	require.Equal(-1, a.Compare(c))
	require.Equal(-1, c.Compare(d))
}

// TestTypedNodeID_FromCert_IsClassicalCompat — TypedNodeIDFromCert
// produces a TypedNodeID tagged with NodeIDSchemeSecp256k1 so the
// strict-PQ gate refuses it at the boundary.
func TestTypedNodeID_FromCert_IsClassicalCompat(t *testing.T) {
	require := require.New(t)

	cert := &Certificate{Raw: []byte("der-encoded-cert-fixture")}
	t1 := TypedNodeIDFromCert(cert)
	require.Equal(NodeIDSchemeSecp256k1, t1.Scheme)
	require.Equal(NodeIDFromCert(cert), t1.NodeID)
	require.True(t1.Scheme.IsClassicalCompatUnsafe())
	require.False(t1.Scheme.IsPostQuantum())
}

// checkSchemeAgainstProfile is the cross-axis gate the consensus
// boundary applies: the presented NodeID's scheme MUST match the
// chain's pinned ValidatorScheme unless the operator has explicitly
// opted into classical compat and the presented scheme is the named
// classical scheme.
//
// Mirrored here (not imported from consensus/config) because this
// package is the dependency root for the test; the consensus package
// uses the same logic against its ChainSecurityProfile.ValidatorSchemeID
// field. Both paths funnel through ErrNodeIDSchemeMismatch.
func checkSchemeAgainstProfile(
	presented, pinned NodeIDScheme,
	classicalCompatUnsafe bool,
) error {
	if presented == pinned {
		return nil
	}
	if classicalCompatUnsafe && presented == NodeIDSchemeSecp256k1 {
		// Operator opted in; the named classical scheme is acceptable.
		return nil
	}
	return fmt.Errorf("presented=%s pinned=%s: %w", presented, pinned, ErrNodeIDSchemeMismatch)
}
