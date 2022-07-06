package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	systemTypes "github.com/cortezaproject/corteza-server/system/types"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
{{- range .imports }}
    {{ . }}
{{- end }}
)


type (
	roleMemberSearcher interface {
		SearchRoleMembers(context.Context, systemTypes.RoleMemberFilter) (systemTypes.RoleMemberSet, systemTypes.RoleMemberFilter, error)
	}

	rbacService interface {
		Can(rbac.Session, string, rbac.Resource) bool
		Trace(rbac.Session, string, rbac.Resource) *rbac.Trace
		Grant(context.Context, ...*rbac.Rule) error
		FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
	}

	accessControl struct {
		actionlog actionlog.Recorder

		store roleMemberSearcher
		rbac  rbacService
	}
)

func AccessControl(rms roleMemberSearcher) *accessControl {
	return &accessControl{
		store:     rms,
		rbac:      rbac.Global(),
		actionlog: DefaultActionlog,
	}
}

func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	return svc.rbac.Can(rbac.ContextToSession(ctx), op, res)
}

// Effective returns a list of effective permissions for all given resource
//
// This function is auto-generated
func (svc accessControl) Effective(ctx context.Context, rr ... rbac.Resource) (ee rbac.EffectiveSet) {
	for _, res := range rr {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee.Push(r, op, svc.can(ctx, op, res))
		}
	}

	return
}

// Evaluate returns a list of permissions evaluated for the given user/roles combo
//
// This function is auto-generated
func (svc accessControl) Trace(ctx context.Context, userID uint64, roles []uint64, rr ...string) (ee []*rbac.Trace, err error) {
	// Reusing the grant permission since this is who the feature is for
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	var (
		resources []rbac.Resource
		members   systemTypes.RoleMemberSet
	)
	if len(rr) > 0 {
		resources = make([]rbac.Resource, 0, len(rr))
		for _, r := range rr {
			if err = rbacResourceValidator(r); err != nil {
				return nil, fmt.Errorf("can not use resource %q: %w", r, err)
			}

			resources = append(resources, rbac.NewResource(r))
		}
	} else {
		resources = svc.Resources()
	}

	// User ID specified, load its roles
	if userID != 0 {
		if len(roles) > 0 {
			// should be prevented on the client
			return nil, fmt.Errorf("userID and roles are mutually exclusive")
		}

		members, _, err = svc.store.SearchRoleMembers(ctx, systemTypes.RoleMemberFilter{UserID: userID})
		if err != nil {
			return nil, err
		}

		for _, m := range members {
			roles = append(roles, m.RoleID)
		}

		for _, r := range internalAuth.AuthenticatedRoles() {
			roles = append(roles, r.ID)
		}
	}

	if len(roles) == 0 {
		// should be prevented on the client
		return nil, fmt.Errorf("no roles specified")
	}

	session := rbac.ParamsToSession(ctx, userID, roles...)
	for _, res := range resources {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee = append(ee, svc.rbac.Trace(session, op, res))
		}
	}

	return
}

// Resources returns list of resources
//
// This function is auto-generated
func (svc accessControl) Resources() []rbac.Resource {
	return []rbac.Resource{
	{{- range .resources }}
		rbac.NewResource({{ .resFunc }}({{ range .references }}0,{{ end }})),
	{{- end }}
	}
}

// List returns list of operations on all resources
//
// This function is auto-generated
func (svc accessControl) List() (out []map[string]string) {
	def := []map[string]string{
	{{- range .operations }}
		{
			"type": {{ .const }},
			"any": {{ .resFunc }}({{ range .references }}0,{{ end }}),
			"op": {{ printf "%q" .op }},
		},
	{{- end }}
	}

	func(svc interface{}) {
		if svc, is := svc.(interface{}).(interface{ list() []map[string]string }); is {
			def = append(def, svc.list()...)
		}
	}(svc)

	return def
}



// Grant applies one or more RBAC rules
//
// This function is auto-generated
func (svc accessControl) Grant(ctx context.Context, rr ...*rbac.Rule) error {
	if !svc.CanGrant(ctx) {
		// @todo should be altered to check grant permissions PER resource
		return AccessControlErrNotAllowedToSetPermissions()
	}

	for _, r := range rr {
		err := rbacResourceValidator(r.Resource, r.Operation)
		if err != nil {
			return err
		}
	}


	if err := svc.rbac.Grant(ctx, rr...); err != nil {
		return AccessControlErrGeneric().Wrap(err)
	}

	svc.logGrants(ctx, rr)

	return nil
}

// This function is auto-generated
func (svc accessControl) logGrants(ctx context.Context, rr []*rbac.Rule) {
	if svc.actionlog == nil {
		return
	}

	for _, r := range rr {
	    g := AccessControlActionGrant(&accessControlActionProps{r})
	    g.log = r.String()
	    g.resource = r.Resource

	    svc.actionlog.Record(ctx, g.ToAction())
	}
}

// FindRules find all rules based on filters
//
// This function is auto-generated
func (svc accessControl) FindRules(ctx context.Context, roleID uint64, specific filter.State, rr ...string) (out rbac.RuleSet, err error) {
	if !svc.CanGrant(ctx) {
        return nil, AccessControlErrNotAllowedToSetPermissions()
    }

    rules, err := svc.FindRulesByRoleID(ctx, roleID)
    if err != nil {
        return
    }

    var (
        resources  []rbac.Resource
        ruleMap    = make(map[string]bool)
        uniqRuleID = func(r *rbac.Rule) string {
            return fmt.Sprintf("%s|%s|%d", r.Resource, r.Operation, r.RoleID)
        }
    )

    // Filter based on resource
    if len(rr) > 0 {
        resources = make([]rbac.Resource, 0, len(rr))
        for _, r := range rr {
            if err = rbacResourceValidator(r); err != nil {
                return nil, fmt.Errorf("can not use resource %q: %w", r, err)
            }

            resources = append(resources, rbac.NewResource(r))
        }
    } else {
        resources = svc.Resources()
    }

    for _, res := range resources {
        for _, rule := range rules.FilterResource(res.RbacResource()) {
            if _, ok := ruleMap[uniqRuleID(rule)]; !ok {
                out = append(out, rule)
                ruleMap[uniqRuleID(rule)] = true
            }
        }
    }

    // Filter for Excluded, Include, or Exclusive specific rules
    switch specific {
    // Exclude all the specific rules
    case filter.StateExcluded:
        out = out.FilterRules(false)
    // Returns only all the specific rules
    case filter.StateExclusive:
        out = out.FilterRules(true)
    }

    return
}

// FindRulesByRoleID find all rules for a specific role
//
// This function is auto-generated
func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.FindRulesByRoleID(roleID), nil
}

{{- range .operations }}
	// {{ .checkFuncName }} checks if current user can {{ lower .description }}
	//
	// This function is auto-generated
	func (svc accessControl) {{ .checkFuncName }}(ctx context.Context{{ if not .component }}, r *{{ .goType }}{{ end }}) bool {
		{{- if .component }}r := &{{ .goType }}{}{{ end }}
		return svc.can(ctx, {{ printf "%q" .op }}, r)
	}
{{- end }}

// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceType(r) {
	{{- range .validation }}
		case {{ .const }}:
		  return {{ .funcName }}(r, oo...)
	{{- end }}
	}

	return fmt.Errorf("unknown resource type '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceType(r) {
	{{- range .validation }}
	case {{ .const }}:
		return map[string]bool{
		{{- range .operations }}
			{{ printf "%q" . }}: true,
		{{- end }}
		}
	{{- end }}
	}

	return nil
}

{{- range .validation }}

// {{ .funcName }} checks validity of RBAC resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func {{ .funcName }}(r string, oo ...string) error {
	{{- if .references }}
	if !strings.HasPrefix(r, {{ .const }}) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}
	{{ else }}
	if r != {{ .const }} + "/" {
		// expecting resource to always include path
		return fmt.Errorf("invalid component resource, expecting " + {{ .const }} + "/")
	}
	{{ end }}

	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for {{ .label }} resource", o)
		}
	}

	{{ if .references }}
		const sep = "/"
		var (
			pp = strings.Split(strings.Trim(r[len({{ .const }}):], sep), sep)
			prc = []string{
		{{- range .references }}
			{{ printf "%q" . }},
		{{- end }}
			}
		)

		if len(pp) != len(prc) {
			return fmt.Errorf("invalid resource path structure")
		}

		for i := 0; i < len(pp); i++ {
			if pp[i] != "*" {
				if i > 0 && pp[i-1] == "*" {
					return fmt.Errorf("invalid path wildcard level (%d) for {{ .label }} resource", i)
				}

				if _, err := cast.ToUint64E(pp[i]); err != nil {
					return fmt.Errorf("invalid reference for %s: '%s'", prc[i], pp[i])
				}
			}
		}
	{{- end }}
	return nil
}
{{- end }}

