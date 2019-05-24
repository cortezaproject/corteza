package request

//lint:file-ignore U1000 Ignore unused code, part of request pkg toolset

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
)

var truthy = regexp.MustCompile(`^\s*(t(rue)?|y(es)?|1)\s*$`)

func parseJSONText(s string) (types.JSONText, error) {
	result := &types.JSONText{}
	err := errors.Wrap(result.Scan(s), "error when parsing JSONText")
	return *result, err
}

// parseInt parses a string to int
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parseInt64 parses a string to int64
func parseInt64(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// parseUInt64 parses a string to uint64
func parseUInt64(s string) uint64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 64)
	return i
}

// parseUInt64 parses a string to uint64
func parseUint(s string) uint {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
}

func parseUInt64A(values []string) []uint64 {
	var result []uint64
	if len(values) > 0 {
		for _, val := range values {
			result = append(result, parseUInt64(val))
		}
	}
	return result
}

// parseUInt64 parses a string to uint64
func parseBool(s string) bool {
	return truthy.MatchString(strings.ToLower(s))
}

// is checks if string s is contained in matches
func is(s string, matches ...string) bool {
	for _, v := range matches {
		if s == v {
			return true
		}
	}
	return false
}
