package str

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/handle"
	"strings"
)

const (
	// DefaultLevenshteinDistance is the default levenshtein distance
	DefaultLevenshteinDistance = 3

	CaseInSensitiveMatch = iota
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
		return strings.Compare(str1, str2) == 0
	case CaseInSensitiveMatch:
		return strings.EqualFold(str1, str2)
	default:
		return false
	}
}

func ParseStrings(ss []string) (m map[string]string, err error) {
	if len(ss) == 0 {
		return nil, nil
	}

	m = make(map[string]string)

	for _, s := range ss {
		if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
			// assume json
			if err = json.Unmarshal([]byte(s), &m); err != nil {
				return nil, err
			}

			continue
		}

		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid label format")
		}

		if !handle.IsValid(kv[0]) {
			return nil, fmt.Errorf("invalid label key format")
		}

		m[kv[0]] = kv[1]
	}

	return m, nil
}
