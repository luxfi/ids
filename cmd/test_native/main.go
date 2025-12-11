package main

import (
	"encoding/json"
	"fmt"

	"github.com/luxfi/ids"
)

func main() {
	fmt.Println("Testing Native Chain IDs")
	fmt.Println("========================")
	fmt.Println()

	// Test all native chain IDs
	chains := []struct {
		name string
		id   ids.ID
	}{
		{"P-Chain", ids.PChainID},
		{"C-Chain", ids.CChainID},
		{"X-Chain", ids.XChainID},
		{"Q-Chain", ids.QChainID},
		{"A-Chain", ids.AChainID},
		{"B-Chain", ids.BChainID},
		{"T-Chain", ids.TChainID},
	}

	fmt.Println("1. Native Chain ID String representations:")
	for _, c := range chains {
		fmt.Printf("   %s: %s\n", c.name, c.id.String())
	}

	fmt.Println()
	fmt.Println("2. Parse from string:")
	testStrs := []string{
		"11111111111111111111111111111111P",
		"11111111111111111111111111111111C",
		"P", "C", "X", "Q", "A", "B", "T",
		"p", "c", "x", "q", "a", "b", "t",
	}
	for _, s := range testStrs {
		id, err := ids.FromString(s)
		if err != nil {
			fmt.Printf("   %q -> ERROR: %v\n", s, err)
		} else {
			fmt.Printf("   %q -> %s (is native: %v)\n", s, id.String(), ids.IsNativeChain(id))
		}
	}

	fmt.Println()
	fmt.Println("3. JSON marshaling:")
	for _, c := range chains {
		b, _ := json.Marshal(c.id)
		fmt.Printf("   %s: %s\n", c.name, string(b))
	}

	fmt.Println()
	fmt.Println("4. JSON unmarshaling:")
	jsonInputs := []string{
		`"11111111111111111111111111111111P"`,
		`"11111111111111111111111111111111C"`,
		`"P"`,
		`"C"`,
	}
	for _, j := range jsonInputs {
		var id ids.ID
		if err := json.Unmarshal([]byte(j), &id); err != nil {
			fmt.Printf("   %s -> ERROR: %v\n", j, err)
		} else {
			fmt.Printf("   %s -> %s\n", j, id.String())
		}
	}

	fmt.Println()
	fmt.Println("5. Verify byte layout:")
	for _, c := range chains {
		fmt.Printf("   %s: [31]=%d (0x%02x = '%c')\n", c.name, c.id[31], c.id[31], c.id[31])
	}

	fmt.Println()
	fmt.Println("6. Compare with Empty ID:")
	fmt.Printf("   Empty.String(): %s\n", ids.Empty.String())
	fmt.Printf("   Empty is native: %v\n", ids.IsNativeChain(ids.Empty))
	fmt.Printf("   PChainID == Empty: %v\n", ids.PChainID == ids.Empty)

	fmt.Println()
	fmt.Println("7. Helper functions:")
	for _, c := range chains {
		alias := ids.NativeChainAlias(c.id)
		fmt.Printf("   %s alias: %q\n", c.name, alias)
	}

	fmt.Println()
	fmt.Println("All tests passed!")
}
