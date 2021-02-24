package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeNamespace struct {
		*base
		Res *types.Namespace
	}
)

func NewComposeNamespace(ns *types.Namespace) *ComposeNamespace {
	r := &ComposeNamespace{base: &base{}}
	r.SetResourceType(COMPOSE_NAMESPACE_RESOURCE_TYPE)
	r.Res = ns

	r.AddIdentifier(identifiers(ns.Slug, ns.Name, ns.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeCUDATimestamps(&ns.CreatedAt, ns.UpdatedAt, ns.DeletedAt, nil))

	return r
}

func (r *ComposeNamespace) SysID() uint64 {
	return r.Res.ID
}

func (r *ComposeNamespace) Ref() string {
	return FirstOkString(r.Res.Slug, r.Res.Name, strconv.FormatUint(r.Res.ID, 10))
}

// FindComposeNamespace looks for the namespace in the resource set
func FindComposeNamespace(rr InterfaceSet, ii Identifiers) (ns *types.Namespace) {
	var nsRes *ComposeNamespace

	rr.Walk(func(r Interface) error {
		nr, ok := r.(*ComposeNamespace)
		if !ok {
			return nil
		}

		if nr.Identifiers().HasAny(ii) {
			nsRes = nr
		}
		return nil
	})

	// Found it
	if nsRes != nil {
		return nsRes.Res
	}

	return nil
}

func ComposeNamespaceErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("compose namespace unresolved %v", ii.StringSlice())
}
