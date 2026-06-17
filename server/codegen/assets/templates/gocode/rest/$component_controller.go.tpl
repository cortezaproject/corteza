package {{ .package }}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"context"
{{- if .hasDelete }}
	"github.com/crusttech/human/server/pkg/api"
{{- end }}
	"github.com/crusttech/human/server/{{ .app }}/rest/request"
{{- if or .hasCreate .hasUpdate }}
	"github.com/crusttech/human/server/{{ .app }}/types"
{{- end }}
)

{{ range $a := .apis }}
{{- if $a.isCRUD }}
{{ if eq $a.name "list" }}
func (ctrl *{{ $.controllerName }}) List(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.{{ $.serviceField }}.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}
{{ else if eq $a.name "read" }}
func (ctrl *{{ $.controllerName }}) Read(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	res, err := ctrl.{{ $.serviceField }}.FindByID(ctx, {{ $a.pathIDArgs }})
	return ctrl.makePayload(ctx, res, err)
}
{{ else if eq $a.name "delete" }}
func (ctrl *{{ $.controllerName }}) Delete(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	return api.OK(), ctrl.{{ $.serviceField }}.DeleteByID(ctx, {{ $a.pathIDArgs }}{{ $a.deleteExtraArgsCall }})
}
{{ else if eq $a.name "undelete" }}
func (ctrl *{{ $.controllerName }}) Undelete(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	return api.OK(), ctrl.{{ $.serviceField }}.UndeleteByID(ctx, {{ $a.pathIDArgs }})
}
{{ else if eq $a.name "create" }}
func (ctrl *{{ $.controllerName }}) Create(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	res := &types.{{ $.controllerName }}{
	{{- range $p := $a.scopeIDParams }}
		{{ $p.exportedName }}: r.{{ $p.exportedName }},
	{{- end }}
	{{- range $p := $a.plainPostParams }}
		{{ $p.exportedName }}: r.{{ $p.exportedName }},
	{{- end }}
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.{{ $.serviceField }}.Create(ctx, res)
{{- if $.genAfterCreate }}
	if err == nil {
		err = ctrl.afterCreate(ctx, res, r)
	}
{{- end }}
	return ctrl.makePayload(ctx, res, err)
}
{{ else if eq $a.name "update" }}
func (ctrl *{{ $.controllerName }}) Update(ctx context.Context, r *request.{{ $a.reqIdent }}) (interface{}, error) {
	res := &types.{{ $.controllerName }}{
		ID: {{ $a.resourceIDArg }},
	{{- range $p := $a.scopeIDParams }}
		{{ $p.exportedName }}: r.{{ $p.exportedName }},
	{{- end }}
	{{- range $p := $a.plainPostParams }}
		{{ $p.exportedName }}: r.{{ $p.exportedName }},
	{{- end }}
	{{- if $a.hasUpdatedAt }}
		UpdatedAt: r.UpdatedAt,
	{{- end }}
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.{{ $.serviceField }}.Update(ctx, res)
{{- if $.genAfterUpdate }}
	if err == nil {
		err = ctrl.afterUpdate(ctx, res, r)
	}
{{- end }}
	return ctrl.makePayload(ctx, res, err)
}
{{ end }}
{{- end }}
{{- end }}
