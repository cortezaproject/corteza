package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/spf13/cast"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type (
	// YamlDecoder is responsible for decoding YAML documents into Corteza resources
	// which are then managed by envoy and imported via an encoder.
	YamlDecoder     struct{}
	documentContext struct {
		references map[string]string
	}
	auxYamlDoc struct {
		nodes envoyx.NodeSet
	}
)

// Decode returns a set of envoy nodes based on the provided params
//
// YamlDecoder expects the DecodeParam of `stream` which conforms
// to the io.Reader interface.
func (d YamlDecoder) Decode(ctx context.Context, p envoyx.DecodeParams) (out envoyx.NodeSet, err error) {
	// Get the reader
	r, err := d.getReader(ctx, p)
	if err != nil {
		return
	}

	// Offload decoding to the aux document
	doc := &auxYamlDoc{}
	err = yaml.NewDecoder(r).Decode(doc)
	if err != nil {
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
		case "chart", "charts", "chrt":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalChartMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
				return err
			}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalChartSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			return err

		case "module", "modules", "mod":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalModuleMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
				return err
			}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalModuleSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			return err

		case "modulefield", "module_fields", "modulefields", "fields":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalModuleFieldMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
				return err
			}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalModuleFieldSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			return err

		case "namespace", "namespaces", "ns":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalNamespaceMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
				return err
			}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalNamespaceSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			return err

		case "page", "pages", "pg":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalPageMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
				return err
			}
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalPageSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			return err

		// Offload to custom handlers
		default:
			aux, err = d.unmarshalYAML(kv, v)
			d.nodes = append(d.nodes, aux...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource chart
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalChartSeq unmarshals Chart when provided as a sequence node
func (d *auxYamlDoc) unmarshalChartSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalChartNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalChartMap unmarshals Chart when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalChartMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalChartNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalChartNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalChartNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Chart

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

	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "Handle", &r.Handle)
		ii = ii.Add(r.Handle)
	}

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		rbacNodes   envoyx.NodeSet
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "config":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalChartConfigNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

		case "handle":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "handle", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "namespaceid", "namespace", "namespace_id", "ns":
			// Handle references
			err = y7s.DecodeScalar(n, "namespaceID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["NamespaceID"] = envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

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
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(ref),
		}
	}

	// Define the scope
	//
	// This resource is scoped to the first parent (generally the namespace)
	// when talking about Compose resources (the only supported scenario at the moment).
	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if ref.ResourceType != scope.ResourceType {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under chart such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

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

			a.References["ChartID"] = envoyx.Ref{
				ResourceType: types.ChartResourceType,
				Identifiers:  ii,
				Scope:        scope,
			}

			for f, ref := range a.References {
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ChartResourceType,
		Identifiers:  ii,
		References:   refs,

		Scope: scope,

		Config: envoyConfig,
	}
	// Update RBAC resource nodes with references regarding the resource
	for _, rn := range rbacNodes {
		// Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["1"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource module
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalModuleSeq unmarshals Module when provided as a sequence node
func (d *auxYamlDoc) unmarshalModuleSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalModuleNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalModuleMap unmarshals Module when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalModuleMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalModuleNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalSourceExtendedSeq unmarshals Source when provided as a sequence node
func (d *auxYamlDoc) unmarshalExtendedSourceSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalSourceExtendedNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalModuleNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalModuleNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Module

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

	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "Handle", &r.Handle)
		ii = ii.Add(r.Handle)
	}

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		rbacNodes   envoyx.NodeSet
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "handle":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "handle", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "namespaceid", "namespace", "namespace_id", "ns", "ns_id":
			// Handle references
			err = y7s.DecodeScalar(n, "namespaceID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["NamespaceID"] = envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

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
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(ref),
		}
	}

	// Define the scope
	//
	// This resource is scoped to the first parent (generally the namespace)
	// when talking about Compose resources (the only supported scenario at the moment).
	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if ref.ResourceType != scope.ResourceType {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under module such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

		case "modulefield", "module_fields", "modulefields", "fields":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalModuleFieldSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalModuleFieldMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break

		case "source", "datasource", "records":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalExtendedSourceSeq(dctx, n)
				if err != nil {
					return err
				}
			}
			break
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

			a.References["ModuleID"] = envoyx.Ref{
				ResourceType: types.ModuleResourceType,
				Identifiers:  ii,
				Scope:        scope,
			}

			for f, ref := range a.References {
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
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

	auxNestedNodes, err = d.postProcessNestedModuleNodes(auxNestedNodes)
	if err != nil {
		return
	}

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ModuleResourceType,
		Identifiers:  ii,
		References:   refs,

		Scope: scope,

		Config: envoyConfig,
	}
	// Update RBAC resource nodes with references regarding the resource
	for _, rn := range rbacNodes {
		// Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["1"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource moduleField
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalModuleFieldSeq unmarshals ModuleField when provided as a sequence node
func (d *auxYamlDoc) unmarshalModuleFieldSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalModuleFieldNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalModuleFieldMap unmarshals ModuleField when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalModuleFieldMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalModuleFieldNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalModuleFieldNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalModuleFieldNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.ModuleField

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

	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "Name", &r.Name)
		ii = ii.Add(r.Name)
	}

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		rbacNodes   envoyx.NodeSet
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "defaultvalue":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalModuleFieldDefaultValueNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

		case "expressions":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalModuleFieldExpressionsNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "moduleid":
			// Handle references
			err = y7s.DecodeScalar(n, "moduleID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["ModuleID"] = envoyx.Ref{
				ResourceType: "corteza::compose:module",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "name":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "name", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "options":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalModuleFieldOptionsNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

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
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(ref),
		}
	}

	// Define the scope
	//
	// This resource is scoped to the first parent (generally the namespace)
	// when talking about Compose resources (the only supported scenario at the moment).
	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if ref.ResourceType != scope.ResourceType {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under moduleField such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

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

			a.References["ModuleFieldID"] = envoyx.Ref{
				ResourceType: types.ModuleFieldResourceType,
				Identifiers:  ii,
				Scope:        scope,
			}

			for f, ref := range a.References {
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ModuleFieldResourceType,
		Identifiers:  ii,
		References:   refs,

		Scope: scope,

		Config: envoyConfig,
	}
	// Update RBAC resource nodes with references regarding the resource
	for _, rn := range rbacNodes {
		// Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["2"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource namespace
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalNamespaceSeq unmarshals Namespace when provided as a sequence node
func (d *auxYamlDoc) unmarshalNamespaceSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalNamespaceNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalNamespaceMap unmarshals Namespace when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalNamespaceMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalNamespaceNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalNamespaceNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalNamespaceNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Namespace

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

	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "Slug", &r.Slug)
		ii = ii.Add(r.Slug)
	}

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		rbacNodes   envoyx.NodeSet
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "slug":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "slug", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

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
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(ref),
		}
	}

	// Define the scope
	//
	// This resource is scoped with no parent resources so this resource is the
	// root itself (generally the namespace -- the only currently supported scenario).
	scope = envoyx.Scope{
		ResourceType: types.NamespaceResourceType,
		Identifiers:  ii,
	}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if ref.ResourceType != scope.ResourceType {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under namespace such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

		case "chart", "charts", "chrt":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalChartSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalChartMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break

		case "module", "modules", "mod":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalModuleSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalModuleMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break

		case "modulefield", "module_fields", "modulefields", "fields":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalModuleFieldSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalModuleFieldMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break

		case "page", "pages", "pg":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalPageSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalPageMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break

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

			a.References["NamespaceID"] = envoyx.Ref{
				ResourceType: types.NamespaceResourceType,
				Identifiers:  ii,
				Scope:        scope,
			}

			for f, ref := range a.References {
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.NamespaceResourceType,
		Identifiers:  ii,
		References:   refs,

		Scope: scope,

		Config: envoyConfig,
	}
	// Update RBAC resource nodes with references regarding the resource
	for _, rn := range rbacNodes {
		// Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource page
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalPageSeq unmarshals Page when provided as a sequence node
func (d *auxYamlDoc) unmarshalPageSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalPageNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalPageMap unmarshals Page when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalPageMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalPageNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalPagesExtendedSeq unmarshals Pages when provided as a sequence node
func (d *auxYamlDoc) unmarshalExtendedPagesSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalPagesExtendedNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalPagesExtendedMap unmarshals Pages when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalExtendedPagesMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalPagesExtendedNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalPageNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalPageNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Page

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

	// When a resource supports mapped input, the key is passed as meta which
	// needs to be registered as an identifier (since it is)
	if len(meta) > 0 {
		y7s.DecodeScalar(meta[0], "Handle", &r.Handle)
		ii = ii.Add(r.Handle)
	}

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
		rbacNodes   envoyx.NodeSet
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "blocks":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalPageBlocksNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

		case "handle":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "handle", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "moduleid", "module":
			// Handle references
			err = y7s.DecodeScalar(n, "moduleID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["ModuleID"] = envoyx.Ref{
				ResourceType: "corteza::compose:module",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "namespaceid", "namespace":
			// Handle references
			err = y7s.DecodeScalar(n, "namespaceID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["NamespaceID"] = envoyx.Ref{
				ResourceType: "corteza::compose:namespace",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "selfid", "parent":
			// Handle references
			err = y7s.DecodeScalar(n, "selfID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["SelfID"] = envoyx.Ref{
				ResourceType: "corteza::compose:page",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "weight", "order":
			// Handle field alias
			//
			// @todo consider adding an is empty check before overwriting
			err = y7s.DecodeScalar(n, "weight", &r.Weight)
			if err != nil {
				return err
			}

			break

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
		case "(envoy)":
			envoyConfig = d.decodeEnvoyConfig(n)
		}

		return nil
	})
	if err != nil {
		return
	}

	// Handle global namespace reference which can be provided as the doc. context
	//
	// @todo this is a temporary solution and should be extended when the document
	//       context needs to be extended.
	//       Limit this only to the compose resource since that is the only scenario
	//       the previous implementation supports.
	if ref, ok := dctx.references["namespace"]; ok {
		refs["NamespaceID"] = envoyx.Ref{
			ResourceType: types.NamespaceResourceType,
			Identifiers:  envoyx.MakeIdentifiers(ref),
		}
	}

	// Define the scope
	//
	// This resource is scoped to the first parent (generally the namespace)
	// when talking about Compose resources (the only supported scenario at the moment).
	scope = envoyx.Scope{
		ResourceType: refs["NamespaceID"].ResourceType,
		Identifiers:  refs["NamespaceID"].Identifiers,
	}

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if ref.ResourceType != scope.ResourceType {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under page such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

		case "children", "pages":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalExtendedPagesSeq(dctx, n)
				if err != nil {
					return err
				}
			} else {
				nestedNodes, err = d.unmarshalExtendedPagesMap(dctx, n)
				if err != nil {
					return err
				}
			}
			break
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

			a.References["PageID"] = envoyx.Ref{
				ResourceType: types.PageResourceType,
				Identifiers:  ii,
				Scope:        scope,
			}

			for f, ref := range a.References {
				ref.Scope = scope
				a.References[f] = ref
			}

			for f, ref := range refs {
				// Only inherit root references
				// @todo improve; this is a hack
				if strings.Contains(f, ".") {
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.PageResourceType,
		Identifiers:  ii,
		References:   refs,

		Scope: scope,

		Config: envoyConfig,
	}
	// Update RBAC resource nodes with references regarding the resource
	for _, rn := range rbacNodes {
		// Since the rule belongs to the resource, it will have the same
		// subset of references as the parent resource.
		rn.References = envoyx.MergeRefs(rn.References, a.References)

		// The RBAC rule's most specific identifier is the resource itself.
		// Using this we can hardcode it to point to the location after the parent resource.
		//
		// @todo consider using a more descriptive identifier for the position
		//       such as `index-%d`.
		rn.References["1"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

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
		out, err = unmarshalNestedRBACNode(n, acc)
		if err != nil {
			return
		}
	} else {
		out, err = unmarshalFlatRBACNode(n, acc)
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

// // // // // // // // // // // // // // // // // // // // // // // // //
// i18n unmarshal logic
// // // // // // // // // // // // // // // // // // // // // // // // //

func unmarshalLocaleNode(n *yaml.Node) (out envoyx.NodeSet, err error) {
	return out, y7s.EachMap(n, func(lang, loc *yaml.Node) error {
		langTag := systemTypes.Lang{Tag: language.Make(lang.Value)}

		return y7s.EachMap(loc, func(res, kv *yaml.Node) error {
			return y7s.EachMap(kv, func(k, msg *yaml.Node) error {
				out = append(out, &envoyx.Node{
					Resource: &systemTypes.ResourceTranslation{
						Lang:    langTag,
						K:       k.Value,
						Message: msg.Value,
					},
					// Providing resource type as plain text to reduce cross component references
					ResourceType: "corteza::system:resource-translation",
					References:   envoyx.SplitResourceIdentifier(res.Value),
				})
				return nil
			})
		})
	})
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
	aux, ok := p.Params["stream"]
	if ok {
		r, ok = aux.(io.Reader)
		if ok {
			return
		}
	}

	// @todo consider adding support for managing files from a location
	err = fmt.Errorf("YAML decoder expects a stream conforming to io.Reader interface")
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
