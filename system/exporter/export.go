package exporter

import (
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/permissions"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	ruleFinder interface {
		FindRulesByRoleID(uint64) permissions.RuleSet
	}
)

func ExportableServicePermissions(roles types.RoleSet, rf ruleFinder, access permissions.Access) map[string]map[string][]string {
	var (
		has   bool
		res   string
		rules permissions.RuleSet
		sp    = make(map[string]map[string][]string)
	)

	for _, r := range roles {
		rules = rf.FindRulesByRoleID(r.ID)

		if len(rules) == 0 {
			continue
		}

		for _, rule := range rules {
			if rule.Resource.GetService() != rule.Resource && !rule.Resource.HasWildcard() {
				continue
			}

			if rule.Access != access {
				continue
			}

			res = strings.TrimRight(rule.Resource.String(), ":*")

			if _, has = sp[r.Handle]; !has {
				sp[r.Handle] = map[string][]string{}
			}

			if _, has = sp[r.Handle][res]; !has {
				sp[r.Handle][res] = make([]string, 0)
			}

			sp[r.Handle][res] = append(sp[r.Handle][res], rule.Operation.String())
		}
	}

	return sp
}

func ExportableResourcePermissions(roles types.RoleSet, rf ruleFinder, access permissions.Access, resource permissions.Resource) map[string][]string {
	var (
		has   bool
		rules permissions.RuleSet
		sp    = make(map[string][]string)
	)

	for _, r := range roles {
		rules = rf.FindRulesByRoleID(r.ID)

		if len(rules) == 0 {
			continue
		}

		for _, rule := range rules {
			if rule.Resource != resource {
				continue
			}

			if rule.Access != access {
				continue
			}

			if _, has = sp[r.Handle]; !has {
				sp[r.Handle] = make([]string, 0)
			}

			sp[r.Handle] = append(sp[r.Handle], rule.Operation.String())
		}
	}

	return sp
}
