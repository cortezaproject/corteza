package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/compose/types"
)

const (
	COMPOSE_RECORD_RESOURCE_TYPE = "ComposeRecordSet"
)

type (
	ComposeRecord struct {
		*base
		Res     *types.Record
		nsRef   string
		modRef  string
		userRef map[string]string
	}
)

func NewComposeRecord(res *types.Record, nsRef, modRef string, userRef map[string]string) *ComposeRecord {
	r := &ComposeRecord{base: &base{}}
	r.SetResourceType(COMPOSE_RECORD_RESOURCE_TYPE)
	r.Res = res
	r.nsRef = nsRef
	r.modRef = modRef

	r.AddIdentifier(identifiers(res.ID)...)

	r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)

	return r
}

func (m *ComposeRecord) SearchQuery() types.RecordFilter {
	f := types.RecordFilter{}

	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("recordID=%d", m.Res.ID)
	}

	return f
}
