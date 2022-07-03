package service

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/types"
	sysService "github.com/cortezaproject/corteza-server/system/service"
	sysTypes "github.com/cortezaproject/corteza-server/system/types"
)

type (
	dataPrivacy struct {
		ns      NamespaceService
		m       ModuleService
		dalConn dalConnectionService
	}

	dalConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (q *sysTypes.DalConnection, err error)
	}

	moduleSetPayload struct {
		Set []types.PrivacyModule `json:"set"`
	}

	DataPrivacyService interface {
		FindModules(ctx context.Context, filter types.PrivacyModuleFilter) (types.PrivacyModuleSet, types.PrivacyModuleFilter, error)
	}
)

func DataPrivacy() *dataPrivacy {
	return &dataPrivacy{
		ns:      DefaultNamespace,
		m:       DefaultModule,
		dalConn: sysService.DefaultDalConnection,
	}
}

func (svc dataPrivacy) FindModules(ctx context.Context, filter types.PrivacyModuleFilter) (out types.PrivacyModuleSet, f types.PrivacyModuleFilter, err error) {
	var (
		modules []types.PrivacyModule
		cc      = make(map[uint64]*sysTypes.DalConnection, 0)
	)

	reqConnes := make(map[uint64]bool)
	hasReqConnes := len(filter.ConnectionID) > 0
	for _, connectionID := range filter.ConnectionID {
		reqConnes[connectionID] = true
	}

	// All namespaces
	namespaces, _, err := svc.ns.Find(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	for _, n := range namespaces {
		// Sensitive modules only
		modules, f, err = svc.m.FindSensitive(ctx, types.PrivacyModuleFilter{NamespaceID: n.ID})
		if err != nil {
			return
		}
		if len(modules) == 0 {
			continue
		}

		for _, m := range modules {
			connID := m.ConnectionID
			if hasReqConnes && !reqConnes[connID] {
				continue
			}

			var c *sysTypes.DalConnection
			if val, ok := cc[connID]; ok {
				c = val
			} else {
				c, err = svc.dalConn.FindByID(ctx, connID)
				if err != nil {
					cc[connID] = c
				}
			}

			out = append(out, &types.PrivacyModule{
				ID:         m.ID,
				Name:       m.Name,
				Handle:     m.Handle,
				Owner:      m.Owner,
				Connection: c,
			})
		}
	}

	return out, f, nil
}
