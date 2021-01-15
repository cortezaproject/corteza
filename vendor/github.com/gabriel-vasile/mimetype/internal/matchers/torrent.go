package matchers

import (
	"bytes"
)

// Torrent has bencoded text in the beginning
func Torrent(in []byte) bool {
	return len(in) > 11 &&
		bytes.Equal(in[:11], []byte("d8:announce"))
}
