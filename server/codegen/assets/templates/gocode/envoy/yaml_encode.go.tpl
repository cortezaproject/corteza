package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"io"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
	"github.com/pkg/errors"
{{- range .imports }}
    "{{ . }}"
{{- end }}
)

{{ $rootRes := .resources }}

type (
  // YamlEncoder is responsible for encoding Corteza resources into
  // a YAML supported format
	YamlEncoder struct{}
)

const (
	paramsKeyWriter = "writer"
)

// Encode encodes the given Corteza resources into some YAML supported format
//
// Encoding should not do any additional processing apart from matching with
// dependencies and runtime validation
//
// Preparation runs validation, default value initialization, matching with
// already existing instances, ...
//
// The prepare function receives a set of nodes grouped by the resource type.
// This enables some batching optimization and simplifications when it comes to
// matching with existing resources.
//
// Prepare does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e YamlEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (err error) {
	var (
		out *yaml.Node
		aux *yaml.Node
	)
	_ = aux

	w, err := e.getWriter(p)
	if err != nil {
		return
	}

	switch rt {
{{- range .resources }}
	{{- if or .envoy.omit .envoy.yaml.omitEncoder}}
		{{continue}}
	{{ end -}}

	case types.{{.expIdent}}ResourceType:
		aux, err = e.encode{{.expIdent}}s(ctx, p, nodes, tt)
		if err != nil {
			return
		}
    // Root level resources are always encoded as a map
		out, err = y7s.AddMap(out, "{{.ident}}", aux)
		if err != nil {
			return
		}
{{ end -}}
	default:
		out, err = e.encode(ctx, out, p, rt, nodes, tt)
		if err != nil {
			return
		}
	}

	// Don't output nil values since that will produce broken yaml docs
	if out == nil {
		return
	}

	return yaml.NewEncoder(w).Encode(out)
}

{{- range .resources }}
	{{- if .envoy.omit}}
		{{continue}}
	{{ end -}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource {{.ident}}
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encode{{.expIdent}}s(ctx context.Context, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	var aux *yaml.Node
	for _, n := range nodes {
		aux, err = e.encode{{.expIdent}}(ctx, p, n, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddSeq(out, aux)
		if err != nil {
			return
		}
	}

	return
}

// encode{{.expIdent}} focuses on the specific resource invoked by the Encode method
func (e YamlEncoder) encode{{.expIdent}}(ctx context.Context, p envoyx.EncodeParams, node *envoyx.Node, tt envoyx.Traverser) (out *yaml.Node, err error) {
	res := node.Resource.(*types.{{.expIdent}})
	{{ $res := .expIdent }}

	// Pre-compute some map values so we can omit error checking when encoding yaml nodes
	{{ range .model.attributes -}}
		{{- if eq .dal.type "Timestamp" -}}
			{{- if .dal.nullable -}}
		aux{{.expIdent}}, err := e.encodeTimestampNil(p, res.{{.expIdent}})
		if err != nil {
			return
		}
			{{- else -}}
		aux{{.expIdent}}, err := e.encodeTimestamp(p, res.{{.expIdent}})
		if err != nil {
			return
		}
			{{- end -}}
		{{- else if eq .dal.type "Ref" -}}
		aux{{.expIdent}}, err := e.encodeRef(p, res.{{.expIdent}}, "{{.expIdent}}", node, tt)
		if err != nil {
			return
		}
		{{- else if .envoy.yaml.customEncoder -}}
		aux{{.expIdent}}, err := e.encode{{$res}}{{.expIdent}}C(ctx, p, tt, node, res, res.{{.expIdent}})
		if err != nil {
			return
		}
		{{- end }}
	{{end}}

	out, err = y7s.AddMap(out,
	{{ range .model.attributes -}}
		{{- if .envoy.yaml.omitEncoder }}{{continue}}{{ end -}}
		{{- if .envoy.yaml.customEncoder -}}
		"{{.envoy.yaml.identKeyEncode}}", aux{{.expIdent}},
		{{- else if eq .dal.type "Timestamp" -}}
			{{- if .dal.nullable -}}
		"{{.envoy.yaml.identKeyEncode}}", aux{{.expIdent}},
			{{- else -}}
		"{{.envoy.yaml.identKeyEncode}}", aux{{.expIdent}},
			{{- end -}}
		{{- else if eq .dal.type "Ref" -}}
		"{{.envoy.yaml.identKeyEncode}}", aux{{.expIdent}},
		{{- else -}}
		"{{.envoy.yaml.identKeyEncode}}", res.{{.expIdent}},
		{{- end }}
	{{end}}
	)
	if err != nil {
		return
	}

	// Handle nested resources
	var aux *yaml.Node
	_ = aux
	{{ $a := . }}
	{{- range $cmp := $rootRes }}
		{{ if $cmp.envoy.omit }}{{continue}}{{ end }}
		{{if not $cmp.parents}}{{continue}}{{end}}

		{{/*
			Only handle resources where the current resource would appear last
			in the list of parents.
			Since parents are ordered by _importance_ this removes the danger of
			multiple parents decoding the same resource.
		*/}}
		{{ $p := (index $cmp.parents (sub (len $cmp.parents) 1)) }}
		{{ if (eq $p.handle $a.ident) }}
			{{ if eq $cmp.ident "moduleField"}}
	// When processing module fields, we need to filter out the ones that
	// don't belong to this module
	//
	// @todo offload this to dependency resolution; this is a hack
	children := tt.ChildrenForResourceType(node, types.ModuleFieldResourceType)

	var proc envoyx.NodeSet
	selfRef := node.ToRef()
	for _, c := range children {
		if c.References["ModuleID"].Equals(selfRef) {
			proc = append(proc, c)
		}
	}

	aux, err = e.encodeModuleFields(ctx, p, proc, tt)
			{{- else }}
	aux, err = e.encode{{$cmp.expIdent}}s(ctx, p, tt.ChildrenForResourceType(node, types.{{$cmp.expIdent}}ResourceType), tt)
			{{- end }}
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"{{$cmp.ident}}", aux,
	)
	if err != nil {
		return
	}
		{{ end }}
	{{- end }}

	{{- range .envoy.yaml.extendedResourceEncoders }}
	aux, err = e.encode{{.expIdent}}s(ctx, p, tt.ChildrenForResourceType(node, types.{{.expIdent}}ResourceType), tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out,
		"{{.identKey}}", aux,
	)
	if err != nil {
		return
	}
	{{- end }}


	return
}

{{ end -}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Encoding utils
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) encodeTimestamp(p envoyx.EncodeParams, t time.Time) (any, error) {
	if t.IsZero() {
		return nil, nil
	}

	tz := p.Encoder.PreferredTimezone
	if tz != "" {
		tzL, err := time.LoadLocation(tz)
		if err != nil {
			return nil, err
		}
		t = t.In(tzL)
	}

	ly := p.Encoder.PreferredTimeLayout
	if ly == "" {
		ly = time.RFC3339
	}

	return t.Format(ly), nil
}

func (e YamlEncoder) encodeTimestampNil(p envoyx.EncodeParams, t *time.Time) (any, error) {
	if t == nil { return nil, nil }

	return e.encodeTimestamp(p, *t)
}

func (e YamlEncoder) encodeRef(p envoyx.EncodeParams, id uint64, field string, node *envoyx.Node, tt envoyx.Traverser) (any, error) {
	parent := tt.ParentForRef(node, node.References[field])

	// @todo should we panic instead?
	//       for now gracefully fallback to the ID
	if parent == nil {
		return id, nil
	}

	return parent.Identifiers.FriendlyIdentifier(), nil
}


// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility functions
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e YamlEncoder) getWriter(p envoyx.EncodeParams) (out io.Writer, err error) {
	aux, ok := p.Params[paramsKeyWriter]
	if ok {
		out, ok = aux.(io.Writer)
		if ok {
			return
		}
	}

	// @todo consider adding support for managing files from a location
	err = errors.Errorf("YAML encoder expects a writer conforming to io.Writer interface")
	return
}

func safeParentIdentifier(tt envoyx.Traverser, n *envoyx.Node, ref envoyx.Ref) (out string) {
	aux := tt.ParentForRef(n, ref)
	if aux == nil {
		return ref.Identifiers.FriendlyIdentifier()
	}

	return aux.Identifiers.FriendlyIdentifier()
}
