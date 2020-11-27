package resource

import (
	"strings"

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
		ID        string
		Values    map[string]string
		SysValues *rawSysValues
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

// ApplyValues takes in a raw map of things and creates a proper structure for it
func (cr *ComposeRecordRaw) ApplyValues(vv map[string]string) {
	if cr.RefUsers == nil {
		cr.RefUsers = make(map[string]string)
	}
	if cr.SysValues == nil {
		cr.SysValues = &rawSysValues{}
	}
	if cr.Values == nil {
		cr.Values = make(map[string]string)
	}

	for k, v := range vv {
		switch strings.ToLower(k) {
		case "ownedby":
			cr.SysValues.OwnedBy = v
			cr.RefUsers[v] = ""
		case "createdat":
			cr.SysValues.CreatedAt = v
		case "createdby":
			cr.SysValues.CreatedBy = v
			cr.RefUsers[v] = ""
		case "updatedat":
			cr.SysValues.UpdatedAt = v
		case "updatedby":
			cr.SysValues.UpdatedBy = v
			cr.RefUsers[v] = ""
		case "deletedat":
			cr.SysValues.DeletedAt = v
		case "deletedby":
			cr.SysValues.DeletedBy = v
			cr.RefUsers[v] = ""

		default:
			cr.Values[k] = v
		}
	}
}
