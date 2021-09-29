package provision

import (
	"context"
	"strconv"
	"strings"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

type (
	resourceIndex struct {
		fields         map[uint64]*composeTypes.ModuleField
		modules        map[uint64]*composeTypes.Module
		charts         map[uint64]*composeTypes.Chart
		pages          map[uint64]*composeTypes.Page
		exposedModules map[uint64]*federationTypes.ExposedModule
		sharedModules  map[uint64]*federationTypes.SharedModule
	}
)

// MigrateOperations creates system roles
func migratePre202109RbacRules(ctx context.Context, log *zap.Logger, s store.Storer) error {
	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		if err != nil {
			return err
		}

		log.Info("migrating RBAC rules to new format", zap.Int("rules", len(rr)))

		rx, err := preloadResourceIndex(ctx, s)
		if err != nil {
			return err
		}

		for _, r := range rr {
			var (
				cr     = *r
				action = migratePre202109RbacRule(r, rx)
			)

			if action != 0 {
				err = store.DeleteRbacRule(ctx, s, &cr)
				if err != nil {
					return err
				}
			}

			if action == 1 {
				err = store.CreateRbacRule(ctx, s, r)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

//  0 - no action
// -1 - remove
//  1 - update
func migratePre202109RbacRule(r *rbac.Rule, rx *resourceIndex) (op int) {
	const (
		nsSep = "::"
		nsDef = "corteza"
	)

	if strings.Contains(r.Resource, nsSep) {
		return
	}

	if r.Resource == "" {
		return -1
	}

	// split old format
	parts := strings.SplitN(r.Resource, ":", 3)

	switch {
	case parts[0] == "messaging":
		return -1
	case len(parts) > 1 && parts[1] == "automation-script":
		return -1
	case len(parts) == 1:
		r.Resource = nsDef + nsSep + strings.Join(parts[:1], ":")
	default:
		r.Resource = nsDef + nsSep + strings.Join(parts[:2], ":")
	}

	op = 1

	rType := rbac.ResourceType(r.Resource)
	switch {
	case rType == systemTypes.UserResourceType && strings.HasPrefix(r.Operation, "unmask."):
		// flipping terms in user unmask operations
		r.Operation = strings.TrimPrefix(r.Operation, "unmask.") + ".unmask"

	case rType == "corteza::federation:module":
		// fed. module resource was split into two resources - exposed & shared
		if r.Operation == "manage" {
			rType = federationTypes.ExposedModuleResourceType
		} else if r.Operation == "map" {
			rType = federationTypes.SharedModuleResourceType
		} else {
			return -1
		}

	case rType == composeTypes.ModuleResourceType && strings.HasPrefix(r.Operation, "record.") && r.Operation != "record.create":
		// change resource type from module to record on record read, delete, update operations & remove the prefix
		rType = composeTypes.RecordResourceType
		r.Operation = strings.TrimPrefix(r.Operation, "record.")
	}

	if len(parts) == 3 {
		var ID, _ = strconv.ParseUint(parts[2], 10, 64)

		p1 := uint64(0)
		p2 := uint64(0)

		// exceptions with nested references
		switch rType {
		case composeTypes.ModuleFieldResourceType:
			if ID > 0 {
				if f, ok := rx.fields[ID]; ok {
					p1 = f.NamespaceID
					p2 = f.ModuleID
				}
			}
			r.Resource = composeTypes.ModuleFieldRbacResource(p1, p2, ID)

		case composeTypes.ModuleResourceType:
			if ID > 0 {
				if r, ok := rx.modules[ID]; ok {
					p1 = r.NamespaceID
				}
			}
			r.Resource = composeTypes.ModuleRbacResource(p1, ID)

		// ID belongs to module!
		case composeTypes.RecordResourceType:
			if ID > 0 {
				if r, ok := rx.modules[ID]; ok {
					p1 = r.NamespaceID
				}
			}
			r.Resource = composeTypes.RecordRbacResource(p1, ID, 0)

		case composeTypes.ChartResourceType:
			if ID > 0 {
				if r, ok := rx.charts[ID]; ok {
					p1 = r.NamespaceID
				}
			}
			r.Resource = composeTypes.ChartRbacResource(p1, ID)

		case composeTypes.PageResourceType:
			if ID > 0 {
				if r, ok := rx.pages[ID]; ok {
					p1 = r.NamespaceID
				}
			}
			r.Resource = composeTypes.PageRbacResource(p1, ID)

		case federationTypes.ExposedModuleResourceType:
			if ID > 0 {
				if r, ok := rx.exposedModules[ID]; ok {
					p1 = r.NodeID
				}
			}
			r.Resource = federationTypes.ExposedModuleRbacResource(p1, ID)

		case federationTypes.SharedModuleResourceType:
			if ID > 0 {
				if r, ok := rx.sharedModules[ID]; ok {
					p1 = r.NodeID
				}
			}
			r.Resource = federationTypes.SharedModuleRbacResource(p1, ID)

		default:
			r.Resource = rType + "/" + func() string {
				if ID == 0 {
					return "*"
				}
				return parts[2]
			}()
		}

	} else {
		r.Resource = rType + "/"
	}

	return
}

// helper to preloadresources that may be used when properly constructing rules
func preloadResourceIndex(ctx context.Context, s store.Storer) (*resourceIndex, error) {
	rx := &resourceIndex{}

	rx.modules = make(map[uint64]*composeTypes.Module)
	modules, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{Paging: filter.Paging{Limit: 0}})
	modIDs := make([]uint64, 0, len(modules))
	if err != nil {
		return nil, err
	}
	for _, r := range modules {
		rx.modules[r.ID] = r
		modIDs = append(modIDs, r.ID)
	}

	if len(modIDs) > 0 {
		rx.fields = make(map[uint64]*composeTypes.ModuleField)
		fields, _, err := store.SearchComposeModuleFields(ctx, s, composeTypes.ModuleFieldFilter{ModuleID: modIDs})
		if err != nil {
			return nil, err
		}
		for _, r := range fields {
			rx.fields[r.ID] = r
		}
	}

	rx.charts = make(map[uint64]*composeTypes.Chart)
	chart, _, err := store.SearchComposeCharts(ctx, s, composeTypes.ChartFilter{Paging: filter.Paging{Limit: 0}})
	if err != nil {
		return nil, err
	}
	for _, r := range chart {
		rx.charts[r.ID] = r
	}

	rx.pages = make(map[uint64]*composeTypes.Page)
	page, _, err := store.SearchComposePages(ctx, s, composeTypes.PageFilter{Paging: filter.Paging{Limit: 0}})
	if err != nil {
		return nil, err
	}
	for _, r := range page {
		rx.pages[r.ID] = r
	}

	rx.exposedModules = make(map[uint64]*federationTypes.ExposedModule)
	exposedModule, _, err := store.SearchFederationExposedModules(ctx, s, federationTypes.ExposedModuleFilter{Paging: filter.Paging{Limit: 0}})
	if err != nil {
		return nil, err
	}
	for _, r := range exposedModule {
		rx.exposedModules[r.ID] = r
	}

	rx.sharedModules = make(map[uint64]*federationTypes.SharedModule)
	sharedModule, _, err := store.SearchFederationSharedModules(ctx, s, federationTypes.SharedModuleFilter{Paging: filter.Paging{Limit: 0}})
	if err != nil {
		return nil, err
	}
	for _, r := range sharedModule {
		rx.sharedModules[r.ID] = r
	}

	return rx, nil
}
