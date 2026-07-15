package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

{{- range .types }}

// ProjectRef returns the ID of the project {{ .goType }} belongs to.
//
// It implements actionlog.ProjectResourcer, letting the action log attribute an
// event to the project owning the affected resource.
//
// This function is auto-generated
func (r {{ .goType }}) ProjectRef() uint64 {
	return r.{{ .field }}
}
{{- end }}
