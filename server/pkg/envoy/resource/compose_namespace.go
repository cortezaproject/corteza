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
	r.SetResourceType(types.NamespaceResourceType)
	r.Res = ns

	r.AddIdentifier(identifiers(ns.Slug, ns.Name, ns.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDA(&ns.CreatedAt, ns.UpdatedAt, ns.DeletedAt, nil))

	return r
}

func (r *ComposeNamespace) Resource() interface{} {
	return r.Res
}

func (r *ComposeNamespace) RBACParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.NamespaceRbacResourceTpl(), types.NamespaceResourceType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Slug))

	return
}

func (r *ComposeNamespace) ResourceTranslationParts() (resource string, ref *Ref, path []*Ref) {
	ref = r.Ref()
	path = nil
	resource = fmt.Sprintf(types.NamespaceResourceTranslationTpl(), types.NamespaceResourceTranslationType, firstOkString(strconv.FormatUint(r.Res.ID, 10), r.Res.Slug))

	return
}

func (r *ComposeNamespace) SysID() uint64 {
	return r.Res.ID
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
