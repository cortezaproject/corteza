package types

import (
	"io"
)

var (
	SfDateTime = "2006-01-02 15:04:05"
)

type (
	JoinedNodeEntry   map[string]string
	JoinedNodeRecords []JoinedNodeEntry

	Migrateable struct {
		Name string
		Path string

		Header *[]string

		Source io.Reader
		// map is used for stream splitting
		Map io.Reader

		// join is used for source joining
		Join  io.Reader
		Joins []*JoinedNode
		// alias.ID: [value]
		FieldMap map[string]JoinedNodeRecords
		// helps us determine what value field to use for linking
		AliasMap map[string]string

		// value is used for field value mapping
		// field: value from: value to
		ValueMap map[string]map[string]string
	}
)
