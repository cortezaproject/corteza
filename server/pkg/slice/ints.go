package slice

import (
	"bytes"
	"strconv"
)

type (
	// UInt64s represents slice of uint64 with some extra powers:
	//  - encodes (JSON) uint64s as strings
	UInt64s []uint64
)

func HasUint64(ss []uint64, s uint64) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}

	return false
}

func (uu UInt64s) MarshalJSON() ([]byte, error) {
	var (
		buf = bytes.Buffer{}
	)

	buf.WriteByte('[')
	for i, u := range uu {
		if i > 0 {
			buf.WriteByte(',')
		}

		buf.WriteByte('"')
		buf.WriteString(strconv.FormatUint(u, 10))
		buf.WriteByte('"')
	}
	buf.WriteByte(']')

	return buf.Bytes(), nil
}
