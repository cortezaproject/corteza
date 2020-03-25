package types

import (
	"io"
)

type (
	// JoinNodeEntry represents a { field: value } map for the given record in the
	// join node.
	JoinNodeEntry map[string]string
	// JoinNodeRecords represents a set of records that should be used together when
	// accessing the specified alias.
	JoinNodeRecords []JoinNodeEntry

	// JoinNode represents an ImportSource that will be used in combination with another,
	// to create a record based on multiple sources
	JoinNode struct {
		Mg   *ImportSource
		Name string
		// Records represents a set of records, available in the given import source.
		// [{ field: value }]
		Records []map[string]string
	}

	// ImportSource helps us perform some pre-proc operations before the actual import,
	// such as data mapping and source joining.
	ImportSource struct {
		Name string
		Path string

		Header *[]string
		Source io.Reader

		// DataMap allows us to specify what values from the original source should
		// map into what fields of what module
		DataMap io.Reader

		// SourceJoin allows us to specify what import sources should be joined
		// when mapping values.
		SourceJoin io.Reader
		// FieldMap stores records from the joined import source.
		// Records are indexed by {alias: [record]}
		FieldMap map[string]JoinNodeRecords
		// AliasMap helps us determine what fields are stored under the given alias.
		AliasMap map[string][]string

		// Value Map allows us to map specific values from the given import source into
		// a specified value used by Corteza.
		ValueMap map[string]map[string]string
	}
)
