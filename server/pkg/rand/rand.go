package rand

import (
	crand "crypto/rand"
	"math/big"
)

const (
	letterSpecials = "~=+%^*/()[]{}/!@#$?|"
	letterDigits   = "0123456789"
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + letterDigits
)

// Bytes returns n random letters and digits
//
// Output is used for tokens and secrets, so it comes from the OS random source
func Bytes(n int) []byte {
	b := make([]byte, n)

	for i := range b {
		b[i] = letterBytes[Intn(len(letterBytes))]
	}

	return b
}

// Password generates a random ASCII string with at least one digit and one special character
func Password(n int) string {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		var s string
		if i == 0 {
			s = letterDigits
		} else if i == 1 {
			s = letterSpecials
		} else {
			s = letterBytes
		}
		b[i] = s[Intn(len(s))]
	}

	for i := len(b) - 1; i > 0; i-- {
		j := Intn(i + 1)
		b[i], b[j] = b[j], b[i]
	}

	return string(b) // E.g. "3i[g0|)z"
}

// Intn returns a uniform random number in [0, n) from the OS random source
func Intn(n int) int {
	v, err := crand.Int(crand.Reader, big.NewInt(int64(n)))
	if err != nil {
		// no random source, nothing sensible can be done
		panic(err)
	}

	return int(v.Int64())
}
