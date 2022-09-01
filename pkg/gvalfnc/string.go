package gvalfnc

import (
	"strings"

	"github.com/spf13/cast"
)

// @todo
func ConcatStrings(parts ...any) (string, error) {
	pp, err := cast.ToStringSliceE(parts)
	if err != nil {
		return "", err
	}

	return strings.Join(pp, ""), nil
}
