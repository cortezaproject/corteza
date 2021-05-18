package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	"fmt"
	"github.com/spf13/cast"
	"strings"
	"context"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
{{- range .Imports }}
    {{ normalizeImport . }}
{{- end }}
)


type (
	accessControl struct {
		actionlog actionlog.Recorder

		rbac interface {
			Can(rbac.Session, string, rbac.Resource) bool
			Grant(context.Context, ...*rbac.Rule) error
			FindRulesByRoleID(roleID uint64) (rr rbac.RuleSet)
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
	return []map[string]string{
	{{- range .Def }}
	{{- $Schema := .RBAC.Schema }}
	{{- range .RBAC.Operations }}
		{ "resource": {{ printf "%q" $Schema }}, "operation": {{ printf "%q" .Operation }} },
	{{- end }}
	{{- end }}
	}
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


{{- range .Def }}
	{{ $GoType   := printf "types.%s" (.Resource) }}

	{{ if .IsComponentResource }}

		{{- range .RBAC.Operations }}
		// {{ export .CanFnName }} checks if current user can {{ lower .Description }}
		//
		// This function is auto-generated
		func (svc accessControl) {{ export .CanFnName }}(ctx context.Context) bool {
			return svc.can(ctx, {{ printf "%q" .Operation }}, &types.Component{})
		}
		{{- end }}

	{{ else }}

		{{ $ResStruct := .RBAC.Resource.Elements }}

		{{- range .RBAC.Operations }}
		// {{ export .CanFnName }} checks if current user can {{ lower .Description }}
		//
		// This function is auto-generated
		func (svc accessControl) {{ export .CanFnName }}(ctx context.Context, r * {{ $GoType }}) bool {
			return svc.can(ctx, {{ printf "%q" .Operation }}, r)
		}
		{{- end }}

{{ end }}
{{- end }}


// rbacResourceValidator validates known component's resource by routing it to the appropriate validator
//
// This function is auto-generated
func rbacResourceValidator(r string, oo ...string) error {
	switch rbac.ResourceSchema(r) {
	{{- range .Def }}
	case {{ printf "%q" .RBAC.Schema }}:
		return rbac{{ .Resource }}ResourceValidator(r, oo...)
	{{- end }}
	}

	return fmt.Errorf("unknown resource schema '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceSchema(r) {
	{{- range .Def }}
	case {{ printf "%q" .RBAC.Schema }}:
		return map[string]bool{
		{{- range .RBAC.Operations }}
			{{ printf "%q" .Operation }}:  true,
		{{- end }}
		}
	{{- end }}
	}

	return nil
}

{{- range .Def }}

{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" (.Resource) }}

// rbac{{ .Resource }}ResourceValidator checks validity of rbac resource and operations
//
// Can be called without operations to check for validity of resource string only
//
// This function is auto-generated
func rbac{{ .Resource }}ResourceValidator(r string, oo ...string) error {
	defOps := rbacResourceOperations(r)
	for _, o := range oo {
		if !defOps[o] {
			return fmt.Errorf("invalid operation '%s' for {{ .Component }}{{ if not .IsComponentResource }} {{ .Resource }}{{end }} resource", o)
		}
	}

	if !strings.HasPrefix(r, {{ $GoType }}RbacResourceSchema + ":/") {
		return fmt.Errorf("invalid schema")
	}

{{ if not .IsComponentResource }}
	pp := strings.Split(r[len({{ $GoType }}RbacResourceSchema)+2:], "/")
	if len(pp) != {{ len .RBAC.Resource.Elements }} {
		return fmt.Errorf("invalid resource path")
	}
{{- end }}

{{ if .RBAC.Resource.Elements }}
	var (
		ppWildcard   bool
		pathElements = []string{
	{{- range .RBAC.Resource.Elements }}
		{{ printf "%q" . }},
	{{- end }}
		}
	)

	for i, p := range pp {
		if p == "*" {
			ppWildcard = true
			continue
		}

		if !ppWildcard {
			return fmt.Errorf("invalid resource path wildcard level")
		}

		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid ID for %s: '%s'", pathElements[i], p)
		}
	}
{{- end }}

	return nil
}
{{- end }}

