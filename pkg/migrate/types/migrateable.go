package types

import (
	"io"
)

var (
	SfDateTime = "2006-01-02 15:04:05"
)

type (
	Migrateable struct {
		Name string
		Path string

		Source io.ReadSeeker
	}
)
