package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"strconv"
)

type (
	// Component struct serves as a virtual resource type for the {{ .cmpIdent }} component
	//
	// This struct is auto-generated
	Component struct {}
)

var (
	{{/*
		making sure that generated code does not break
		when these packages are not used
  */}}
	_ = fmt.Printf
	_ = strconv.FormatUint
)

{{- range .types }}

// RbacResource returns string representation of RBAC resource for {{ .goType }} by calling {{ .resFunc }} fn
//
// RBAC resource is in the {{ .type }}/... format
//
// This function is auto-generated
func (r {{ .goType }}) RbacResource() string {
	return {{ .resFunc }}({{ if not .component }}{{ range .references }}r.{{ .refField }},{{ end }}{{ end }})
}

// {{ .resFunc }} returns string representation of RBAC resource for {{ .goType }}
//
// RBAC resource is in the {{ .type }}/{{- if .references }}...{{ end }} format
//
// This function is auto-generated
func {{ .resFunc }}({{ if not .component }}{{ range .references }}{{ .param }} uint64,{{ end }}{{ end }}) string {
	{{- if .references }}
	cpts := []interface{{"{}"}}{{"{"}}{{ .goType }}ResourceType{{"}"}}
	{{- range .references }}
		if {{ .param }} != 0 {
			cpts = append(cpts, strconv.FormatUint({{ .param }}, 10))
		} else {
			cpts = append(cpts, "*")
		}

	{{ end }}
	return fmt.Sprintf({{ .tplFunc }}(), cpts...)
	{{- else }}
	return {{ .goType }}ResourceType + "/"
	{{- end }}

}

func {{ .tplFunc }}() string {
	{{- if .references }}
	return "%s
	{{- range .references }}/%s{{- end }}"

	{{- else }}
	return "%s"
	{{- end }}
}

{{- end }}

