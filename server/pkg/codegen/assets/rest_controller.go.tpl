package rest

{{ $resource := $.Endpoint.ControllerName -}}
{{ $svc := $.Endpoint.ServiceField -}}
{{- /* Determine which CRUD apis are present so we can emit conditional imports */ -}}
{{- $hasCreate := false -}}
{{- $hasUpdate := false -}}
{{- $hasDelete := false -}}
{{- range $a := $.Endpoint.Apis -}}
	{{- if $a.IsCRUD -}}
		{{- if eq $a.Name "create" }}{{ $hasCreate = true }}{{ end -}}
		{{- if eq $a.Name "update" }}{{ $hasUpdate = true }}{{ end -}}
		{{- if or (eq $a.Name "delete") (eq $a.Name "undelete") }}{{ $hasDelete = true }}{{ end -}}
	{{- end -}}
{{- end -}}

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// {{ .Source }}

import (
	"context"
{{- if $hasDelete }}
	"github.com/crusttech/human/server/pkg/api"
{{- end }}
	"github.com/crusttech/human/server/{{ .App }}/rest/request"
{{- if or $hasCreate $hasUpdate }}
	"github.com/crusttech/human/server/{{ .App }}/types"
{{- end }}
)

{{ range $a := $.Endpoint.Apis }}
{{- if $a.IsCRUD }}
{{- $req := export $.Endpoint.Entrypoint $a.Name }}
{{- $idField := $a.ResourceIDField }}
{{ if eq $a.Name "list" }}
func (ctrl *{{ $resource }}) List(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	f, err := ctrl.makeFilter(ctx, r)
	if err != nil {
		return nil, err
	}

	set, filter, err := ctrl.{{ $svc }}.Search(ctx, f)
	return ctrl.makeFilterPayload(ctx, set, filter, err)
}
{{ else if eq $a.Name "read" }}
func (ctrl *{{ $resource }}) Read(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	res, err := ctrl.{{ $svc }}.FindByID(ctx, {{ $a.PathIDArgs }})
	return ctrl.makePayload(ctx, res, err)
}
{{ else if eq $a.Name "delete" }}
func (ctrl *{{ $resource }}) Delete(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	return api.OK(), ctrl.{{ $svc }}.DeleteByID(ctx, {{ $a.PathIDArgs }}{{ $a.DeleteExtraArgsCall }})
}
{{ else if eq $a.Name "undelete" }}
func (ctrl *{{ $resource }}) Undelete(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	return api.OK(), ctrl.{{ $svc }}.UndeleteByID(ctx, {{ $a.PathIDArgs }})
}
{{ else if eq $a.Name "create" }}
func (ctrl *{{ $resource }}) Create(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	res := &types.{{ $resource }}{
	{{- range $p := $a.ScopeIDParams }}
		{{ $p.ExportedName }}: r.{{ $p.ExportedName }},
	{{- end }}
	{{- range $p := $a.PlainPostParams }}
		{{ $p.ExportedName }}: r.{{ $p.ExportedName }},
	{{- end }}
	}

	if err := ctrl.beforeCreate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.{{ $svc }}.Create(ctx, res)
{{- if $.Endpoint.GenAfterCreate }}
	if err == nil {
		err = ctrl.afterCreate(ctx, res, r)
	}
{{- end }}
	return ctrl.makePayload(ctx, res, err)
}
{{ else if eq $a.Name "update" }}
func (ctrl *{{ $resource }}) Update(ctx context.Context, r *request.{{ $req }}) (interface{}, error) {
	res := &types.{{ $resource }}{
		ID: {{ $a.ResourceIDArg }},
	{{- range $p := $a.ScopeIDParams }}
		{{ $p.ExportedName }}: r.{{ $p.ExportedName }},
	{{- end }}
	{{- range $p := $a.PlainPostParams }}
		{{ $p.ExportedName }}: r.{{ $p.ExportedName }},
	{{- end }}
	{{- if $a.HasUpdatedAt }}
		UpdatedAt: r.UpdatedAt,
	{{- end }}
	}

	if err := ctrl.beforeUpdate(ctx, res, r); err != nil {
		return nil, err
	}

	res, err := ctrl.{{ $svc }}.Update(ctx, res)
{{- if $.Endpoint.GenAfterUpdate }}
	if err == nil {
		err = ctrl.afterUpdate(ctx, res, r)
	}
{{- end }}
	return ctrl.makePayload(ctx, res, err)
}
{{ end }}
{{- end }}
{{- end }}
