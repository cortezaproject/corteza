package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// Template represents a Template
	Template struct {
		*base
		Res *types.Template
	}
)

func NewTemplate(t *types.Template) *Template {
	r := &Template{base: &base{}}
	r.SetResourceType(types.TemplateResourceType)
	r.Res = t

	r.AddIdentifier(identifiers(t.Handle, t.Meta.Short, t.ID)...)

	// Initial timestamps
	r.SetTimestamps(MakeTimestampsCUDAS(&t.CreatedAt, t.UpdatedAt, t.DeletedAt, nil, nil))

	return r
}

func (r *Template) SysID() uint64 {
	return r.Res.ID
}

// FindTemplate looks for the template in the resources
func FindTemplate(rr InterfaceSet, ii Identifiers) (u *types.Template) {
	var tRes *Template

	rr.Walk(func(r Interface) error {
		tr, ok := r.(*Template)
		if !ok {
			return nil
		}

		if tr.Identifiers().HasAny(ii) {
			tRes = tr
		}
		return nil
	})

	// Found it
	if tRes != nil {
		return tRes.Res
	}

	return nil
}

func TemplateErrUnresolved(ii Identifiers) error {
	return fmt.Errorf("template unresolved %v", ii.StringSlice())
}
