package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"

{{- range .imports }}
    "{{ . }}"
{{- end }}
)

type (
	// StoreDecoder is responsible for generating Envoy nodes from already stored
	// resources which can then be managed by Envoy and imported via an encoder.
	StoreDecoder struct{}

	filterWrap struct {
			rt string
			f  envoyx.ResourceFilter
		}
)

const (
	paramsKeyStorer = "storer"
	paramsKeyDAL = "dal"
)

var (
	// @todo temporary fix to make unused pkg/id not throw errors
	_ = id.Next
)

// Decode returns a set of envoy nodes based on the provided params
//
// StoreDecoder expects the DecodeParam of `storer` and `dal` which conform
// to the store.Storer and dal.FullService interfaces.
func (d StoreDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// @todo we can optionally not require them based on what we're doing
	s, err := d.getStorer(p)
	if err != nil {
		return
	}
	dl, err := d.getDal(p)
	if err != nil {
		return
	}

	return d.decode(ctx, s, dl, p)
}

func (d StoreDecoder) decode(ctx context.Context, s store.Storer, dl dal.FullService, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Preprocessing and basic filtering (to omit what this decoder can't handle)
	wrappedFilters := d.prepFilters(p.Filter)
	
	// Get all scoped nodes
	scopedNodes, err := d.getScopeNodes(ctx, s, dl, wrappedFilters)
	if err != nil {
		return
	}

	// Get all reference nodes
	refNodes, refRefs, err := d.getReferenceNodes(ctx, s, dl, wrappedFilters)
	if err != nil {
		return
	}

	// Process filters to get the envoy nodes
	err = func() (err error) {
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
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode filters")
		return
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
		var n *envoyx.Node
		n, err = {{.expIdent}}ToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
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

func {{.expIdent}}ToEnvoyNode(r *types.{{.expIdent}}) (node *envoyx.Node, err error) {
	// Identifiers
	ii := envoyx.MakeIdentifiers(
{{- range .model.attributes -}}
		{{- if not .envoy.identifier -}}
			{{continue}}
		{{ end }}
		r.{{.expIdent}},
{{- end }}
	)

	// Handle references
	// Omit any non-defined values
	refs := map[string]envoyx.Ref{}
{{- range .model.attributes -}}
	{{- if eq .dal.type "Ref" }}
	if r.{{.expIdent}} > 0 {
		refs["{{ .expIdent }}"] = envoyx.Ref{
			ResourceType: "{{ .dal.refModelResType }}",
			Identifiers: envoyx.MakeIdentifiers(r.{{.expIdent}}),
		}
	}
	{{- end }}
{{- end }}

{{ if .envoy.store.extendedRefDecoder }}
	refs = envoyx.MergeRefs(refs, decode{{.expIdent}}Refs(r))
{{ end }}

var scope envoyx.Scope
{{if and .envoy.scoped .parents}}
	scope = envoyx.Scope{
		ResourceType: refs["{{(index .parents 0).refField}}"].ResourceType,
		Identifiers:  refs["{{(index .parents 0).refField}}"].Identifiers,
	}
	for k, ref := range refs {
		// @todo temporary solution to not needlessly scope resources.
		//       Optimally, this would be selectively handled by codegen.
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}

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

	node = &envoyx.Node{
		Resource: r,

		ResourceType: types.{{.expIdent}}ResourceType,
		Identifiers:  ii,
		References: refs,
		Scope: scope,
	}
	return
}

	{{ if not .envoy.store.customFilterBuilder }}
func (d StoreDecoder) make{{.expIdent}}Filter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.{{.expIdent}}Filter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.{{.expIdent}}ID = id.Strings(ids...)

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

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d StoreDecoder) getStorer(p envoyx.DecodeParams) (s store.Storer, err error) {
	aux, ok := p.Params[paramsKeyStorer]
	if ok {
		s, ok = aux.(store.Storer)
		if ok {
			return
		}
	}

	err = errors.Errorf("store decoder expects a storer conforming to store.Storer interface")
	return
}

func (d StoreDecoder) getDal(p envoyx.DecodeParams) (dl dal.FullService, err error) {
	aux, ok := p.Params[paramsKeyDAL]
	if ok {
		dl, ok = aux.(dal.FullService)
		if ok {
			return
		}
	}

	err = errors.Errorf("store decoder expects a DAL conforming to dal.FullService interface")
	return
}

func (d StoreDecoder) prepFilters(ff map[string]envoyx.ResourceFilter) (out []filterWrap) {
	out = make([]filterWrap, 0, len(ff))
	for rt, f := range ff {
		// Handle resources that don't belong to this decoder
		if !strings.HasPrefix(rt, "corteza::{{.componentIdent}}") {
			continue
		}

		out = append(out, filterWrap{rt: rt, f: f})
	}

	return
}

func (d StoreDecoder) getScopeNodes(ctx context.Context, s store.Storer, dl dal.FullService, ff []filterWrap) (scopes envoyx.NodeSet, err error) {
	// Get all requested scopes
	scopes = make(envoyx.NodeSet, len(ff))
	{{ if eq .componentIdent "compose" }}
	err = func() (err error) {
		for i, fw := range ff {
			if fw.f.Scope.ResourceType == "" {
				continue
			}

			// For now the scope can only point to namespace so this will do
			var nn envoyx.NodeSet
			nn, err = d.decodeNamespace(ctx, s, dl, d.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: fw.f.Scope.Identifiers}))
			if err != nil {
				return
			}
			if len(nn) > 1 {
				err = fmt.Errorf("ambiguous scope %v: matches multiple resources", fw.f.Scope)
				return
			}
			if len(nn) == 0 {
				err = fmt.Errorf("invalid scope %v: resource not found", fw.f)
				return
			}

			scopes[i] = nn[0]
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node scopes")
		return
	}
	{{ else }}
	// @note skipping scope logic since it's currently only supported within
	//       Compose resources.
	{{ end }}
	return
}

// getReferenceNodes returns all of the nodes referenced by the nodes defined by the filters
//
// The nodes are provided as a slice (the same order as the filters) and as a map for easier lookups.
func (d StoreDecoder) getReferenceNodes(ctx context.Context, s store.Storer, dl dal.FullService, ff []filterWrap) (nodes []map[string]*envoyx.Node, refs []map[string]envoyx.Ref, err error) {
	nodes = make([]map[string]*envoyx.Node, len(ff))
	refs = make([]map[string]envoyx.Ref, len(ff))
	err = func() (err error) {
		for i, a := range ff {
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
					return err
				}

				// @todo consider changing this.
				//       Currently it's required because the .decode may return some
				//       nested nodes as well.
				//       Consider a flag or a new function.
				aux = envoyx.NodesForResourceType(ref.ResourceType, aux...)
				if len(aux) == 0 {
					return fmt.Errorf("invalid reference %v", ref)
				}
				if len(aux) > 1 {
					return fmt.Errorf("ambiguous reference: too many resources returned %v", a.f)
				}

				auxr[field] = aux[0]
				auxa[field] = aux[0].ToRef()
			}

			nodes[i] = auxr
			refs[i] = auxa
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to decode node references")
		return
	}

	return
}
