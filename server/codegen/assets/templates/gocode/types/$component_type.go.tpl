package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}
{{- /* $structs: names of in-package generated structs (own Clone/Diff exist, so
       Clone deep-copies and Diff dot-path recurses through them). $needsReflect:
       any field compared via reflect.DeepEqual (slice/map/any/ptr-to-non-struct),
       gating the reflect import. */ -}}
{{- $structs := list -}}
{{- range .structDefs }}{{ $structs = append $structs .name }}{{ end -}}
{{- $sliceTypes := .sliceTypes -}}
{{- $reflectTypes := .fields -}}
{{- range .structDefs }}{{ $reflectTypes = concat $reflectTypes .fields }}{{ end -}}
{{- if .needsLabel }}{{ $reflectTypes = append $reflectTypes (dict "goType" "map[string]labelTypes.LabelValue") }}{{ end -}}
{{- if .flags }}{{ $reflectTypes = append $reflectTypes (dict "goType" "[]string") }}{{ end -}}
{{- $needsReflect := false -}}
{{- range $reflectTypes -}}
	{{- $g := .goType -}}
	{{- $base := $g | trimPrefix "*" | trimPrefix "[]" -}}
	{{- if or (hasPrefix "[]" $g) (hasPrefix "map[" $g) (eq $base "any") (and (hasPrefix "*" $g) (not (has $base $structs))) (has $base $sliceTypes) (and (not (has $base $structs)) (or (regexMatch "^[A-Z]" $base) (contains "." $base))) -}}
		{{- $needsReflect = true -}}
	{{- end -}}
{{- end -}}
{{- /* $skip: structDefs with hand-written (custom) Scan/Value. $hasScanValue:
       any structDef gets generated Scan/Value, gating sql/driver/json imports. */ -}}
{{- $skip := .jsonTypesSkip -}}
{{- $hasScanValue := false -}}
{{- range .structDefs }}{{ if not (has .name $skip) }}{{ $hasScanValue = true }}{{ end }}{{ end }}

import (
	{{ if or .jsonTypes $hasScanValue }}"database/sql/driver"
	{{ end }}{{ if or .jsonTypes .idLists $hasScanValue }}"encoding/json"
	{{ end }}{{ if $needsReflect }}"reflect"
	{{ end }}{{ if .idLists }}"strconv"
	{{ end }}{{ if .needsTime }}"time"
	{{ end }}"github.com/crusttech/human/server/pkg/revisions"
	{{ if or .jsonTypes $hasScanValue }}"github.com/crusttech/human/server/pkg/sql"
	{{ end }}{{ if .needsLabel }}labelTypes "github.com/crusttech/human/server/pkg/label/types"
	{{ end }}{{ range .imports }}"{{ . }}"
	{{ end }}
)

type (
	{{ .expIdent }} struct {
	{{ range .fields }}	{{ .expIdent }} {{ .goType }} `{{ .jsonTag }}`
	{{ end }}{{ if .needsLabel }}	Labels map[string]labelTypes.LabelValue `json:"labels,omitempty"`
	{{ end }}{{ if .flags }}	Flags []string `json:"flags,omitempty"`
	{{ end }}}
{{ range .structDefs }}
	{{ if .doc }}// {{ .doc }}
	{{ end }}{{ .name }} struct {
	{{ range .fields }}{{ if .doc }}	// {{ .doc }}
	{{ end }}	{{ .name }} {{ .goType }} `{{ .jsonTag }}`
	{{ end }}}
{{ end }}{{ range .enumDefs }}
	{{ if .doc }}// {{ .doc }}
	{{ end }}{{ .name }} {{ .goKind }}
{{ end }}{{ range .idLists }}
	// {{ . }} is a []uint64 that serializes each element as a JSON string.
	{{ . }} []uint64
{{ end }})

func (r {{ .expIdent }}) Clone() *{{ .expIdent }} {
	dup := r
{{ range .fields }}{{ template "gocode/types/clone-field" (dict "Name" .expIdent "Type" .goType "Structs" $structs) }}{{ end }}{{ if .needsLabel }}	if r.Labels != nil {
		dup.Labels = make(map[string]labelTypes.LabelValue, len(r.Labels))
		for k, v := range r.Labels {
			dup.Labels[k] = v
		}
	}
{{ end }}{{ if .flags }}	if r.Flags != nil {
		dup.Flags = make([]string, len(r.Flags))
		copy(dup.Flags, r.Flags)
	}
{{ end }}	return &dup
}

func (r {{ .expIdent }}) Diff(cmp *{{ .expIdent }}) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &{{ .expIdent }}{}
	}
{{ range .fields }}{{ $key := .jsonTag | trimPrefix `json:"` | splitList "," | first | trimSuffix `"` }}{{ if and (ne $key "-") (ne $key "") }}{{ template "gocode/types/diff-field" (dict "Name" .expIdent "Type" .goType "Key" $key "Structs" $structs "SliceTypes" $sliceTypes) }}{{ end }}{{ end }}{{ if .needsLabel }}	if !reflect.DeepEqual(r.Labels, cmp.Labels) {
		out = append(out, &revisions.Change{Key: "labels", Old: []any{cmp.Labels}, New: []any{r.Labels}})
	}
{{ end }}{{ if .flags }}	if !reflect.DeepEqual(r.Flags, cmp.Flags) {
		out = append(out, &revisions.Change{Key: "flags", Old: []any{cmp.Flags}, New: []any{r.Flags}})
	}
{{ end }}	return out
}
{{ range .structDefs }}{{ $sname := .name }}
func (r {{ $sname }}) Clone() *{{ $sname }} {
	dup := r
{{ range .fields }}{{ template "gocode/types/clone-field" (dict "Name" .name "Type" .goType "Structs" $structs) }}{{ end }}	return &dup
}

func (r {{ $sname }}) Diff(cmp *{{ $sname }}) []*revisions.Change {
	out := make([]*revisions.Change, 0)
	if cmp == nil {
		cmp = &{{ $sname }}{}
	}
{{ range .fields }}{{ $key := .jsonTag | trimPrefix `json:"` | splitList "," | first | trimSuffix `"` }}{{ if and (ne $key "-") (ne $key "") }}{{ template "gocode/types/diff-field" (dict "Name" .name "Type" .goType "Key" $key "Structs" $structs "SliceTypes" $sliceTypes) }}{{ end }}{{ end }}	return out
}
{{ if not (has $sname $skip) }}
func (r *{{ $sname }}) Scan(src any) error          { return sql.ParseJSON(src, r) }
func (r {{ $sname }}) Value() (driver.Value, error) { return json.Marshal(r) }
{{ end }}{{ end }}{{ range .enumDefs }}
const (
{{ range .values }}{{ if .doc }}	// {{ .doc }}
{{ end }}	{{ .constIdent }} {{ .constType }} = {{ .lit }}
{{ end }})
{{ end }}{{ range .idLists }}
func (ll {{ . }}) MarshalJSON() ([]byte, error) {
	ss := make([]string, len(ll))
	for i, id := range ll {
		ss[i] = strconv.FormatUint(id, 10)
	}
	return json.Marshal(ss)
}

func (ll *{{ . }}) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*ll = make({{ . }}, 0, len(raw))
	for _, r := range raw {
		s := string(r)
		if len(s) >= 2 && s[0] == '"' {
			s = s[1 : len(s)-1]
		}
		id, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		*ll = append(*ll, id)
	}
	return nil
}
{{ end }}
{{ range .jsonTypes }}{{ if not (has .name $structs) }}
func (m *{{ .name }}) Scan(src any) error          { return sql.ParseJSON(src, m) }
func (m {{ .name }}) Value() (driver.Value, error) { return json.Marshal(m) }
{{ end }}
{{ if .ptr }}func Parse{{ .name }}(ss []string) (p *{{ .name }}, err error) {
	p = &{{ .name }}{}
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), p)
}{{ else }}func Parse{{ .name }}(ss []string) (p {{ .name }}, err error) {
	if len(ss) == 0 {
		return
	}
	return p, json.Unmarshal([]byte(ss[0]), &p)
}{{ end }}
{{ end }}
{{- define "gocode/types/clone-field" -}}
{{- $g := .Type -}}
{{- $isPtr := hasPrefix "*" $g -}}
{{- $isSlice := hasPrefix "[]" $g -}}
{{- $isMap := hasPrefix "map[" $g -}}
{{- $base := $g | trimPrefix "*" | trimPrefix "[]" -}}
{{- $isGen := has $base .Structs -}}
{{- if and $isGen $isSlice }}	if r.{{ .Name }} != nil {
		dup.{{ .Name }} = make({{ $g }}, len(r.{{ .Name }}))
		for i := range r.{{ .Name }} {
			dup.{{ .Name }}[i] = *r.{{ .Name }}[i].Clone()
		}
	}

{{ else if and $isGen $isPtr }}	if r.{{ .Name }} != nil {
		dup.{{ .Name }} = r.{{ .Name }}.Clone()
	}

{{ else if $isGen }}	dup.{{ .Name }} = *r.{{ .Name }}.Clone()

{{ else if or $isSlice $isMap }}	if r.{{ .Name }} != nil {
		dup.{{ .Name }} = make({{ $g }}, len(r.{{ .Name }}))
{{ if $isMap }}		for k, v := range r.{{ .Name }} {
			dup.{{ .Name }}[k] = v
		}
{{ else }}		copy(dup.{{ .Name }}, r.{{ .Name }})
{{ end }}	}

{{ else if $isPtr }}	if r.{{ .Name }} != nil {
		v := *r.{{ .Name }}
		dup.{{ .Name }} = &v
	}

{{ end -}}
{{- end -}}
{{- define "gocode/types/diff-field" -}}
{{- $g := .Type -}}
{{- $isPtr := hasPrefix "*" $g -}}
{{- $isSlice := hasPrefix "[]" $g -}}
{{- $isMap := hasPrefix "map[" $g -}}
{{- $base := $g | trimPrefix "*" | trimPrefix "[]" -}}
{{- $isGen := has $base .Structs -}}
{{- $isNamedSlice := has $base .SliceTypes -}}
{{- if and $isGen (not $isSlice) (not $isPtr) }}	for _, c := range r.{{ .Name }}.Diff(&cmp.{{ .Name }}) {
		c.Key = "{{ .Key }}." + c.Key
		out = append(out, c)
	}

{{ else if and $isGen $isPtr (not $isSlice) }}	if (r.{{ .Name }} == nil) != (cmp.{{ .Name }} == nil) {
		out = append(out, &revisions.Change{Key: "{{ .Key }}", Old: []any{cmp.{{ .Name }}}, New: []any{r.{{ .Name }}}})
	} else if r.{{ .Name }} != nil {
		for _, c := range r.{{ .Name }}.Diff(cmp.{{ .Name }}) {
			c.Key = "{{ .Key }}." + c.Key
			out = append(out, c)
		}
	}

{{ else if or $isSlice $isMap $isPtr (eq $base "any") $isNamedSlice }}	if !reflect.DeepEqual(r.{{ .Name }}, cmp.{{ .Name }}) {
		out = append(out, &revisions.Change{Key: "{{ .Key }}", Old: []any{cmp.{{ .Name }}}, New: []any{r.{{ .Name }}}})
	}

{{ else if or (regexMatch "^[A-Z]" $base) (contains "." $base) }}	if !reflect.DeepEqual(r.{{ .Name }}, cmp.{{ .Name }}) {
		out = append(out, &revisions.Change{Key: "{{ .Key }}", Old: []any{cmp.{{ .Name }}}, New: []any{r.{{ .Name }}}})
	}

{{ else }}	if r.{{ .Name }} != cmp.{{ .Name }} {
		out = append(out, &revisions.Change{Key: "{{ .Key }}", Old: []any{cmp.{{ .Name }}}, New: []any{r.{{ .Name }}}})
	}

{{ end -}}
{{- end -}}
