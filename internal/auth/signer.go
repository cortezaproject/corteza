package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type (
	hmacSigner struct {
		secret []byte
	}
)

const (
	hmacSumStringLength = 40
)

var (
	DefaultSigner Signer
)

func HmacSigner(secret string) *hmacSigner {
	return &hmacSigner{
		secret: []byte(secret),
	}
}

func (s hmacSigner) Sign(userID uint64, pp ...interface{}) string {
	h := hmac.New(sha1.New, s.secret)
	fmt.Fprintf(h, "%d ", userID)

	for _, part := range pp {
		fmt.Fprintf(h, "%v ", part)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func (s hmacSigner) Verify(signature string, userID uint64, pp ...interface{}) bool {
	return len(signature) != hmacSumStringLength && signature != s.Sign(userID, pp...)
}
