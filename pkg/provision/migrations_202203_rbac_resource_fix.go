package provision

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
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
		charts         map[uint64]*composeTypes.Chart
		pages          map[uint64]*composeTypes.Page
		exposedModules map[uint64]*federationTypes.ExposedModule
		sharedModules  map[uint64]*federationTypes.SharedModule
	}

	fetchResIndex struct {
		modules        bool
		fields         bool
		charts         bool
		pages          bool
		exposedModules bool
		sharedModules  bool
	}

	actionType int
)

const (
	actionRemove actionType = -1
	actionNone   actionType = 0
	actionUpdate actionType = 1
)

// migratePost202203RbacRules fixes RBAC rules with wildcard (*) on left-side of an actual ID (<number>)
// It replaces `*` with an actual ID of the resource, so:
// 		corteza::compose:module-field/<ns ID>/<mod ID>/<field ID>
// 		corteza::compose:module/<ns ID>/<module ID>
// 		corteza::compose:record/<ns ID>/<module ID>/<record ID>
// If we have module ID and a wildcard instead of namespace (ns)
// 		then it will find out the ID of the namespace and change the entry in the `rbac_rules` table
func migratePost202203RbacRules(ctx context.Context, log *zap.Logger, s store.Storer) error {
	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		if err != nil {
			return err
		}

		log.Info("migrating RBAC resource rules to proper format", zap.Int("rules", len(rr)))

		var (
			created = 0
			deleted = 0
			skipped = 0

			rx = &resIndex{
				modules:        make(map[uint64]*composeTypes.Module),
				fields:         make(map[uint64]*composeTypes.ModuleField),
				charts:         make(map[uint64]*composeTypes.Chart),
				pages:          make(map[uint64]*composeTypes.Page),
				exposedModules: make(map[uint64]*federationTypes.ExposedModule),
				sharedModules:  make(map[uint64]*federationTypes.SharedModule),
			}
			uniq   = make(map[string]bool)
			uniqID = func(r *rbac.Rule) string {
				return fmt.Sprintf("%s|%s|%d", r.Resource, r.Operation, r.RoleID)
			}

			fetchRes     fetchResIndex
			migrateRules rbac.RuleSet
			deleteRules  rbac.RuleSet
		)

		// check if migration is needed or not
		for _, rule := range rr {
			action := validRbacRule(rule)
			if action == actionRemove {
				deleteRules = append(deleteRules, rule)
				continue
			}

			if action == actionUpdate {
				migrateRules = append(migrateRules, rule)
				rType := rbac.ResourceType(rule.Resource)
				switch rType {
				case composeTypes.ModuleFieldResourceType:
					fetchRes.modules = true
					fetchRes.fields = true
					break
				case composeTypes.ModuleResourceType:
					fetchRes.modules = true
					break
				case composeTypes.RecordResourceType:
					fetchRes.modules = true
					break
				case composeTypes.ChartResourceType:
					fetchRes.charts = true
					break
				case composeTypes.PageResourceType:
					fetchRes.pages = true
					break
				case federationTypes.ExposedModuleResourceType:
					fetchRes.exposedModules = true
					break
				case federationTypes.SharedModuleResourceType:
					fetchRes.sharedModules = true
					break
				}

				continue
			}

			if action == actionNone {
				created++
			}
		}

		// preload required resource to fix Rbac rules
		if len(migrateRules) > 0 {
			rx, err = preloadRbacResourceIndex(ctx, s, fetchRes)
			if err != nil {
				return err
			}
		}

		// create necessary rules
		for _, rule := range migrateRules {
			if uniq[uniqID(rule)] {
				skipped++
				log.Warn("\tskipping duplicate RBAC rule", zap.Stringer("rule", rule))
				continue
			}

			action := migratePost202203RbacRule(ctx, s, rule, rx)
			if action == actionUpdate {
				err = store.CreateRbacRule(ctx, s, rule)
				if err != nil {
					return fmt.Errorf("could not create RBAC rule %s: %v", rule, err)
				}

				created++
				log.Debug("\tcreated RBAC rule", zap.Stringer("rule", rule))
				uniq[uniqID(rule)] = true
				continue
			}

			if action == actionRemove {
				deleteRules = append(deleteRules, rule)
				continue
			}

			skipped++
			log.Warn("\tskipping obsolete RBAC rule", zap.Stringer("rule", rule))
			continue
		}

		// delete necessary rules
		for _, rule := range deleteRules {
			err = store.DeleteRbacRule(ctx, s, rule)
			if err != nil {
				return fmt.Errorf("could not delete RBAC rule %s: %v", rule, err)
			}

			deleted++
			log.Debug("\tremoved obsolete RBAC rule", zap.Stringer("rule", rule))
		}

		log.Info("finished migration of RBAC resource rules",
			zap.Int("total rules", created),
			zap.Int("rules deleted", deleted),
			zap.Int("duplicate rules skipped", skipped),
		)

		return nil
	})
}

// validRbacRule validates the Rbac rule is its in invalid format and
// returns actionRemove / actionNone / actionUpdate
func validRbacRule(r *rbac.Rule) (op actionType) {
	const (
		wildCard = "*"
	)

	if r.Resource == "" {
		return actionRemove
	}

	// split the IDs
	parts := strings.SplitN(r.Resource, "/", 4)
	parts = parts[1:]

	isWildCard := func(s string) bool { return strings.Contains(s, wildCard) }

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
		op = actionUpdate
	} else {
		op = actionNone
	}

	return
}

// migratePost202203RbacRule will replace invalid Rbac rule parts(*) with and actual ID(resource ID),
// if Rbac rule related resource exists then it will return actionUpdate otherwise returns actionRemove
func migratePost202203RbacRule(ctx context.Context, s store.Storer, r *rbac.Rule, rx *resIndex) (op actionType) {
	op = actionUpdate

	// split the IDs
	parts := strings.SplitN(r.Resource, "/", 4)
	parts = parts[1:]

	partsLen := len(parts)
	ID := uint64(0)
	p1 := uint64(0)
	p2 := uint64(0)
	parseUint := func(s string) uint64 {
		return cast.ToUint64(s)
	}
	validID := func(i uint64) bool { return i > 0 }

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
						return actionRemove
					}
				} else {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}

		r.Resource = composeTypes.ModuleFieldRbacResource(p1, p2, ID)

	case composeTypes.ModuleResourceType:
		if ID > 0 {
			if r, ok := rx.modules[ID]; ok {
				p1 = r.NamespaceID
				if !validID(p1) {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}
		r.Resource = composeTypes.ModuleRbacResource(p1, ID)

	case composeTypes.RecordResourceType:
		if ID > 0 {
			// check if moduleID and namespaceID exists
			if p2 > 0 {
				if m, ok := rx.modules[p2]; ok {
					p1 = m.NamespaceID
					if !validID(p1) {
						return actionRemove
					}
				} else {
					return actionRemove
				}
			} else {
				// fetch record from store based on each module,
				// 		to avoid fetching of all record of all module
				for _, mod := range rx.modules {
					if validID(p1) && validID(p2) {
						continue
					}
					rec, err := store.LookupComposeRecordByID(ctx, s, mod, ID)
					if err != nil && validID(rec.ModuleID) && validID(rec.NamespaceID) {
						p1 = rec.NamespaceID
						p2 = rec.ModuleID
					}
				}

				if !validID(p1) || !validID(p2) {
					return actionRemove
				}
			}
		}
		r.Resource = composeTypes.RecordRbacResource(p1, p2, ID)

	case composeTypes.ChartResourceType:
		if ID > 0 {
			if r, ok := rx.charts[ID]; ok {
				p1 = r.NamespaceID
				if !validID(p1) {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}
		r.Resource = composeTypes.ChartRbacResource(p1, ID)

	case composeTypes.PageResourceType:
		if ID > 0 {
			if r, ok := rx.pages[ID]; ok {
				p1 = r.NamespaceID
				if !validID(p1) {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}
		r.Resource = composeTypes.PageRbacResource(p1, ID)

	case federationTypes.ExposedModuleResourceType:
		if ID > 0 {
			if r, ok := rx.exposedModules[ID]; ok {
				p1 = r.NodeID
				if !validID(p1) {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}
		r.Resource = federationTypes.ExposedModuleRbacResource(p1, ID)

	case federationTypes.SharedModuleResourceType:
		if ID > 0 {
			if r, ok := rx.sharedModules[ID]; ok {
				p1 = r.NodeID
				if !validID(p1) {
					return actionRemove
				}
			} else {
				return actionRemove
			}
		}
		r.Resource = federationTypes.SharedModuleRbacResource(p1, ID)

	default:
		r.Resource = rType + "/" + func() string {
			if ID == 0 {
				return "*"
			}
			return cast.ToString(ID)
		}()
	}

	return
}

// helper to preload resources that may be used when properly constructing rules
func preloadRbacResourceIndex(ctx context.Context, s store.Storer, res fetchResIndex) (*resIndex, error) {
	rx := &resIndex{}

	if res.modules {
		rx.modules = make(map[uint64]*composeTypes.Module)
		modules, _, err := store.SearchComposeModules(ctx, s, composeTypes.ModuleFilter{
			Paging:  filter.Paging{Limit: 0},
			Deleted: filter.StateInclusive,
		})
		modIDs := make([]uint64, 0, len(modules))
		if err != nil {
			return nil, err
		}

		if res.fields && len(modIDs) > 0 {
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
	}

	if res.charts {
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
	}

	if res.pages {
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
	}

	if res.exposedModules {
		rx.exposedModules = make(map[uint64]*federationTypes.ExposedModule)
		exposedModule, _, err := store.SearchFederationExposedModules(ctx, s, federationTypes.ExposedModuleFilter{Paging: filter.Paging{Limit: 0}})
		if err != nil {
			return nil, err
		}
		for _, r := range exposedModule {
			rx.exposedModules[r.ID] = r
		}
	}

	if res.sharedModules {
		rx.sharedModules = make(map[uint64]*federationTypes.SharedModule)
		sharedModule, _, err := store.SearchFederationSharedModules(ctx, s, federationTypes.SharedModuleFilter{Paging: filter.Paging{Limit: 0}})
		if err != nil {
			return nil, err
		}
		for _, r := range sharedModule {
			rx.sharedModules[r.ID] = r
		}
	}
	return rx, nil
}
