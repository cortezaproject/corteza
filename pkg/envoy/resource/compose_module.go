package resource

import (
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza-server/compose/types"
)

type (
	ComposeModule struct {
		*base
		Res *types.Module
	}
)

func NewComposeModule(mod *types.Module) *ComposeModule {
	r := &ComposeModule{base: &base{}}
	r.SetResourceType("compose:module")
	r.Res = mod

	if mod.Handle != "" {
		r.AddIdentifier(mod.Handle)
	}
	if mod.Name != "" {
		r.AddIdentifier(mod.Name)
	}
	if mod.ID > 0 {
		r.AddIdentifier(strconv.FormatUint(mod.ID, 10))
	}

	// Field deps.
	for _, f := range mod.Fields {
		switch f.Kind {
		case "Record":
			refM := f.Options.String("module")
			if refM != "" && refM != "0" {
				r.AddRef("compose:module", refM)
			}
		}
	}

	return r
}

func (m *ComposeModule) SearchQuery() types.ModuleFilter {
	f := types.ModuleFilter{Query: ""}

	f.Handle = m.Res.Handle
	f.Name = m.Res.Name
	if m.Res.ID > 0 {
		f.Query = fmt.Sprintf("moduleID=%d", m.Res.ID)
	}

	return f
}
