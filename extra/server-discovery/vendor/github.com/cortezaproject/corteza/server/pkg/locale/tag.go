package locale

import (
	"errors"

	"golang.org/x/text/language"
)

const (
	// longest realistic BCP 47 tag is well under this
	maxTagLength = 64

	// browsers send well under 200 bytes, leave room for proxies that append
	maxAcceptLanguageLength = 1024
)

var (
	ErrTagTooLong = errors.New("language tag too long")
)

// ParseTag parses a single language tag from untrusted input
//
// Length is checked first, parsing time grows quadratically with the input
func ParseTag(s string) (language.Tag, error) {
	if len(s) > maxTagLength {
		return language.Und, ErrTagTooLong
	}

	return language.Parse(s)
}

// MakeTag is ParseTag that falls back to und
func MakeTag(s string) language.Tag {
	if len(s) > maxTagLength {
		return language.Und
	}

	return language.Make(s)
}

// ParseAcceptLanguage parses Accept-Language style list from untrusted input
func ParseAcceptLanguage(s string) ([]language.Tag, []float32, error) {
	if len(s) > maxAcceptLanguageLength {
		return nil, nil, ErrTagTooLong
	}

	return language.ParseAcceptLanguage(s)
}
