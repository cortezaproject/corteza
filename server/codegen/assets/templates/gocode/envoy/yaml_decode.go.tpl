package {{ .package }}

{{ template "gocode/header-gentext.tpl" }}

import (
	"strings"
	"context"
	"io"
	"os"

	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"golang.org/x/text/language"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v3"
	"github.com/pkg/errors"

{{- range .imports }}
    "{{ . }}"
{{- end }}
)

{{$cmpIdent := .componentIdent}}
{{$rootRes := .resources}}

type (
	// YamlDecoder is responsible for decoding YAML documents into Corteza resources
	// which are then managed by envoy and imported via an encoder.
	YamlDecoder struct{}

	// documentContext provides a bit of metadata to the decoder such as
	// root-level reference definitions (such as the namespace)
	documentContext struct {
		// references holds any references defined on the root of the document
		//
		// Example of defining a namespace reference:
		// namespace: test_import_namespace
		// modules:
		//   - name: Test Import Module
		//		 handle: test_import_module
		references map[string]string

		// parentIdent holds the identifier of the parent resource (if nested).
		// The parent ident can be handy if the resource requires the info before the
		// original handling is finished (such as record datasource nodes).
		parentIdent envoyx.Identifiers
	}

	// auxYamlDoc is a helper struct that registers custom YAML decoders
	// to assist the package
	auxYamlDoc  struct {
		nodes envoyx.NodeSet
	}
)

const (
	paramsKeyStream = "stream"
)

// CanFile returns true if the provided file can be decoded with this decoder
func (d YamlDecoder) CanFile(f *os.File) (ok bool) {
	// @todo improve this check; for now a simple extension check should do the trick
	return d.CanExt(f.Name())
}

// CanExt returns true if the provided file extension can be decoded with this decoder
func (d YamlDecoder) CanExt(name string) (ok bool) {
	var (
			pt  = strings.Split(name, ".")
			ext = strings.TrimSpace(pt[len(pt)-1])
		)
		return ext == "yaml" || ext == "yml"
}

// Decode returns a set of envoy nodes based on the provided params
//
// YamlDecoder expects the DecodeParam of `stream` which conforms
// to the io.Reader interface.
func (d YamlDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	r, err := d.getReader(ctx, p)
	if err != nil {
		return
	}

	// Offload decoding to the aux document
	doc := &auxYamlDoc{}
	err = yaml.NewDecoder(r).Decode(doc)
	if err != nil {
		err = errors.Wrap(err, "{{$cmpIdent}} yaml decoder: failed to decode document")
		return
	}

	return doc.nodes, nil
}

func (d *auxYamlDoc) UnmarshalYAML(n *yaml.Node) (err error) {
	// Get the document context from the root level
	dctx, err := d.getDocumentContext(n)
	if err != nil {
		return
	}

	var aux envoyx.NodeSet
	return y7s.EachMap(n, func(k, v *yaml.Node) error {
		kv := strings.ToLower(k.Value)

		switch kv {
		// Decode all resources under the {{$cmpIdent}} component
{{- range .resources -}}
	{{- if .envoy.omit -}}
		{{continue}}
	{{- end -}}

		{{ $identKeys := .envoy.yaml.identKeys }}
		case {{ range $i, $l := $identKeys -}}
			"{{ $l }}"{{if not (eq $i (sub (len $identKeys) 1))}},{{end}}
		{{- end}}:
	{{- if .envoy.yaml.supportMappedInput }}
			if y7s.IsMapping(v) {
				aux, err = d.unmarshal{{.expIdent}}Map(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
	{{ else }}
			// @note {{ .ident }} doesn't support mapped inputs. This can be
			//       changed in the .cue definition under the
			//       .envoy.yaml.supportMappedInput field.
	{{- end }}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshal{{.expIdent}}Seq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: {{.ident}}")
			}
			return err
{{ end }}

	{{/* @todo this should be improved and offloaded to the cue instead of template sys */}}
	{{- if eq .componentIdent "system" }}
		// Decode access control nodes
		case "allow":
			aux, err = unmarshalAllowNode(v)
			d.nodes = append(d.nodes, aux...)
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: RBAC allow")
			}

		case "deny":
			aux, err = unmarshalDenyNode(v)
			d.nodes = append(d.nodes, aux...)
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: RBAC deny")
			}

		// Resource translation nodes
		case "locale", "translation", "translations", "i18n":
			aux, err = unmarshalLocaleNode(v)
			d.nodes = append(d.nodes, aux...)
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: locale")
			}
	{{ end }}

		// Offload to custom handlers
		default:
			aux, err = d.unmarshalYAML(kv, v)
			d.nodes = append(d.nodes, aux...)
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node")
			}
		}
		return nil
	})
}

{{- range .resources }}
	{{- if .envoy.omit}}
		{{continue}}
	{{ end -}}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource {{.ident}}
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshal{{.expIdent}}Seq unmarshals {{.expIdent}} when provided as a sequence node
func (d *auxYamlDoc) unmarshal{{.expIdent}}Seq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshal{{ .expIdent }}Node(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshal{{.expIdent}}Map unmarshals {{.expIdent}} when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
{{- if not .envoy.yaml.supportMappedInput }}
// @note this resource does not support map encoding.
//       Refer to the corresponding definition files to adjust if needed.
{{ else }}
func (d *auxYamlDoc) unmarshal{{ .expIdent }}Map(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshal{{ .expIdent }}Node(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}
{{ end }}

{{ range .envoy.yaml.extendedResourceDecoders -}}
// unmarshal{{.expIdent}}ExtendedSeq unmarshals {{.expIdent}} when provided as a sequence node
func (d *auxYamlDoc) unmarshalExtended{{.expIdent}}Seq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshal{{ .expIdent }}ExtendedNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}
{{if not .supportMappedInput}}{{continue}}{{end}}
// unmarshal{{.expIdent}}ExtendedMap unmarshals {{.expIdent}} when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalExtended{{ .expIdent }}Map(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshal{{ .expIdent }}ExtendedNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}
{{ end }}

// unmarshal{{ .expIdent }}Node is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshal{{ .expIdent }}Node(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.{{ .expIdent }}

	// @todo we're omitting errors because there will be a bunch due to invalid
	//       resource field types. This might be a bit unstable as other errors may
	//       also get ignored.
	//
	//       A potential fix would be to firstly unmarshal into an any, check errors
	//       and then unmarshal into the resource while omitting errors.
	n.Decode(&r)

	// Identifiers are determined manually when iterating the yaml node.
	// This is to help assure there are no duplicates and everything
	// was accounted for especially when working with aliases such as
	// user_name instead of userName.
	ii := envoyx.Identifiers{}

	{{ if .envoy.yaml.supportMappedInput}}
	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "{{ .envoy.yaml.mappedField }}", &r.{{ .envoy.yaml.mappedField }})
		ii = ii.Add(r.{{ .envoy.yaml.mappedField }})
	}
	{{ end }}

	var (
		refs = make(map[string]envoyx.Ref)
		auxOut envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		{{- if .rbac }}
		rbacNodes envoyx.NodeSet
		{{- end }}
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {
		{{ $resource := . }}
		{{/*
			Iterate over all model attributes and handle
				- attribute aliases
				- identifiers
				- reference
				- custom decoding logic

			Generic field decoding is already handled at the top with the generic
			node.Decode(...) function call.
		*/}}
		{{- range $attr := .model.attributes }}

		{{/*
		  We can skip the attributes when...
			- there aren't any ident key aliases
			- not an identifier
			- not a reference
			- no custom logic
		*/}}
		{{
			if and
				(or (not $attr.envoy.yaml.identKeyAlias) (eq (len $attr.envoy.yaml.identKeyAlias) 0))
				(not $attr.envoy.identifier)
				(not (eq $attr.dal.type "Ref"))
				(not $attr.envoy.yaml.customDecoder)
		}}
			{{continue}}
		{{ end }}

		{{ $identKeys := .envoy.yaml.identKeys }}
		case {{ range $i, $l := $identKeys -}}
		"{{ $l }}"{{if not (eq $i (sub (len $identKeys) 1))}},{{end}}
		{{- end}}:

			{{- if and $attr.envoy.yaml.identKeyAlias (gt (len $attr.envoy.yaml.identKeyAlias) 0) (not (eq $attr.dal.type "Ref")) }}
				// Handle field alias
				//
				// @todo consider adding an is empty check before overwriting
				err = y7s.DecodeScalar(n, "{{ $attr.ident }}", &r.{{ $attr.expIdent }})
				if err != nil {
					return err
				}
			{{- end }}

			{{- if $attr.envoy.identifier }}
				// Handle identifiers
				err = y7s.DecodeScalar(n, "{{ $attr.ident }}", &auxNodeValue)
				if err != nil {
					return err
				}
				ii = ii.Add(auxNodeValue)
			{{- end }}

			{{- if eq $attr.dal.type "Ref" }}
				// Handle references
				err = y7s.DecodeScalar(n, "{{ $attr.ident }}", &auxNodeValue)
				if err != nil {
					return err
				}

				// Omit if not defined
				tmp := cast.ToString(auxNodeValue)
				if tmp == "0" || tmp == "" {
					break
				}
				refs["{{ $attr.expIdent }}"] = envoyx.Ref{
					ResourceType: "{{ $attr.dal.refModelResType }}",
					Identifiers: envoyx.MakeIdentifiers(auxNodeValue),
				}
			{{- end }}

			{{ if $attr.envoy.yaml.customDecoder }}
			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshal{{ $resource.expIdent }}{{ $attr.expIdent }}Node(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)
			{{ end }}
			break
		{{- end }}
	{{ if .rbac }}
		// Handle RBAC rules
		case "allow":
			auxOut, err = unmarshalAllowNode(n)
			if err != nil {
				return err
			}
			rbacNodes = append(rbacNodes, auxOut...)
			auxOut = nil

		case "deny":
			auxOut, err = unmarshalDenyNode(n)
			if err != nil {
				return err
			}
			rbacNodes = append(rbacNodes, auxOut...)
			auxOut = nil
	{{- end }}
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	{{ if eq $cmpIdent "compose" }}
	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers: envoyx.MakeIdentifiers(ref),
		}
	}
	{{ end }}

	{{- if and .envoy.scoped .parents}}
	// Define the scope
	//
	// This resource is scoped to the first parent (generally the namespace)
	// when talking about Compose resources (the only supported scenario at the moment).
	scope = envoyx.Scope{
		ResourceType: refs["{{(index .parents 0).refField}}"].ResourceType,
		Identifiers:  refs["{{(index .parents 0).refField}}"].Identifiers,
	}
	{{end}}
	{{if and .envoy.scoped (not .parents)}}
	// Define the scope
	//
	// This resource is scoped with no parent resources so this resource is the
	// root itself (generally the namespace -- the only currently supported scenario).
	scope = envoyx.Scope{
		ResourceType: types.{{ .expIdent }}ResourceType,
		Identifiers:  ii,
	}
	{{end}}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under {{.ident}} such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {
	{{/*
		@note each resource can only nest resources from the same component.
					Iterate resources of the current component and if this one appears as any
					of the parent resources, attempt to unmarshal it.

		@todo consider limiting the supported set to only the root parent
		      (like we do with store encoders) to limit strange cases and reduce
		      potential problems.
	*/}}
	{{ $a := . }}
	{{- range $cmp := $rootRes }}
		{{ if $cmp.envoy.omit }}
			{{continue}}
		{{ end }}

		{{- range $p := $cmp.parents }}
			{{ if not (eq $p.handle $a.ident) }}
				{{continue}}
			{{ end }}

		{{ $identKeys := $cmp.envoy.yaml.identKeys }}
		case {{ range $i, $l := $identKeys -}}
		"{{ $l }}"{{if not (eq $i (sub (len $identKeys) 1))}},{{end}}
		{{- end}}:
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshal{{$cmp.expIdent}}Seq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshal{{$cmp.expIdent}}Map(dctx, n)
				if err != nil {
					return err
				}
			}
			break

				{{/* As long as one parent matches we're golden; avoid potential duplicates */}}
			{{break}}
		{{- end }}
	{{- end -}}

	{{- range .envoy.yaml.extendedResourceDecoders }}
		{{ $identKeys := .identKeys }}
		case {{ range $i, $l := $identKeys -}}
		"{{ $l }}"{{if not (eq $i (sub (len $identKeys) 1))}},{{end}}
		{{- end}}:
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalExtended{{.expIdent}}Seq(dctx, n)
				if err != nil {
					return err
				}
			} {{- if .supportMappedInput }} else {
				nestedNodes, err = d.unmarshalExtended{{.expIdent}}Map(dctx, n)
				if err != nil {
					return err
				}
			}{{ end }}
			break
		{{ end -}}
		}

		// Iterate nested nodes and update their reference to the current resource
		//
		// Any reference to the parent resource from the child resource is overwritten
		// to avoid potential user-error edge cases.
		for _, a := range nestedNodes {
		  // @note all nested resources fall under the same component and the same scope.
			//       Simply assign the same scope to all -- if it shouldn't be scoped
			//       the parent won't have it (saving CPU ticks :)
			a.Scope = scope

			if a.References == nil {
				a.References = make(map[string]envoyx.Ref)
			}

			if _, ok := a.References["{{if .envoy.yaml.extendedResourceRefIdent}}{{.envoy.yaml.extendedResourceRefIdent}}{{else}}{{.expIdent}}ID{{end}}"]; !ok {
				a.References["{{if .envoy.yaml.extendedResourceRefIdent}}{{.envoy.yaml.extendedResourceRefIdent}}{{else}}{{.expIdent}}ID{{end}}"] = envoyx.Ref{
					ResourceType: types.{{ .expIdent }}ResourceType,
					Identifiers: ii,
					Scope: scope,
				}
			}

			for f, ref := range a.References {
				if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
					continue
				}
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
					continue
				}

				// Only assume refs if they're not yet set
				if _, ok := a.References[f]; ok {
					continue
				}

				a.References[f] = ref
			}
		}
		auxNestedNodes = append(auxNestedNodes, nestedNodes...)
		return nil
	})
	if err != nil {
		return
	}

	{{ if .envoy.yaml.extendedResourcePostProcess}}
	auxNestedNodes, err = d.postProcessNested{{.expIdent}}Nodes(auxNestedNodes)
	if err!= nil {
		return
	}
	{{ end }}
	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.{{ .expIdent }}ResourceType,
		Identifiers:  ii,
		References: refs,
		{{if or .envoy.scoped}}
		Scope: scope,
		{{end}}
		Config: envoyConfig,
	}
	{{- if .rbac }}
	// Update RBAC resource nodes with references regarding the resource
	rres := r.RbacResource()
	for _, rn := range rbacNodes {
		aux := rn.Resource.(*rbac.Rule)
		if aux.Resource == "" {
			aux.Resource = rres
		}

	  // Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["{{len .parents}}"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope: scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}
	{{ end }}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	{{- if .rbac }}
	out = append(out, rbacNodes...)
	{{ end }}

	return
}

{{ end }}

// // // // // // // // // // // // // // // // // // // // // // // // //
// RBAC unmarshal logic
// // // // // // // // // // // // // // // // // // // // // // // // //

func unmarshalAllowNode(n *yaml.Node) (out envoyx.NodeSet, err error) {
	return unmarshalRBACNode(n, rbac.Allow)
}

func unmarshalDenyNode(n *yaml.Node) (out envoyx.NodeSet, err error) {
	return unmarshalRBACNode(n, rbac.Deny)
}

func unmarshalRBACNode(n *yaml.Node, acc rbac.Access) (out envoyx.NodeSet, err error) {
	if y7s.IsMapping(n.Content[1]) {
		out, err = unmarshalFlatRBACNode(n, acc)
		if err != nil {
			return
		}
	} else {
		out, err = unmarshalNestedRBACNode(n, acc)
		if err != nil {
			return
		}
	}

	for _, o := range out {
		for _, r := range o.References {
			if r.Scope.IsEmpty() {
				continue
			}
			o.Scope = r.Scope
			break
		}
	}

	return
}

// unmarshalNestedRBACNode handles RBAC rules when they are nested inside a resource
//
// The edge-case exists since the node doesn't explicitly specify the resource
// it belongs to.
//
// Example:
//
// modules:
//   module1:
//     name: "module 1"
//     fields: ...
//     allow:
//       role1:
//         - read
//         - delete
func unmarshalNestedRBACNode(n *yaml.Node, acc rbac.Access) (out envoyx.NodeSet, err error) {
	return out, y7s.EachMap(n, func(role, op *yaml.Node) error {
		out = append(out, &envoyx.Node{
			Resource: &rbac.Rule{
				Operation: op.Value,
				Access:    acc,
			},
			ResourceType: rbac.RuleResourceType,
			References: map[string]envoyx.Ref{
				"RoleID": {
					// Providing resource type as plain text to reduce cross component references
					ResourceType: "corteza::system:role",
					Identifiers:  envoyx.MakeIdentifiers(role.Value),
				},
			},
		})
		return nil
	})
}

// unmarshalFlatRBACNode handles RBAC rules when they are provided on the root level
//
// Example:
//
// allow:
//   role1:
//     corteza::system/:
//       - users.search
//       - users.create
func unmarshalFlatRBACNode(n *yaml.Node, acc rbac.Access) (out envoyx.NodeSet, err error) {
	// Handles role
	return out, y7s.EachMap(n, func(role, perm *yaml.Node) error {
		// Handles operations
		return y7s.EachMap(perm, func(res, ops *yaml.Node) error {
			// Handle operation (one RBAC rule per op)
			return y7s.EachSeq(ops, func(op *yaml.Node) error {
				out = append(out, &envoyx.Node{
					Resource: &rbac.Rule{
						Resource:  res.Value,
						Operation: op.Value,
						Access:    acc,
					},
					ResourceType: rbac.RuleResourceType,
					References: envoyx.MergeRefs(
						map[string]envoyx.Ref{"RoleID": {
							// Providing resource type as plain text to reduce cross component references
							ResourceType: "corteza::system:role",
							Identifiers:  envoyx.MakeIdentifiers(role.Value),
						}},
						envoyx.SplitResourceIdentifier(res.Value),
					),
				})
				return nil
			})
		})
	})
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// i18n unmarshal logic
// // // // // // // // // // // // // // // // // // // // // // // // //

func unmarshalLocaleNode(n *yaml.Node) (out envoyx.NodeSet, err error) {
	err = y7s.EachMap(n, func(lang, loc *yaml.Node) error {
		langTag := systemTypes.Lang{Tag: language.Make(lang.Value)}

		return y7s.EachMap(loc, func(res, kv *yaml.Node) error {
			return y7s.EachMap(kv, func(k, msg *yaml.Node) error {
				out = append(out, &envoyx.Node{
					Resource: &systemTypes.ResourceTranslation{
						Resource: res.Value,
						Lang:     langTag,
						K:        k.Value,
						Message:  msg.Value,
					},
					// Providing resource type as plain text to reduce cross component references
					ResourceType: "corteza::system:resource-translation",
					References:   envoyx.SplitResourceIdentifier(res.Value),
				})
				return nil
			})
		})
	})
	if err != nil {
		return
	}

	for _, o := range out {
		for _, r := range o.References {
			if r.Scope.IsEmpty() {
				continue
			}
			o.Scope = r.Scope
			break
		}
	}

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Envoy config unmarshal logic
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d *auxYamlDoc) decodeEnvoyConfig(n *yaml.Node) (out envoyx.EnvoyConfig) {
	y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "skipif", "skip":
			return y7s.DecodeScalar(v, "decode skip if", &out.SkipIf)
		case "onexisting", "mergealg":
			out.MergeAlg = envoyx.CastMergeAlg(v.Value)
		}

		return nil
	})

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Utilities
// // // // // // // // // // // // // // // // // // // // // // // // //

func (d YamlDecoder) getReader(ctx context.Context, p envoyx.DecodeParams) (r io.Reader, err error) {
	aux, ok := p.Params[paramsKeyStream]
	if ok {
		r, ok = aux.(io.Reader)
		if ok {
			return
		}
	}

	// @todo consider adding support for managing files from a location
	err = errors.Errorf("YAML decoder expects a stream conforming to io.Reader interface")
	return
}

func (d *auxYamlDoc) getDocumentContext(n *yaml.Node) (dctx documentContext, err error) {
	dctx = documentContext{
		references: make(map[string]string),
	}

	err = y7s.EachMap(n, func(k, v *yaml.Node) error {
		// @todo expand when needed. The previous implementation only supported
		//       namespaces on the root of the document.

		if y7s.IsKind(v, yaml.ScalarNode) {
			dctx.references[k.Value] = v.Value
		}

		return nil
	})

	return
}
