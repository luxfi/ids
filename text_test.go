package ids

import (
	"encoding/json"
	"testing"
)

// A TextUnmarshaler is handed the word ITSELF, never a quoted one — that is the
// contract encoding/json keeps for a map key, encoding/xml for an attribute,
// flag.Value for an argument, and a URL binder for a query value. NodeID and
// ShortID delegated to UnmarshalJSON, so all four met "first and last
// characters should be quotes" for a value that was never JSON, and a node id
// named in a URL arrived as the zero node with a 200 on the answer.
func TestTextIsNotQuotedJSON(t *testing.T) {
	node := GenerateTestNodeID()
	var gotNode NodeID
	if err := gotNode.UnmarshalText([]byte(node.String())); err != nil || gotNode != node {
		t.Errorf("NodeID: got %v, err %v, want %v", gotNode, err, node)
	}

	short := GenerateTestShortID()
	var gotShort ShortID
	if err := gotShort.UnmarshalText([]byte(short.String())); err != nil || gotShort != short {
		t.Errorf("ShortID: got %v, err %v, want %v", gotShort, err, short)
	}

	id := GenerateTestID()
	var gotID ID
	if err := gotID.UnmarshalText([]byte(id.String())); err != nil || gotID != id {
		t.Errorf("ID: got %v, err %v, want %v", gotID, err, id)
	}
}

// And the JSON spelling is unchanged: it is the same word in quotes, read
// through the same one reading.
func TestJSONIsTheTextInQuotes(t *testing.T) {
	node := GenerateTestNodeID()
	raw, err := json.Marshal(node)
	if err != nil {
		t.Fatal(err)
	}
	if want := `"` + node.String() + `"`; string(raw) != want {
		t.Errorf("marshal: got %s, want %s", raw, want)
	}
	var back NodeID
	if err := json.Unmarshal(raw, &back); err != nil || back != node {
		t.Errorf("round trip: got %v, err %v", back, err)
	}
	if err := json.Unmarshal([]byte(`NodeID-unquoted`), &back); err == nil {
		t.Error("JSON accepted an unquoted word")
	}

	// A map keyed by a node id decodes, which is the case that made ID's own
	// correction necessary: the stdlib hands a KEY to UnmarshalText.
	var keyed map[NodeID]int
	if err := json.Unmarshal([]byte(`{"`+node.String()+`":1}`), &keyed); err != nil || keyed[node] != 1 {
		t.Errorf("map key: got %v, err %v", keyed, err)
	}
}
