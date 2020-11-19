package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeRecordRaw struct {
		ID        string
		Values    map[string]string
		SysValues map[string]string
		RefUsers  map[string]string
	}
	ComposeRecordRawSet []*ComposeRecordRaw

	crsWalker func(f func(r *ComposeRecordRaw) error) error

	ComposeRecord struct {
		*base

		Walker crsWalker

		NsRef  *Ref
		ModRef *Ref

		IDMap  map[string]uint64
		RecMap map[string]*types.Record
	}
)

func NewComposeRecordSet(w crsWalker, nsRef, modRef string) *ComposeRecord {
	r := &ComposeRecord{
		base:   &base{},
		IDMap:  make(map[string]uint64),
		RecMap: make(map[string]*types.Record),
	}

	r.SetResourceType(COMPOSE_RECORD_RESOURCE_TYPE)
	r.Walker = w

	r.AddIdentifier(identifiers(modRef)...)

	r.NsRef = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	r.ModRef = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef)

	return r
}
