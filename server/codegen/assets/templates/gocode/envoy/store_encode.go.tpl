package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/cortezaproject/corteza/server/pkg/id"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/pkg/errors"
{{- range .imports }}
    "{{ . }}"
{{- end }}
)

type (
  // StoreEncoder is responsible for encoding Corteza resources into the
  // database via the Storer or the DAL interface
	//
	// @todo consider having a different encoder for the DAL resources
	StoreEncoder struct{}
)

{{ $rootRes := .resources }}

// Prepare performs some initial processing on the resource before it can be encoded
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
func (e StoreEncoder) Prepare(ctx context.Context, p envoyx.EncodeParams, rt string, nn envoyx.NodeSet) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
{{- range .resources }}
  {{- if .envoy.omit}}
    {{continue}}
  {{end -}}

	case types.{{.expIdent}}ResourceType:
		return e.prepare{{.expIdent}}(ctx, p, s, nn)
{{ end -}}
	default:
		return e.prepare(ctx, p, s, rt, nn)
	}

	return
}

// Encode encodes the given Corteza resources into the primary store
//
// Encoding should not do any additional processing apart from matching with
// dependencies and runtime validation
//
// The Encode function is called for every resource type where the resource
// appears at the root of the dependency tree.
// All of the root-level resources for that resource type are passed into the function.
// The encoding function must traverse the branches to encode all of the dependencies.
//
// This flow is used to simplify the flow of how resources are encoded into YAML
// (and other documents) as well as to simplify batching.
//
// Encode does not receive any placeholder nodes which are used solely
// for dependency resolution.
func (e StoreEncoder) Encode(ctx context.Context, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	s, err := e.grabStorer(p)
	if err != nil {
		return
	}

	switch rt {
{{- range .resources }}
{{- if .envoy.omit -}}
	{{continue}}
{{end}}
	case types.{{.expIdent}}ResourceType:
		return e.encode{{.expIdent}}s(ctx, p, s, nodes, tree)
{{ end -}}
	default:
		return e.encode(ctx, p, s, rt, nodes, tree)
	}
}

{{- range .resources }}
  {{- if .envoy.omit}}
    {{continue}}
  {{end}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource {{.ident}}
// // // // // // // // // // // // // // // // // // // // // // // // //

// prepare{{.expIdent}} prepares the resources of the given type for encoding
func (e StoreEncoder) prepare{{.expIdent}}(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet) (err error) {
  // Grab an index of already existing resources of this type
  // @note since these resources should be fairly low-volume and existing for
	//       a short time (and because we batch by resource type); fetching them all
	//       into memory shouldn't hurt too much.
	// @todo do some benchmarks and potentially implement some smarter check such as
	//       a bloom filter or something similar.
	
	// Get node scopes
	scopedNodes, err := e.getScopeNodes(ctx, s, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to get scope nodes")
		return
	}

	// Initializing the index here (and using a hashmap) so it's not escaped to the heap
	existing := make(map[int]types.{{.expIdent}}, len(nn))
	err = e.matchup{{.expIdent}}s(ctx, s, existing, scopedNodes, nn)
	if err != nil {
		err = errors.Wrap(err, "failed to matchup existing {{.expIdent}}s")
		return
	}

	for i, n := range nn {
		if n.Resource == nil {
			panic("unexpected state: cannot call prepare{{.expIdent}} with nodes without a defined Resource")
		}

		res, ok := n.Resource.(*types.{{.expIdent}})
		if !ok {
			panic("unexpected resource type: node expecting type of {{.ident}}")
		}

		existing, hasExisting := existing[i]

		// Run expressions on the nodes
		err = e.runEvals(ctx, hasExisting, n)
		if err != nil {
			return
		}

		if hasExisting {
			// On existing, we don't need to re-do identifiers and references; simply
			// changing up the internal resource is enough.
			//
			// In the future, we can pass down the tree and re-do the deps like that
			switch n.Config.MergeAlg {
			case envoyx.OnConflictPanic:
				err = errors.Errorf("resource %v already exists", n.Identifiers.Slice)
				return

			case envoyx.OnConflictReplace:
				// Replace; simple ID change should do the trick
				res.ID = existing.ID

			case envoyx.OnConflictSkip:
				// Replace the node's resource with the fetched one
				res = &existing

				// @todo merging
			}
		} else {
			// @todo actually a bottleneck. As per sonyflake docs, it can at most
			//       generate up to 2**8 (256) IDs per 10ms in a single thread.
			//       How can we improve this?
			res.ID = id.Next()
		}

		// We can skip validation/defaults when the resource is overwritten by
		// the one already stored (the panic one errors out anyway) since it
		// should already be ok.
		if !hasExisting || n.Config.MergeAlg != envoyx.OnConflictSkip {
			err = e.set{{.expIdent}}Defaults(res)
			if err != nil {
				return err
			}

			err = e.validate{{.expIdent}}(res)
			if err != nil {
				return err
			}
		}

		n.Resource = res
	}

	return
}

// encode{{.expIdent}}s encodes a set of resource into the database
func (e StoreEncoder) encode{{.expIdent}}s(ctx context.Context, p envoyx.EncodeParams, s store.Storer, nn envoyx.NodeSet, tree envoyx.Traverser) (err error) {
	for _, n := range nn {
		err = e.encode{{.expIdent}}(ctx, p, s, n, tree)
		if err != nil {
			return
		}
	}

	{{- if .envoy.store.postSetEncoder }}
	err = e.post{{.expIdent}}sEncode(ctx, p, s, tree, nn)
	{{end }}

	return
}

// encode{{.expIdent}} encodes the resource into the database
func (e StoreEncoder) encode{{.expIdent}}(ctx context.Context, p envoyx.EncodeParams, s store.Storer, n *envoyx.Node, tree envoyx.Traverser) (err error) {
	// Grab dependency references
	var auxID uint64
	err = func() (err error) {
		for fieldLabel, ref := range n.References {
			rn := tree.ParentForRef(n, ref)
			if rn == nil {
				err = fmt.Errorf("parent reference %v not found", ref)
				return
			}

			auxID = rn.Resource.GetID()
			if auxID == 0 {
				err = fmt.Errorf("parent reference does not provide an identifier")
				return
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxID)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to set dependency references for %s %v", n.ResourceType, n.Identifiers.Slice))
		return
	}

	{{ if .envoy.store.sanitizeBeforeSave }}
	// Custom resource sanitization before saving.
	// This can be used to cleanup arbitrary config fields.
	e.sanitize{{.expIdent}}BeforeSave(n.Resource.(*types.{{.expIdent}}))
	{{ end }}

  // Flush to the DB
	if !n.Evaluated.Skip {
		err = store.Upsert{{.store.expIdent}}(ctx, s, n.Resource.(*types.{{.expIdent}}))
		if err != nil {
			err = errors.Wrap(err, "failed to upsert {{.expIdent}}")
			return
		}
	}

	{{ $a := . }}

  // Handle resources nested under it
	//
	// @todo how can we remove the OmitPlaceholderNodes call the same way we did for
	//       the root function calls?
	{{/*
		@note this setup will not duplicate encode calls since we only take the
					most specific parent resource.
	*/}}
	{{ $extendedEncoder := .envoy.store.extendedEncoder}}
	{{ if $extendedEncoder }}
	nested := make(envoyx.NodeSet, 0, 10)
	{{ end }}
	err = func() (err error) {
		for rt, nn := range envoyx.NodesByResourceType(tree.Children(n)...) {
			nn = envoyx.OmitPlaceholderNodes(nn...)

			switch rt {
			{{- range $cmp := $rootRes }}
				{{ if $cmp.envoy.omit }}
					{{continue}}
				{{ end }}
				{{ if not $cmp.parents }}
					{{continue}}
				{{ end }}

				{{ $p := index $cmp.parents (sub (len $cmp.parents) 1)}}
					{{ if not (eq $p.handle $a.ident) }}
						{{continue}}
					{{ end }}

				case types.{{$cmp.expIdent}}ResourceType:
					err = e.encode{{$cmp.expIdent}}s(ctx, p, s, nn, tree)
					if err != nil {
						return
					}
					{{ if $extendedEncoder }}
					nested = append(nested, nn...)
					{{ end }}

			{{- end }}
			}
		}

		return
	}()
	if err != nil {
		err = errors.Wrap(err, "failed to encode nested resources")
		return
	}

	{{ if .envoy.store.extendedEncoder }}
	err = e.encode{{.expIdent}}Extend(ctx, p, s, n, nested, tree)
	if err != nil {
		err = errors.Wrap(err, "post encode logic failed with errors")
		return
	}
	{{ end }}

	{{ if .envoy.store.extendedSubResources }}
	err = e.encode{{.expIdent}}ExtendSubResources(ctx, p, s, n, tree)
	if err != nil {
		err = errors.Wrap(err, "failed to encode extended sub resources")
		return
	}
	{{ end }}
	return
}

// matchup{{.expIdent}}s returns an index with indicates what resources already exist
func (e StoreEncoder) matchup{{.expIdent}}s(ctx context.Context, s store.Storer, uu map[int]types.{{.expIdent}}, scopes envoyx.NodeSet, nn envoyx.NodeSet) (err error) {
  // @todo might need to do it smarter then this.
  //       Most resources won't really be that vast so this should be acceptable for now.
	aa, _, err := store.Search{{.store.expIdentPlural}}(ctx, s, types.{{.expIdent}}Filter{})
	if err != nil {
		return
	}

	idMap := make(map[uint64]*types.{{.expIdent}}, len(aa))
	strMap := make(map[string]*types.{{.expIdent}}, len(aa))

	for _, a := range aa {
	{{ range .model.attributes }}
		{{- if not .envoy.identifier -}}
			{{continue}}
		{{- end -}}
		{{- if eq .goType "uint64" -}}
		idMap[a.{{.expIdent}}] = a
		{{- else -}}
		strMap[a.{{.expIdent}}] = a
		{{- end}}
	{{ end }}
	}

	var aux *types.{{.expIdent}}
	var ok bool
	for i, n := range nn {
		{{ if eq .component "compose" }}
		scope := scopes[i]
		if scope == nil {
			continue
		}
		{{ end }}

		for _, idf := range n.Identifiers.Slice {
			if id, err := strconv.ParseUint(idf, 10, 64); err == nil {
				aux, ok = idMap[id]
				if ok {
					uu[i] = *aux
					// When any identifier matches we can end it
					break
				}
			}

			aux, ok = strMap[idf]
			if ok {
				uu[i] = *aux
				// When any identifier matches we can end it
				break
			}
		}
	}

	return
}

{{end}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utility functions
// // // // // // // // // // // // // // // // // // // // // // // // //

func (e *StoreEncoder) grabStorer(p envoyx.EncodeParams) (s store.Storer, err error) {
	auxs, ok := p.Params[paramsKeyStorer]
	if !ok {
		err = errors.Errorf("store encoder expects a store conforming to store.Storer interface")
		return
	}

	s, ok = auxs.(store.Storer)
	if !ok {
		err = errors.Errorf("store encoder expects a store conforming to store.Storer interface")
		return
	}

	return
}

func (e *StoreEncoder) runEvals(ctx context.Context, existing bool, n *envoyx.Node) (err error) {
	// Skip if
	if n.Config.SkipIfEval == nil {
		return
	}

	aux, err := expr.EmptyVars().Cast(map[string]any{
		"missing": !existing,
	})
	if err != nil {
		return
	}

	n.Evaluated.Skip, err = n.Config.SkipIfEval.Test(ctx, aux.(*expr.Vars))
	return
}

func (e StoreEncoder) getScopeNodes(ctx context.Context, s store.Storer, nn envoyx.NodeSet) (scopes envoyx.NodeSet, err error) {
	// Get all requested scopes
	scopes = make(envoyx.NodeSet, len(nn))
	{{ if eq .componentIdent "compose" }}
	err = func() (err error) {
		for i, n := range nn {
			if n.Scope.ResourceType == "" {
				continue
			}

			// For now the scope can only point to namespace so this will do
			var nn envoyx.NodeSet
			nn, err = e.decodeNamespace(ctx, s, e.makeNamespaceFilter(nil, nil, envoyx.ResourceFilter{Identifiers: n.Scope.Identifiers}))
			if err != nil {
				return
			}
			if len(nn) > 1 {
				err = fmt.Errorf("ambiguous scope %v: matches multiple resources", n.Scope)
				return
			}

			// when encoding, it could be missing
			if len(nn) == 0 {
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

{{ if eq .componentIdent "compose" }}
func (e StoreEncoder) decodeNamespace(ctx context.Context, s store.Storer, f types.NamespaceFilter) (out envoyx.NodeSet, err error) {
	// @todo this might need to be improved.
	//       Currently, no resource is vast enough to pose a problem.
	rr, _, err := store.SearchComposeNamespaces(ctx, s, f)
	if err != nil {
		return
	}

	for _, r := range rr {
		var n *envoyx.Node
		n, err = NamespaceToEnvoyNode(r)
		if err != nil {
			return
		}
		out = append(out, n)
	}

	return
}

func (e StoreEncoder) makeNamespaceFilter(scope *envoyx.Node, refs map[string]*envoyx.Node, auxf envoyx.ResourceFilter) (out types.NamespaceFilter) {
	out.Limit = auxf.Limit

	ids, hh := auxf.Identifiers.Idents()
	_ = ids
	_ = hh

	out.NamespaceID = id.Strings(ids...)

	if len(hh) > 0 {
		out.Slug = hh[0]
	}

	if scope == nil {
		return
	}

	if scope.ResourceType == "" {
		return
	}

	// Overwrite it
	out.NamespaceID = id.Strings(scope.Resource.GetID())

	return
}
{{ end }}