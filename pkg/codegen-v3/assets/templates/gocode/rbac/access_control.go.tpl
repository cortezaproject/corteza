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
	{{- range .Def }}
	{{- $Resource := .Resource }}
	{{- $RbacResource := .RBAC.Resource }}
	{{- range .RBAC.Operations }}
		{
			"type": types.{{ coalesce $Resource }}ResourceType,
			"any": types.{{ coalesce $Resource }}RbacResource({{ range $RbacResource.References }}0,{{ end }}),
			"op": {{ printf "%q" .Operation }},
		},
	{{- end }}
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
	switch rbac.ResourceType(r) {
	{{- range .Def }}
		case types.{{ coalesce .Resource }}ResourceType:
		return rbac{{ .Resource }}ResourceValidator(r, oo...)
	{{- end }}
	}

	return fmt.Errorf("unknown resource type '%q'", r)
}

// rbacResourceOperations returns defined operations for a requested resource
//
// This function is auto-generated
func rbacResourceOperations(r string) map[string]bool {
	switch rbac.ResourceType(r) {
	{{- range .Def }}
	case types.{{ coalesce .Resource }}ResourceType:
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

	if !strings.HasPrefix(r, {{ $GoType }}ResourceType) {
		// expecting resource to always include path
		return fmt.Errorf("invalid resource type")
	}


{{ if .RBAC.Resource.References }}
	const sep = "/"
	var (
		specIdUsed = true

		pp = strings.Split(strings.Trim(r[len({{ $GoType }}ResourceType):], sep), sep)
		prc = []string{
	{{- range .RBAC.Resource.References }}
		{{ printf "%q" .Field }},
	{{- end }}
		}
	)

	if len(pp) != len(prc) {
		return fmt.Errorf("invalid resource path structure")
	}

	for i, p := range pp {
		if p == "*" {
			if !specIdUsed {
				return fmt.Errorf("invalid resource path wildcard level (%d) for {{ .Resource }}", i)
			}

			specIdUsed = false
			continue
		}

		specIdUsed = true
		if _, err := cast.ToUint64E(p); err != nil {
			return fmt.Errorf("invalid reference for %s: '%s'", prc[i], p)
		}
	}
{{- end }}

	return nil
}
{{- end }}

