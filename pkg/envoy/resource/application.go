package resource

import (
	"fmt"

	"github.com/cortezaproject/corteza-server/system/types"
)

const (
	APPLICATION_RESOURCE_TYPE = "application"
)

type (
	// application represents a Application
	application struct {
		*base
		Res *types.Application
	}
)

func Application(res *types.Application) *application {
	r := &application{base: &base{}}
	r.SetResourceType(APPLICATION_RESOURCE_TYPE)
	r.Res = res

	r.AddIdentifier(identifiers(res.Name, res.ID)...)

	return r
}

func (r *application) SearchQuery() types.ApplicationFilter {
	f := types.ApplicationFilter{
		Name: r.Res.Name,
	}

	if r.Res.ID > 0 {
		f.Query = fmt.Sprintf("applicationID=%d", r.Res.ID)
	}

	return f
}
