package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"fmt"
	"strconv"
	"github.com/cortezaproject/corteza/server/pkg/ds"
)

type (
	// Component struct serves as a virtual resource type for the {{ .cmpIdent }} component
	//
	// This struct is auto-generated
	Component struct {}

	indexWrapper struct {
		resource string
		counter  uint
	}
)

var (
	{{/*
		making sure that generated code does not break
		when these packages are not used
  */}}
	_ = fmt.Printf
	_ = strconv.FormatUint
)

var (
	resourceIndex = ds.Trie[uint64, *indexWrapper]()
)

var (
	resourceIndexMaxSize = 1000
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
	cc, ok := ds.TrieSearch[uint64, *indexWrapper](resourceIndex, {{ if not .component }}{{ range .references }}{{ .param }},{{ end }}{{ end }})
	if ok {
		cc.counter++
		return cc.resource
	}

	cpts := []interface{{"{}"}}{{"{"}}{{ .goType }}ResourceType{{"}"}}
	{{- range .references }}
		if {{ .param }} != 0 {
			cpts = append(cpts, strconv.FormatUint({{ .param }}, 10))
		} else {
			cpts = append(cpts, "*")
		}

	{{ end }}
	// Remove the least used ones
	// @todo for now just rebuild the index, later do this properly
	if resourceIndex.Size+1 > resourceIndexMaxSize {
		resourceIndex = ds.Trie[uint64, *indexWrapper]()
	}

	out := fmt.Sprintf({{ .tplFunc }}(), cpts...)
	ds.TrieUpsert[uint64, *indexWrapper](resourceIndex, merge, &indexWrapper{resource: out, counter: 1}, {{ if not .component }}{{ range .references }}{{ .param }},{{ end }}{{ end }})

	return out
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

func merge(a, b *indexWrapper) *indexWrapper {
	a.counter += b.counter
	return a
}
