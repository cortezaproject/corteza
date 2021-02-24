package resource

import (
	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	rawSysValues struct {
		OwnedBy   string
		CreatedAt string
		CreatedBy string
		UpdatedAt string
		UpdatedBy string
		DeletedAt string
		DeletedBy string
	}
	ComposeRecordRaw struct {
		ID     string
		Values map[string]string

		Ts *Timestamps
		Us *Userstamps

		Config *EnvoyConfig
	}
	ComposeRecordRawSet []*ComposeRecordRaw

	CrsWalker func(f func(r *ComposeRecordRaw) error) error

	ComposeRecord struct {
		*base

		Walker CrsWalker

		RefNs *Ref

		RefMod *Ref
		RelMod *types.Module

		IDMap  map[string]uint64
		RecMap map[string]*types.Record
		// UserFlakes help the system by predefining a set of potential sys user references.
		// This should make the operation cheaper for larger datasets.
		UserFlakes UserstampIndex
	}
)

func NewComposeRecordSet(w CrsWalker, nsRef, modRef string) *ComposeRecord {
	r := &ComposeRecord{
		base:   &base{},
		IDMap:  make(map[string]uint64),
		RecMap: make(map[string]*types.Record),
	}

	r.SetResourceType(COMPOSE_RECORD_RESOURCE_TYPE)
	r.Walker = w

	r.AddIdentifier(identifiers(modRef)...)

	r.RefNs = r.AddRef(COMPOSE_NAMESPACE_RESOURCE_TYPE, nsRef)
	r.RefMod = r.AddRef(COMPOSE_MODULE_RESOURCE_TYPE, modRef).Constraint(r.RefNs)

	return r
}

func (r *ComposeRecord) SetUserFlakes(uu UserstampIndex) {
	r.UserFlakes = uu

	// Set user refs as wildflag, indicating it refers to any user resource
	r.AddRef(USER_RESOURCE_TYPE, "*")
}
