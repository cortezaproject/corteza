package resource

import (
	"fmt"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
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
		RelMod *composeTypes.Module

		IDMap map[string]uint64
		// UserFlakes help the system by predefining a set of potential sys user references.
		// This should make the operation cheaper for larger datasets.
		UserFlakes UserstampIndex
	}
)

func NewComposeRecordSet(w CrsWalker, nsRef, modRef string) *ComposeRecord {
	r := &ComposeRecord{
		base:  &base{},
		IDMap: make(map[string]uint64),
	}

	r.SetResourceType(composeTypes.RecordResourceType)
	r.Walker = w

	r.AddIdentifier(identifiers(modRef)...)

	r.RefNs = r.AddRef(composeTypes.NamespaceResourceType, nsRef)
	r.RefMod = r.AddRef(composeTypes.ModuleResourceType, modRef).Constraint(r.RefNs)

	return r
}

func (r *ComposeRecord) SetUserFlakes(uu UserstampIndex) {
	r.UserFlakes = uu

	// Set user refs as wildflag, indicating it refers to any user resource
	r.AddRef(systemTypes.UserResourceType, "*")
}

func (r *ComposeRecord) RBACPath() []*Ref {
	return []*Ref{r.RefNs, r.RefMod}
}

func FindComposeRecordResource(rr InterfaceSet, ii Identifiers) (rec *ComposeRecord) {
	rr.Walk(func(r Interface) error {
		crr, ok := r.(*ComposeRecord)
		if !ok {
			return nil
		}

		if crr.Identifiers().HasAny(ii) {
			rec = crr
		}
		return nil
	})

	return rec
}

func ComposeRecordErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose record unresolved %v", ii.StringSlice())
}
