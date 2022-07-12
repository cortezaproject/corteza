package rest

import (
	"context"
	"github.com/cortezaproject/corteza-server/discovery/rest/internal/mapping"
	"github.com/cortezaproject/corteza-server/discovery/rest/request"
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	mappings struct {
		settings *types.AppSettings

		sys interface {
			Users(context.Context) ([]*mapping.Mapping, error)
		}

		cmp interface {
			Namespaces(context.Context) ([]*mapping.Mapping, error)
			Modules(context.Context) ([]*mapping.Mapping, error)
			Records(context.Context) ([]*mapping.Mapping, error)
		}
	}
)

func Mappings() *mappings {
	return &mappings{
		settings: service.CurrentSettings,
		sys:      mapping.SystemMapping(),
		cmp:      mapping.ComposeMapping(),
	}
}

func (ctrl mappings) List(ctx context.Context, r *request.MappingsList) (interface{}, error) {
	var (
		out = make([]*mapping.Mapping, 0, 100)
		// Collection of all mapping functions we support
		//
		// Each function is responsible to
		mappingFn = make([]func(ctx context.Context) ([]*mapping.Mapping, error), 4)
	)

	if ctrl.settings.Discovery.SystemUsers.Enabled {
		mappingFn = append(mappingFn, ctrl.sys.Users)
	}
	if ctrl.settings.Discovery.ComposeNamespaces.Enabled {
		mappingFn = append(mappingFn, ctrl.cmp.Namespaces)
	}
	if ctrl.settings.Discovery.ComposeModules.Enabled {
		mappingFn = append(mappingFn, ctrl.cmp.Modules)
	}
	if ctrl.settings.Discovery.ComposeRecords.Enabled {
		mappingFn = append(mappingFn, ctrl.cmp.Records)
	}

	return out, func() error {
		for _, fn := range mappingFn {
			if fn == nil {
				continue
			}

			mm, err := fn(ctx)
			if err != nil {
				return err
			}

			out = append(out, mm...)
		}

		return nil
	}()
}
