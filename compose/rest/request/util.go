package request

//lint:file-ignore U1000 Ignore unused code, part of request pkg toolset

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx/types"
	"github.com/pkg/errors"
)

type (
	ProcedureArgs []ProcedureArg

	ProcedureArg struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}
)

var truthy = regexp.MustCompile(`^\s*(t(rue)?|y(es)?|1)\s*$`)

func parseJSONTextWithErr(s string) (types.JSONText, error) {
	result := &types.JSONText{}
	err := errors.Wrap(result.Scan(s), "error when parsing JSONText")
	return *result, err
}

func parseISODateWithErr(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func parseISODatePtrWithErr(s string) (*time.Time, error) {
	t, err := parseISODateWithErr(s)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// parseInt parses a string to int
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	i, _ := strconv.Atoi(s)
	return i
}

// parseInt parses a string to int
func parseUint(s string) uint {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseUint(s, 10, 32)
	return uint(i)
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

func parseUInt64A(values []string) []uint64 {
	var result []uint64
	if len(values) > 0 {
		for _, val := range values {
			result = append(result, parseUInt64(val))
		}
	}
	return result
}

func parseStrings(values []string) []string {
	return values
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

func (args ProcedureArgs) GetUint64(name string) uint64 {
	u, _ := strconv.ParseUint(args.Get(name), 10, 64)
	return u
}

func (args ProcedureArgs) Get(name string) string {
	name = strings.ToLower(name)
	for _, arg := range args {
		if strings.ToLower(arg.Name) == name {
			return arg.Value
		}
	}

	return ""
}
