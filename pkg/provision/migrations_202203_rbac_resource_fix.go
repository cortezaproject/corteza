package provision

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

type (
	resIndex struct {
		fields         map[uint64]*composeTypes.ModuleField
		modules        map[uint64]*composeTypes.Module
		records        map[uint64]*composeTypes.Record
		charts         map[uint64]*composeTypes.Chart
		pages          map[uint64]*composeTypes.Page
		exposedModules map[uint64]*federationTypes.ExposedModule
		sharedModules  map[uint64]*federationTypes.SharedModule
	}
)

// MigrateOperations
func migratePost202203RbacRules(ctx context.Context, log *zap.Logger, s store.Storer) error {
	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		if err != nil {
			return err
		}

		log.Info("migrating RBAC resource rules to proper format", zap.Int("rules", len(rr)))

		rx, err := preloadRbacResourceIndex(ctx, s)
		if err != nil {
			return err
		}

		var uniq = make(map[string]bool)
		var uniqID = func(r *rbac.Rule) string {
			return fmt.Sprintf("%s|%s|%d", r.Resource, r.Operation, r.RoleID)
		}

		for _, r := range rr {
			var (
				cr     = *r
				action = migratePost202203RbacRule(r, rx)
			)

			if action != 0 {
				err = store.DeleteRbacRule(ctx, s, &cr)
				if err != nil {
					return fmt.Errorf("could not delete RBAC rule %s: %v", r, err)
				}

				if action == -1 {
					log.Debug("removed obsolete RBAC rule", zap.Stringer("rule", r))
				}
			}

			if action == 1 {
				if uniq[uniqID(r)] {
					log.Warn("skipping duplicate RBAC rule", zap.Stringer("rule", r))
					continue
				}

				err = store.CreateRbacRule(ctx, s, r)
				if err != nil {
					return fmt.Errorf("could not create RBAC rule %s: %v", r, err)
				}

				uniq[uniqID(r)] = true
			}
		}

		return nil
	})
}

//  0 - no action
// -1 - remove
//  1 - update
func migratePost202203RbacRule(r *rbac.Rule, rx *resIndex) (op int) {
	const (
		wildCard = "*"
	)

	if r.Resource == "" {
		return -1
	}

	// split the IDs
	parts := strings.SplitN(r.Resource, "/", 4)
	parts = parts[1:]

	isWildCard := func(s string) bool { return strings.Contains(s, wildCard) }
	parseUint := func(s string) uint64 {
		ID, _ := strconv.ParseUint(s, 10, 64)
		return ID
	}
	validID := func(i uint64) bool { return i > 0 }

	invalid := false
	for i, p := range parts {
		if i == 0 || len(p) == 0 {
			continue
		}

		if !invalid {
			invalid = isWildCard(parts[i-1]) && !isWildCard(p)
		}
	}

	if invalid {
		op = 1

		partsLen := len(parts)
		ID := uint64(0)
		p1 := uint64(0)
		p2 := uint64(0)

		if partsLen == 1 {
			ID = parseUint(parts[0])
		} else if partsLen == 2 {
			p1 = parseUint(parts[0])
			ID = parseUint(parts[1])
		} else if partsLen == 3 {
			p1 = parseUint(parts[0])
			p2 = parseUint(parts[1])
			ID = parseUint(parts[2])
		}

		rType := rbac.ResourceType(r.Resource)
		switch rType {
		case composeTypes.ModuleFieldResourceType:
			if ID > 0 {
				if f, ok := rx.fields[ID]; ok {
					p2 = f.ModuleID
					if m, ok := rx.modules[p2]; ok {
						p1 = m.NamespaceID
						if !validID(p1) {
							return -1
						}
					} else {
						return -1
					}
				} else {
					return -1
				}
			}

			r.Resource = composeTypes.ModuleFieldRbacResource(p1, p2, ID)

		case composeTypes.ModuleResourceType:
			if ID > 0 {
				if r, ok := rx.modules[ID]; ok {
					p1 = r.NamespaceID
					if !validID(p1) {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = composeTypes.ModuleRbacResource(p1, ID)

		case composeTypes.RecordResourceType:
			if ID > 0 {
				if r, ok := rx.records[ID]; ok {
					p2 = r.ModuleID
					if m, ok := rx.modules[p2]; ok {
						p1 = m.NamespaceID
						if !validID(p1) {
							return -1
						}
					} else {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = composeTypes.RecordRbacResource(p1, p2, ID)

		case composeTypes.ChartResourceType:
			if ID > 0 {
				if r, ok := rx.charts[ID]; ok {
					p1 = r.NamespaceID
					if !validID(p1) {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = composeTypes.ChartRbacResource(p1, ID)

		case composeTypes.PageResourceType:
			if ID > 0 {
				if r, ok := rx.pages[ID]; ok {
					p1 = r.NamespaceID
					if !validID(p1) {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = composeTypes.PageRbacResource(p1, ID)

		case federationTypes.ExposedModuleResourceType:
			if ID > 0 {
				if r, ok := rx.exposedModules[ID]; ok {
					p1 = r.NodeID
					if !validID(p1) {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = federationTypes.ExposedModuleRbacResource(p1, ID)

		case federationTypes.SharedModuleResourceType:
			if ID > 0 {
				if r, ok := rx.sharedModules[ID]; ok {
					p1 = r.NodeID
					if !validID(p1) {
						return -1
					}
				} else {
					return -1
				}
			}
			r.Resource = federationTypes.SharedModuleRbacResource(p1, ID)

		default:
			r.Resource = rType + "/" + func() string {
				if ID == 0 {
					return "*"
				}
				return strconv.FormatUint(ID, 10)
			}()
		}
	}
	return
}

// helper to preload resources that may be used when properly constructing rules
func preloadRbacResourceIndex(ctx context.Context, s store.Storer) (*resIndex, error) {
	rx := &resIndex{}

	rx.modules = make(map[uint64]*composeTypes.Module)
	modules, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{
		Paging:  filter.Paging{Limit: 0},
		Deleted: filter.StateInclusive,
	})
	modIDs := make([]uint64, 0, len(modules))
	if err != nil {
		return nil, err
	}
	for _, r := range modules {
		rx.modules[r.ID] = r
		modIDs = append(modIDs, r.ID)

		rx.records = make(map[uint64]*composeTypes.Record)
		records, _, err := store.SearchComposeRecords(ctx, s, r, composeTypes.RecordFilter{
			Deleted: filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		for _, rec := range records {
			rx.records[rec.ID] = rec
		}
	}

	if len(modIDs) > 0 {
		rx.fields = make(map[uint64]*composeTypes.ModuleField)
		fields, _, err := store.SearchComposeModuleFields(ctx, s, composeTypes.ModuleFieldFilter{
			ModuleID: modIDs,
			Deleted:  filter.StateInclusive,
		})
		if err != nil {
			return nil, err
		}
		for _, r := range fields {
			rx.fields[r.ID] = r
		}
	}

	rx.charts = make(map[uint64]*composeTypes.Chart)
	chart, _, err := store.SearchComposeCharts(ctx, s, composeTypes.ChartFilter{
		Paging:  filter.Paging{Limit: 0},
		Deleted: filter.StateInclusive,
	})
	if err != nil {
		return nil, err
	}
	for _, r := range chart {
		rx.charts[r.ID] = r
	}

	rx.pages = make(map[uint64]*composeTypes.Page)
	page, _, err := store.SearchComposePages(ctx, s, composeTypes.PageFilter{
		Paging:  filter.Paging{Limit: 0},
		Deleted: filter.StateInclusive,
	})
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
