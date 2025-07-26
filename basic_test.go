// Copyright (C) 2020-2025, Lux Industries Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ids

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicIDFunctionality(t *testing.T) {
	require := require.New(t)

	// Test ID creation and string conversion
	id := ID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
	str := id.String()
	require.NotEmpty(str)

	// Test ID from string
	id2, err := FromString(str)
	require.NoError(err)
	require.Equal(id, id2)

	// Test ShortID
	shortID := ShortID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	shortStr := shortID.String()
	require.NotEmpty(shortStr)

	// Test NodeID
	nodeID := NodeID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	nodeStr := nodeID.String()
	require.NotEmpty(nodeStr)
	require.Contains(nodeStr, NodeIDPrefix)
}
