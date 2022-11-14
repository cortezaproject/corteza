package decoder

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeRecord struct {
		types.Record
	}
	ComposeRecordSet []*ComposeRecord

	ComposeRecordFilter struct {
		types.RecordFilter
	}
)
