package {{ .Package }}

{{ template "header-gentext.tpl" }}
{{ template "header-definitions.tpl" . }}

import (
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
	{{ coalesce .Resource "Component" }}RbacResourceSchema = "{{ .RBAC.Schema }}"
{{- end }}
)


{{- range .Def }}
{{ $Resource := .Resource }}
{{ $GoType   := printf "types.%s" .Resource }}


// RbacResource returns string representation of RBAC resource for {{ .Resource }} by calling {{ .Resource }}RbacResource fn
//
// RBAC resource is in the {{ .RBAC.Schema }}:/... format
//
// This function is auto-generated
func (r {{ .Resource }}) RbacResource() string {
	return {{ .Resource }}RbacResource({{ if .RBAC.Resource }}{{ range .RBAC.Resource.Elements }}r.{{ export . }},{{ end }}{{ end }})
}

// {{ .Resource }}RbacResource returns string representation of RBAC resource for {{ .Resource }}
//
// RBAC resource is in the {{ .RBAC.Schema }}:/... format
//
// This function is auto-generated
func {{ .Resource }}RbacResource({{ if .RBAC.Resource }}{{ range .RBAC.Resource.Elements }}{{ unexport . }} uint64,{{ end }}{{ end }}) string {
	out := {{ .Resource }}RbacResourceSchema + ":"
	{{- range .RBAC.Resource.Elements }}
		out += "/"

		if {{ unexport . }} != 0 {
			out += strconv.FormatUint({{ unexport . }}, 10)
		} else {
			out += "*"
		}
	{{- end }}
	return out
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
		// RBAC resource is in the {{ .RBAC.Schema }}:/... format
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

