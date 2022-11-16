package service

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	sysService "github.com/cortezaproject/corteza/server/system/service"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	dataPrivacy struct {
		ns      NamespaceService
		m       ModuleService
		dalConn dalConnectionService
		locale  ResourceTranslationsManagerService
	}

	dalConnectionService interface {
		FindByID(ctx context.Context, ID uint64) (q *sysTypes.DalConnection, err error)
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
		locale:  DefaultResourceTranslation,
	}
}

func (svc dataPrivacy) FindModules(ctx context.Context, filter types.PrivacyModuleFilter) (out types.PrivacyModuleSet, f types.PrivacyModuleFilter, err error) {
	var (
		modules []types.PrivacyModule
		cc      = make(map[uint64]*sysTypes.DalConnection, 0)
	)

	namespaces, _, err := svc.ns.Find(ctx, types.NamespaceFilter{})
	if err != nil {
		return
	}

	for _, n := range namespaces {
		tag := locale.GetAcceptLanguageFromContext(ctx)
		n.DecodeTranslations(svc.locale.Locale().ResourceTranslations(tag, n.ResourceTranslation()))

		modules, f, err = svc.m.SearchSensitive(ctx, types.PrivacyModuleFilter{
			NamespaceID:  n.ID,
			ConnectionID: filter.ConnectionID,
		})
		if err != nil {
			return
		}
		if len(modules) == 0 {
			continue
		}

		for _, m := range modules {
			connID := m.ConnectionID
			if val, ok := cc[connID]; ok {
				m.Connection = val
			} else {
				m.Connection, err = svc.dalConn.FindByID(ctx, connID)
				if err != nil {
					cc[connID] = m.Connection
				}
			}

			m.Namespace = types.PrivacyNamespaceMeta{
				ID:   n.ID,
				Slug: n.Slug,
				Name: n.Name,
			}
			out = append(out, &m)
		}
	}

	return out, f, nil
}
