package decoder

import (
	"github.com/crusttech/human/server/compose/types"
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
