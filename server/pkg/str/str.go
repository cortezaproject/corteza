package str

import (
	"strings"
)

const (
	// DefaultLevenshteinDistance is the default levenshtein distance
	DefaultLevenshteinDistance = 3

	SimpleMatch = iota
	CaseSensitiveMatch
	LevenshteinDistance
	Soundex
)

// Match will match string as per given algorithm
func Match(str1, str2 string, algorithm int) bool {
	switch algorithm {
	case LevenshteinDistance:
		return ToLevenshteinDistance(str1, str2) <= DefaultLevenshteinDistance
	case Soundex:
		return ToSoundex(str1) == ToSoundex(str2)
	case CaseSensitiveMatch:
		return str1 == str2
	case SimpleMatch:
		return ToLower(str1) == ToLower(str2)
	default:
		return false
	}
}

func ToLower(s string) string {
	return strings.ToLower(s)
}
