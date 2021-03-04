package app

import (
	cr "crypto/rand"
	"encoding/binary"
	"math/rand"
)

func init() {
	// seed some randomness
	var b [8]byte
	if _, err := cr.Read(b[:]); err != nil {
		panic(err)
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}
