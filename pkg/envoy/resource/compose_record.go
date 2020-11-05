package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeRecordRaw struct {
		ID        string
		Values    map[string]string
		SysValues map[string]string
		RefUser   map[string]string
	}
	ComposeRecordRawSet []*ComposeRecordRaw

	crsWalker func(f func(r *ComposeRecordRaw) error) error

	ComposeRecord struct {
		*base
		// Res     *types.Record

		Walker crsWalker

		NsRef     *Ref
		ModRef    *Ref
		ModFields types.ModuleFieldSet
		UserRef   map[string]string
	}
)

func NewComposeRecordSet(w crsWalker, nsRef, modRef string) *ComposeRecord {
	r := &ComposeRecord{base: &base{}}
	r.SetResourceType(COMPOSE_RECORD_RESOURCE_TYPE)
	r.Walker = w

	r.AddIdentifier(identifiers(modRef)...)

	// for _, u := range userRef {
	// 	r.AddRef(USER_RESOURCE_TYPE, u)
	// }

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	r.ModRef = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)

	return r
}

func (m *ComposeRecord) SearchQuery() types.RecordFilter {
	f := types.RecordFilter{}

	return f
}
