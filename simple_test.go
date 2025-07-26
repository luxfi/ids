package ids

import "testing"

func TestSimple(t *testing.T) {
	// Test basic ID creation
	id := ID{1, 2, 3}
	if len(id) != IDLen {
		t.Fatalf("Expected ID length %d, got %d", IDLen, len(id))
	}

	// Test ShortID creation
	shortID := ShortID{1, 2, 3}
	if len(shortID) != ShortIDLen {
		t.Fatalf("Expected ShortID length %d, got %d", ShortIDLen, len(shortID))
	}
}
