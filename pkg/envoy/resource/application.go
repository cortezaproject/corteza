package resource

import (
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	// Application represents a Application
	Application struct {
		*base
		Res *types.Application
	}
)

func NewApplication(res *types.Application) *Application {
	r := &Application{base: &base{}}
	r.SetResourceType(APPLICATION_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Name, res.ID)...)

	return r
}

func (r *Application) SysID() uint64 {
	return r.Res.ID
}
