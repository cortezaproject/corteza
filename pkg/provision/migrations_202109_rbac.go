package provision

import (
	"context"
	"strconv"
	"strings"

	composeTypes "github.com/cortezaproject/corteza-server/compose/types"
	federationTypes "github.com/cortezaproject/corteza-server/federation/types"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
)

// MigrateOperations creates system roles
func migratePre202109RbacRules(ctx context.Context, log *zap.Logger, s store.Storer) error {
	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) error {
		rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
		if err != nil {
			return err
		}

		log.Info("migrating RBAC rules to new format", zap.Int("rules", len(rr)))
		for _, r := range rr {
			var (
				cr     = *r
				action = migratePre202109RbacRule(r)
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
func migratePre202109RbacRule(r *rbac.Rule) (op int) {
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

		// exceptions with nested references
		switch rType {
		case composeTypes.ModuleFieldResourceType:
			r.Resource = composeTypes.ModuleFieldRbacResource(0, 0, ID)
		case composeTypes.ModuleResourceType:
			r.Resource = composeTypes.ModuleRbacResource(0, ID)
		case composeTypes.RecordResourceType:
			// ID belongs to module!
			r.Resource = composeTypes.RecordRbacResource(0, ID, 0)
		case composeTypes.ChartResourceType:
			r.Resource = composeTypes.ChartRbacResource(0, ID)
		case composeTypes.PageResourceType:
			r.Resource = composeTypes.PageRbacResource(0, ID)
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
