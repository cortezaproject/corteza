package permissions

import (
	"context"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/cortezaproject/corteza-server/pkg/deinterfacer"
)

type (
	Importer struct {
		whitelist whitelistChecker

		// Rules per role
		rules map[string]RuleSet
	}

	whitelistChecker interface {
		Check(*Rule) bool
	}

	ImportKeeper interface {
		Grant(ctx context.Context, rr ...*Rule) error
	}
)

func NewImporter(wl whitelistChecker) *Importer {
	return &Importer{
		whitelist: wl,
	}
}

// CastSet - resolves permission rules for specific resource:
//   <role>: [<operation>, ...]
func (imp *Importer) CastSet(resource, accessStr string, roles interface{}) (err error) {
	if !deinterfacer.IsMap(roles) {
		return errors.New("expecting map of roles")
	}
	return deinterfacer.Each(roles, func(_ int, roleHandle string, oo interface{}) error {
		return imp.appendPermissionRule(roleHandle, accessStr, resource, oo)
	})
}

// CastResourcesSet - resolves permission rules:
//   <role> { <resource>: [<operation>, ...] }
func (imp *Importer) CastResourcesSet(accessStr string, roles interface{}) (err error) {
	// if !IsIterable(roles) {
	// 	return errors.New("expecting map of roles")
	// }

	return deinterfacer.Each(roles, func(_ int, roleHandle string, perResource interface{}) error {
		if !deinterfacer.IsMap(perResource) {
			return errors.New("expecting map of resources")
		}

		// Each resource
		return deinterfacer.Each(perResource, func(_ int, resource string, oo interface{}) error {
			// We want to make life of the person that's preparing the import data easy, so
			// let's do a little guessing instead of him:
			if strings.Contains(resource, ":") {
				// This is not service-level resource, trim * and : from the end
				resource = strings.TrimRight(resource, ":*")
				resource = Resource(resource + string(resourceDelimiter)).AppendWildcard().String()
			}

			return imp.appendPermissionRule(roleHandle, accessStr, resource, oo)
		})
	})
}

func (imp *Importer) appendPermissionRule(roleHandle, accessStr, res string, oo interface{}) (err error) {
	var access Access

	if err = access.UnmarshalJSON([]byte(accessStr)); err != nil {
		return
	}

	if imp.rules == nil {
		imp.rules = map[string]RuleSet{}
	}

	if imp.rules[roleHandle] == nil {
		imp.rules[roleHandle] = RuleSet{}
	}

	operations := deinterfacer.ToStrings(oo)
	if operations == nil {
		return errors.New("could not resolve permission rule operations")
	}

	sort.Strings(operations)

	for _, op := range operations {
		rule := &Rule{
			Access:    access,
			Resource:  Resource(res),
			Operation: Operation(op),
		}

		if imp.whitelist != nil && !imp.whitelist.Check(rule) {
			return errors.Errorf("invalid rule: operation %q on resource %q", op, res)
		}

		imp.rules[roleHandle] = append(imp.rules[roleHandle], rule)
	}

	return nil
}

// UpdateResources iterates over all rules and replaces resource (foo:bar => foo:42)
func (imp *Importer) UpdateResources(base, handle string, ID uint64) {
	var (
		from = Resource(base).append(handle)
		to   = Resource(base).AppendID(ID)
	)
	for _, rules := range imp.rules {
		for _, rule := range rules {
			if rule.Resource == from {
				rule.Resource = to
			}
		}
	}
}

func (imp *Importer) UpdateRoles(handle string, ID uint64) {
	if imp.rules[handle] != nil {
		for _, rule := range imp.rules[handle] {
			rule.RoleID = ID
		}
	}
}

func (imp *Importer) Store(ctx context.Context, k ImportKeeper) (err error) {
	for _, rr := range imp.rules {
		// Make sure all rules have valid role
		rr, _ = rr.Filter(func(rule *Rule) (b bool, e error) {
			return rule.RoleID > 0, nil
		})

		if err = k.Grant(ctx, rr...); err != nil {
			return
		}
	}

	return
}
