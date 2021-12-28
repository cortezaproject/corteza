package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
{{- range .imports }}
    {{ . }}
{{- end }}
)


type (
	accessControl struct {
		actionlog actionlog.Recorder

		rbac interface {
			Can(rbac.Session, string, rbac.Resource) bool
			Grant(context.Context, ...*rbac.Rule) error
			FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
			CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error
		}
	}
)

func AccessControl() *accessControl {
	return &accessControl{
		rbac:      rbac.Global(),
		actionlog: DefaultActionlog,
	}
}


func (svc accessControl) can(ctx context.Context, op string, res rbac.Resource) bool {
	return svc.rbac.Can(rbac.ContextToSession(ctx), op, res)
}

// Effective returns a list of effective permissions for all given resource
func (svc accessControl) Effective(ctx context.Context, rr ... rbac.Resource) (ee rbac.EffectiveSet) {
	for _, res := range rr {
		r := res.RbacResource()
		for op := range rbacResourceOperations(r) {
			ee.Push(r, op, svc.can(ctx, op, res))
		}
	}

	return
}

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

// FindRulesByRoleID find all rules for a specific role
//
// This function is auto-generated
func (svc accessControl) FindRulesByRoleID(ctx context.Context, roleID uint64) (rbac.RuleSet, error) {
	if !svc.CanGrant(ctx) {
		return nil, AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.FindRulesByRoleID(roleID), nil
}

// CloneRulesByRoleID clone all rules of a Role S to a specific Role T
//
// This function is auto-generated
func (svc accessControl) CloneRulesByRoleID(ctx context.Context, fromRoleID uint64, toRoleID ...uint64) error {
	if !svc.CanGrant(ctx) {
		return AccessControlErrNotAllowedToSetPermissions()
	}

	return svc.rbac.CloneRulesByRoleID(ctx, fromRoleID, toRoleID...)
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
	if !strings.HasPrefix(r, {{ .const }}) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}

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

