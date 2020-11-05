package resource

import (
	"fmt"

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

func (r *Application) SearchQuery() types.ApplicationFilter {
	f := types.ApplicationFilter{
		Name: r.Res.Name,
	}

	if r.Res.ID > 0 {
		f.Query = fmt.Sprintf("applicationID=%d", r.Res.ID)
	}

	return f
}
