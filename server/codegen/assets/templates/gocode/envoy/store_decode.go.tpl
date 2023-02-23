package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"

{{- range .imports }}
    "{{ . }}"
{{- end }}
)

type (
	// StoreDecoder is responsible for fetching already stored Corteza resources
	// which are then managed by envoy and imported via an encoder.
	StoreDecoder struct{}
)

// Decode returns a set of envoy nodes based on the provided params
//
// StoreDecoder expects the DecodeParam of `storer` and `dal` which conform
// to the store.Storer and dal.FullService interfaces.
func (d StoreDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	var (
		s  store.Storer
		dl dal.FullService
	)

	// @todo we can optionally not require them based on what we're doing
	if auxS, ok := p.Params["storer"]; ok {
		s = auxS.(store.Storer)
	}
	if auxDl, ok := p.Params["dal"]; ok {
		dl = auxDl.(dal.FullService)
	}

	return d.decode(ctx, s, dl, p)
}

func (d StoreDecoder) decode(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Transform passed filters into an ordered structure
	type (
		filterWrap struct {
			rt string
			f  envoyx.ResourceFilter
		}
	)
	wrappedFilters := make([]filterWrap, 0, len(p.Filter))
	for rt, f := range p.Filter {
		wrappedFilters = append(wrappedFilters, filterWrap{rt: rt, f: f})
	}

	// Get all requested scopes
	scopedNodes := make(envoyx.NodeSet, len(p.Filter))
	{{ if eq .componentIdent "compose" }}
	for i, a := range wrappedFilters {
		if a.f.Scope.ResourceType == "" {
			continue
		}

		// For now the scope can only point to namespace so this will do
		var nn envoyx.NodeSet
		nn, err = d.decodeNamespace(ctx, s, dl, d.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: a.f.Scope.Identifiers}))
		if err != nil {
			return
		}
		if len(nn) > 1 {
			err = fmt.Errorf("ambiguous scope %v", a.f.Scope)
			return
		}
		if len(nn) == 0 {
			err = fmt.Errorf("invalid scope: resource not found %v", a.f)
			return
		}

		scopedNodes[i] = nn[0]
	}
	{{ else }}
	// @note skipping scope logic since it's currently only supported within
	//       Compose resources.
	{{ end }}

	// Get all requested references
	//
	// Keep an index for the Node and one for the reference to make our
	// lives easier.
	refNodes := make([]map[string]*envoyx.Node, len(p.Filter))
	refRefs := make([]map[string]envoyx.Ref, len(p.Filter))
	for i, a := range wrappedFilters {
		if len(a.f.Refs) == 0 {
			continue
		}

		auxr := make(map[string]*envoyx.Node, len(a.f.Refs))
		auxa := make(map[string]envoyx.Ref)
		for field, ref := range a.f.Refs {
			f := ref.ResourceFilter()
			aux, err := d.decode(ctx, s, dl, envoyx.DecodeParams{
				Type:   envoyx.DecodeTypeStore,
				Filter: f,
			})
			if err != nil {
				return nil, err
			}

			// @todo consider changing this.
			//       Currently it's required because the .decode may return some
			//       nested nodes as well.
			//       Consider a flag or a new function.
			aux = envoyx.NodesForResourceType(ref.ResourceType, aux...)
			if len(aux) == 0 {
				return nil, fmt.Errorf("invalid reference %v", ref)
			}
			if len(aux) > 1 {
				return nil, fmt.Errorf("ambiguous reference: too many resources returned %v", a.f)
			}

			auxr[field] = aux[0]
			auxa[field] = aux[0].ToRef()
		}

		refNodes[i] = auxr
		refRefs[i] = auxa
	}

	var aux envoyx.NodeSet
	for i, wf := range wrappedFilters {
		switch wf.rt {
{{ range .resources -}}
	{{- if .envoy.omit}}{{continue}}{{ end -}}

		case types.{{.expIdent}}ResourceType:
			aux, err = d.decode{{.expIdent}}(ctx, s, dl, d.make{{.expIdent}}Filter(scopedNodes[i], refNodes[i], wf.f))
			if err != nil {
				return
			}
			for _, a := range aux {
				a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
				a.References = envoyx.MergeRefs(a.References, refRefs[i])
			}
			out = append(out, aux...)

{{ end }}
		default:
			aux, err = d.extendDecoder(ctx, s, dl, wf.rt, refNodes[i], wf.f)
			if err!= nil {
				return
			}
			for _, a := range aux {
				a.Identifiers = a.Identifiers.Merge(wf.f.Identifiers)
				a.References = envoyx.MergeRefs(a.References, refRefs[i])
			}
			out = append(out, aux...)
		}
	}

	return
}

{{- range .resources }}
  {{- if .envoy.omit}}
    {{continue}}
  {{ end -}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource {{.ident}}
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) decode{{.expIdent}}(ctx context.Context, s store.Storer, dl dal.FullService, f types.{{.expIdent}}Filter) (out envoyx.NodeSet, err error) {
  // @todo this might need to be improved.
  //       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.Search{{.store.expIdentPlural}}(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		// Identifiers
		ii := envoyx.MakeIdentifiers(
	{{- range .model.attributes -}}
			{{- if not .envoy.identifier -}}
				{{continue}}
			{{ end }}
			r.{{.expIdent}},
	{{- end }}
		)

	refs := map[string]envoyx.Ref{
	{{- range .model.attributes -}}
		{{- if eq .dal.type "Ref" }}
		// Handle references
		"{{ .expIdent }}": envoyx.Ref{
			ResourceType: "{{ .dal.refModelResType }}",
			Identifiers: envoyx.MakeIdentifiers(r.{{.expIdent}}),
		},
		{{- end }}
	{{- end }}
	}

	{{ if .envoy.store.extendedRefDecoder }}
	refs = envoyx.MergeRefs(refs, d.decode{{.expIdent}}Refs(r))
	{{ end }}

	var scope envoyx.Scope
	{{if and .envoy.scoped .parents}}
		scope = envoyx.Scope{
			ResourceType: refs["{{(index .parents 0).refField}}"].ResourceType,
			Identifiers:  refs["{{(index .parents 0).refField}}"].Identifiers,
		}
		for k, ref := range refs {
			ref.Scope = scope
			refs[k] = ref
		}
		{{end}}
		{{if and .envoy.scoped (not .parents)}}
		scope = envoyx.Scope{
			ResourceType: types.{{ .expIdent }}ResourceType,
			Identifiers:  ii,
		}
	{{end}}

		out = append(out, &envoyx.Node{
			Resource: r,

			ResourceType: types.{{.expIdent}}ResourceType,
			Identifiers:  ii,
			References: refs,
			Scope: scope,
		})
	}

	{{ if .envoy.store.extendedDecoder -}}
	aux, err := d.extended{{.expIdent}}Decoder(ctx, s, dl, f, out)
	if err != nil {
		return
	}
	out = append(out, aux...)
	{{- end }}

	return
}

	{{ if not .envoy.store.customFilterBuilder }}
func (d StoreDecoder) make{{.expIdent}}Filter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.{{.expIdent}}Filter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.{{.expIdent}}ID = ids

		{{ if .envoy.store.handleField }}
	if len(hh) > 0 {
		out.{{ .envoy.store.handleField }} = hh[0]
	}
		{{ end }}

	// Refs
	var (
		ar *envoyx.Node
		ok bool
	)
	_ = ar
	_ = ok
		{{ range .model.attributes }}
			{{- if .envoy.store.omitRefFilter }}{{continue}}{{ end }}
			{{ if eq .dal.type "Ref" }}
	ar, ok = refs["{{ .expIdent }}"]
	if ok {
				{{ if .envoy.store.filterRefField -}}
		out.{{ .envoy.store.filterRefField }} = ar.Resource.GetID()
				{{- else -}}
		out.{{ .expIdent }} = ar.Resource.GetID()
				{{- end }}
	}
			{{ end }}
		{{ end }}

		{{- if .envoy.store.extendedFilterBuilder }}
		out = d.extend{{.expIdent}}Filter(scope, refs, auxf, out)
		{{ end -}}
	return
}
	{{ else }}
// Resource should define a custom filter builder
	{{ end }}
{{end}}
