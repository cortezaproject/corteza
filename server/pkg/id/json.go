package id

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	// Uint64 is a uint64 identifier that is written to JSON as a string and
	// read from either a string or a number.
	Uint64 uint64

	// Uint64s is a list of uint64 identifiers that is written to JSON as an
	// array of strings and read from an array of strings or numbers.
	Uint64s []uint64
)

func (v Uint64) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatUint(uint64(v), 10) + `"`), nil
}

func (v *Uint64) UnmarshalJSON(data []byte) error {
	n, err := parseJSONUint(data)
	if err != nil {
		return err
	}

	*v = Uint64(n)
	return nil
}

func (vv Uint64s) MarshalJSON() ([]byte, error) {
	if vv == nil {
		return []byte("null"), nil
	}

	return json.Marshal(Strings(vv...))
}

func (vv *Uint64s) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if raw == nil {
		*vv = nil
		return nil
	}

	out := make(Uint64s, len(raw))
	for i, r := range raw {
		n, err := parseJSONUint(r)
		if err != nil {
			return err
		}
		out[i] = n
	}

	*vv = out
	return nil
}

// parseJSONUint reads a JSON string or number holding a uint64; null and ""
// are zero.
func parseJSONUint(data []byte) (uint64, error) {
	data = bytes.TrimSpace(data)
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	if len(data) == 0 || string(data) == "null" {
		return 0, nil
	}

	n, err := strconv.ParseUint(string(data), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid ID %s: %w", data, err)
	}

	return n, nil
}
