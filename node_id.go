// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/luxfi/crypto/hash"
)

const (
	NodeIDPrefix = "NodeID-"
	NodeIDLen    = ShortIDLen
)

var (
	EmptyNodeID = NodeID{}

	errShortNodeID = errors.New("insufficient NodeID length")

	_ Sortable[NodeID] = NodeID{}
)

type NodeID ShortID

func (id NodeID) String() string {
	return ShortID(id).PrefixedString(NodeIDPrefix)
}

func (id NodeID) Bytes() []byte {
	return id[:]
}

func (id NodeID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

func (id NodeID) MarshalText() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *NodeID) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == nullStr { // If "null", do nothing
		return nil
	} else if len(str) <= 2+len(NodeIDPrefix) {
		return fmt.Errorf("%w: expected to be > %d", errShortNodeID, 2+len(NodeIDPrefix))
	}

	lastIndex := len(str) - 1
	if str[0] != '"' || str[lastIndex] != '"' {
		return errMissingQuotes
	}

	var err error
	*id, err = NodeIDFromString(str[1:lastIndex])
	return err
}

func (id *NodeID) UnmarshalText(text []byte) error {
	return id.UnmarshalJSON(text)
}

func (id NodeID) Compare(other NodeID) int {
	return bytes.Compare(id[:], other[:])
}

// ToNodeID attempt to convert a byte slice into a node id
func ToNodeID(bytes []byte) (NodeID, error) {
	nodeID, err := ToShortID(bytes)
	return NodeID(nodeID), err
}

func NodeIDFromCert(cert *Certificate) NodeID {
	return hash.ComputeHash160Array(
		hash.ComputeHash256(cert.Raw),
	)
}

// NodeIDFromString is the inverse of NodeID.String()
func NodeIDFromString(nodeIDStr string) (NodeID, error) {
	asShort, err := ShortFromPrefixedString(nodeIDStr, NodeIDPrefix)
	if err != nil {
		return NodeID{}, err
	}
	return NodeID(asShort), nil
}

// NodeIDPrefix for ML-DSA derived NodeIDs
const NodeIDMLDSADomainPrefix = "LuxNodeID/v1"

// NodeIDFromMLDSA derives a NodeID from an ML-DSA public key.
// This provides post-quantum secure node identity following FIPS 204.
// NodeID = H("LuxNodeID/v1" || mldsa_pubkey_bytes)
func NodeIDFromMLDSA(mldsaPubKey []byte) NodeID {
	prefix := []byte(NodeIDMLDSADomainPrefix)
	data := make([]byte, len(prefix)+len(mldsaPubKey))
	copy(data, prefix)
	copy(data[len(prefix):], mldsaPubKey)
	return hash.ComputeHash160Array(hash.ComputeHash256(data))
}
