package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
	"fmt"
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the {{ .Component }} component
	//
	// This struct is auto-generated
	Component struct {}
)

const (
{{- range .Def }}
	{{ coalesce .Resource "Component" }}ResourceType = "{{ .RBAC.ResourceType }}"
{{- end }}
)


{{- range .Def }}
{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" .Resource }}


// RbacResource returns string representation of RBAC resource for {{ .Resource }} by calling {{ .Resource }}RbacResource fn
//
// RBAC resource is in the {{ .RBAC.ResourceType }}/... format
//
// This function is auto-generated
func (r {{ .Resource }}) RbacResource() string {
	return {{ .Resource }}RbacResource({{ if .RBAC.Resource }}{{ range .RBAC.Resource.References }}r.{{ export .Field }},{{ end }}{{ end }})
}

// {{ .Resource }}RbacResource returns string representation of RBAC resource for {{ .Resource }}
//
// RBAC resource is in the {{ .RBAC.ResourceType }}/{{- if .RBAC.Resource.References }}...{{ end }} format
//
// This function is auto-generated
func {{ .Resource }}RbacResource({{ if .RBAC.Resource }}{{ range .RBAC.Resource.References }}{{ unexport .Field }} uint64,{{ end }}{{ end }}) string {
	{{- if .RBAC.Resource.References }}
	cpts := []interface{{"{}"}}{{"{"}}{{ .Resource }}ResourceType{{"}"}}
	{{- range .RBAC.Resource.References }}
		if {{ unexport .Field }} != 0 {
			cpts = append(cpts, strconv.FormatUint({{ unexport .Field }}, 10))
		} else {
			cpts = append(cpts, "*")
		}

	{{ end }}
	return fmt.Sprintf({{ .Resource }}RbacResourceTpl(), cpts...)
	{{- else }}
	return {{ .Resource }}ResourceType + "/"
	{{- end }}

}

// @todo template
func {{ .Resource }}RbacResourceTpl() string {
	{{- if .RBAC.Resource.References }}
	return "%s
	{{- range .RBAC.Resource.References }}/%s{{- end }}"

	{{- else }}
	return "%s"
	{{- end }}
}

{{ if .RBAC.Resource.Attributes }}
	// RbacAttributes returns resource attributes used for generating list of contextual roles
	//
	// This function is auto-generated
	func (r {{ .Resource }}) RbacAttributes() map[string]interface{} {
		return {{ unexport .Resource }}RbacAttributes(r)
	}

		{{ if .RBAC.Resource.Attributes.Fields }}
		// {{ .Resource }}RbacResource returns string representation of RBAC resource for {{ .Resource }}
		//
		// RBAC resource is in the {{ .RBAC.ResourceType }}/... format
		//
		// This function is auto-generated
		func {{ unexport .Resource }}RbacAttributes(r {{ .Resource }}) map[string]interface{} {
			return map[string]interface{}{
			{{- range .RBAC.Resource.Attributes.Fields }}
				{{ printf "%q" . }}: r.{{ export . }},
			{{- end }}
			}
		}
	{{- end }}
{{- end }}


{{- end }}

