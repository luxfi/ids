package ids

// EXPERIMENT — see blue report. This states a ZAP wire for the fixed-width id
// types: the raw bytes, at their declared width, with no prefix and no text
// spelling. The width is the type's own length, so it cannot drift from the
// declaration.
//
// It is kept on a branch and NOT merged: an id is an INLINE field of the message
// that carries it, not a message of its own, so a parent codec writes it with
// SetBytesFixed and never calls these. Proven in the report.

// MarshalZAP returns the id's raw bytes.
func (id ID) MarshalZAP() ([]byte, error) {
	out := make([]byte, len(id))
	copy(out, id[:])
	return out, nil
}

// UnmarshalZAP reads the id's raw bytes.
func (id *ID) UnmarshalZAP(b []byte) error {
	if len(b) != len(id) {
		return errWrongZAPLen
	}
	copy(id[:], b)
	return nil
}

func (id ShortID) MarshalZAP() ([]byte, error) {
	out := make([]byte, len(id))
	copy(out, id[:])
	return out, nil
}

func (id *ShortID) UnmarshalZAP(b []byte) error {
	if len(b) != len(id) {
		return errWrongZAPLen
	}
	copy(id[:], b)
	return nil
}

func (id NodeID) MarshalZAP() ([]byte, error) {
	out := make([]byte, len(id))
	copy(out, id[:])
	return out, nil
}

func (id *NodeID) UnmarshalZAP(b []byte) error {
	if len(b) != len(id) {
		return errWrongZAPLen
	}
	copy(id[:], b)
	return nil
}

var errWrongZAPLen = zapLenError("ids: ZAP wire length mismatch")

type zapLenError string

func (e zapLenError) Error() string { return string(e) }
