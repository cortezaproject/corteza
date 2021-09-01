package rand

import (
	"math/rand"
	"sync"
	"time"
)

// credits: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang

const (
	letterSpecials = "~=+%^*/()[]{}/!@#$?|"
	letterDigits   = "0123456789"
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" + letterDigits
	letterIdxBits  = 6                    // 6 bits to represent a letter index
	letterIdxMask  = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax   = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randSrc = rand.NewSource(time.Now().UnixNano())
var mu sync.Mutex

func Bytes(n int) []byte {
	mu.Lock()
	defer mu.Unlock()

	b := make([]byte, n)

	// A randSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// Password generates a random ASCII string with at least one digit and one special character
func Password(n int) string {
	mu.Lock()
	defer mu.Unlock()

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
		b[i] = s[rand.Intn(len(s))]
	}
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})
	return string(b) // E.g. "3i[g0|)z"
}
