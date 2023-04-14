package envoy

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

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
	auxYamlDoc struct {
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
		err = errors.Wrap(err, "system yaml decoder: failed to decode document")
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
		// Decode all resources under the system component
		case "application", "applications", "apps":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalApplicationMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalApplicationSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: application")
			}
			return err

		case "apigwroute", "endpoints":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalApigwRouteMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalApigwRouteSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: apigwRoute")
			}
			return err

		case "apigwfilter":
			// @note apigwFilter doesn't support mapped inputs. This can be
			//       changed in the .cue definition under the
			//       .envoy.yaml.supportMappedInput field.
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalApigwFilterSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: apigwFilter")
			}
			return err

		case "authclient", "authclients":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalAuthClientMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalAuthClientSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: authClient")
			}
			return err

		case "queue":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalQueueMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalQueueSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: queue")
			}
			return err

		case "report", "reports":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalReportMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalReportSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: report")
			}
			return err

		case "role", "roles":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalRoleMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalRoleSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: role")
			}
			return err

		case "template", "templates":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalTemplateMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalTemplateSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: template")
			}
			return err

		case "user", "users", "usr":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalUserMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalUserSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: user")
			}
			return err

		case "dalconnection", "connection", "connections":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalDalConnectionMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalDalConnectionSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: dalConnection")
			}
			return err

		case "dalsensitivitylevel", "sensitivity_level":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalDalSensitivityLevelMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalDalSensitivityLevelSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: dalSensitivityLevel")
			}
			return err

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

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource application
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalApplicationSeq unmarshals Application when provided as a sequence node
func (d *auxYamlDoc) unmarshalApplicationSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalApplicationNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalApplicationMap unmarshals Application when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalApplicationMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalApplicationNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalApplicationNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalApplicationNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Application

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

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "ownerid", "owner":
			// Handle references
			err = y7s.DecodeScalar(n, "ownerID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["OwnerID"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under application such as a module inside a namespace
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

			if _, ok := a.References["ApplicationID"]; !ok {
				a.References["ApplicationID"] = envoyx.Ref{
					ResourceType: types.ApplicationResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ApplicationResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource apigwRoute
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalApigwRouteSeq unmarshals ApigwRoute when provided as a sequence node
func (d *auxYamlDoc) unmarshalApigwRouteSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalApigwRouteNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalApigwRouteMap unmarshals ApigwRoute when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalApigwRouteMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalApigwRouteNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalFiltersExtendedSeq unmarshals Filters when provided as a sequence node
func (d *auxYamlDoc) unmarshalExtendedFiltersSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalFiltersExtendedNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalApigwRouteNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalApigwRouteNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.ApigwRoute

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
		y7s.DecodeScalar(meta[0], "Endpoint", &r.Endpoint)
		ii = ii.Add(r.Endpoint)
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

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "group":
			// Handle references
			err = y7s.DecodeScalar(n, "group", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["Group"] = envoyx.Ref{
				ResourceType: "corteza::system:apigw-group",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under apigwRoute such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

		case "filters":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalExtendedFiltersSeq(dctx, n)
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

			if _, ok := a.References["ApigwRouteID"]; !ok {
				a.References["ApigwRouteID"] = envoyx.Ref{
					ResourceType: types.ApigwRouteResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ApigwRouteResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource apigwFilter
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalApigwFilterSeq unmarshals ApigwFilter when provided as a sequence node
func (d *auxYamlDoc) unmarshalApigwFilterSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalApigwFilterNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalApigwFilterMap unmarshals ApigwFilter when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
// @note this resource does not support map encoding.
//       Refer to the corresponding definition files to adjust if needed.

// unmarshalApigwFilterNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalApigwFilterNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.ApigwFilter

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

	var (
		refs        = make(map[string]envoyx.Ref)
		auxOut      envoyx.NodeSet
		nestedNodes envoyx.NodeSet
		scope       envoyx.Scope
		envoyConfig envoyx.EnvoyConfig
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "route":
			// Handle references
			err = y7s.DecodeScalar(n, "route", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["Route"] = envoyx.Ref{
				ResourceType: "corteza::system:apigw-route",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

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

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under apigwFilter such as a module inside a namespace
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

			if _, ok := a.References["ApigwFilterID"]; !ok {
				a.References["ApigwFilterID"] = envoyx.Ref{
					ResourceType: types.ApigwFilterResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ApigwFilterResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource authClient
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalAuthClientSeq unmarshals AuthClient when provided as a sequence node
func (d *auxYamlDoc) unmarshalAuthClientSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalAuthClientNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalAuthClientMap unmarshals AuthClient when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalAuthClientMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalAuthClientNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalAuthClientNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalAuthClientNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.AuthClient

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

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

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

		case "ownedby":
			// Handle references
			err = y7s.DecodeScalar(n, "ownedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "security":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalAuthClientSecurityNode(r, n)
			if err != nil {
				return err
			}
			refs = envoyx.MergeRefs(refs, auxRefs)
			ii = ii.Merge(auxIdents)

			break

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under authClient such as a module inside a namespace
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

			if _, ok := a.References["AuthClientID"]; !ok {
				a.References["AuthClientID"] = envoyx.Ref{
					ResourceType: types.AuthClientResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.AuthClientResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource queue
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalQueueSeq unmarshals Queue when provided as a sequence node
func (d *auxYamlDoc) unmarshalQueueSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalQueueNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalQueueMap unmarshals Queue when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalQueueMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalQueueNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalQueueNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalQueueNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Queue

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
		y7s.DecodeScalar(meta[0], "Queue", &r.Queue)
		ii = ii.Add(r.Queue)
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

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "id":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "id", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "queue":
			// Handle identifiers
			err = y7s.DecodeScalar(n, "queue", &auxNodeValue)
			if err != nil {
				return err
			}
			ii = ii.Add(auxNodeValue)

			break

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under queue such as a module inside a namespace
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

			if _, ok := a.References["QueueID"]; !ok {
				a.References["QueueID"] = envoyx.Ref{
					ResourceType: types.QueueResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.QueueResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource report
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalReportSeq unmarshals Report when provided as a sequence node
func (d *auxYamlDoc) unmarshalReportSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalReportNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalReportMap unmarshals Report when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalReportMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalReportNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalReportNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalReportNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Report

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

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

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

		case "ownedby":
			// Handle references
			err = y7s.DecodeScalar(n, "ownedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["OwnedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under report such as a module inside a namespace
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

			if _, ok := a.References["ReportID"]; !ok {
				a.References["ReportID"] = envoyx.Ref{
					ResourceType: types.ReportResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.ReportResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource role
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalRoleSeq unmarshals Role when provided as a sequence node
func (d *auxYamlDoc) unmarshalRoleSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalRoleNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalRoleMap unmarshals Role when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalRoleMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalRoleNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalRoleNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalRoleNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Role

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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under role such as a module inside a namespace
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

			if _, ok := a.References["RoleID"]; !ok {
				a.References["RoleID"] = envoyx.Ref{
					ResourceType: types.RoleResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.RoleResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource template
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalTemplateSeq unmarshals Template when provided as a sequence node
func (d *auxYamlDoc) unmarshalTemplateSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalTemplateNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalTemplateMap unmarshals Template when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalTemplateMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalTemplateNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalTemplateNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalTemplateNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Template

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

		case "ownerid":
			// Handle references
			err = y7s.DecodeScalar(n, "ownerID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["OwnerID"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under template such as a module inside a namespace
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

			if _, ok := a.References["TemplateID"]; !ok {
				a.References["TemplateID"] = envoyx.Ref{
					ResourceType: types.TemplateResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.TemplateResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource user
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalUserSeq unmarshals User when provided as a sequence node
func (d *auxYamlDoc) unmarshalUserSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalUserNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalUserMap unmarshals User when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalUserMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalUserNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalUserNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalUserNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.User

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

		case "roles":

			// Handle custom node decoder
			//
			// The decoder may update the passed resource with arbitrary values
			// as well as provide additional references and identifiers for the node.
			var (
				auxRefs   map[string]envoyx.Ref
				auxIdents envoyx.Identifiers
			)
			auxRefs, auxIdents, err = d.unmarshalUserRolesNode(r, n)
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under user such as a module inside a namespace
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

			if _, ok := a.References["UserID"]; !ok {
				a.References["UserID"] = envoyx.Ref{
					ResourceType: types.UserResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.UserResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource dalConnection
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalDalConnectionSeq unmarshals DalConnection when provided as a sequence node
func (d *auxYamlDoc) unmarshalDalConnectionSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalDalConnectionNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalDalConnectionMap unmarshals DalConnection when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalDalConnectionMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalDalConnectionNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalDalConnectionNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalDalConnectionNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.DalConnection

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

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

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

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
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

	// Make parent identifiers available through the dctx
	dctx.parentIdent = ii

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under dalConnection such as a module inside a namespace
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

			if _, ok := a.References["DalConnectionID"]; !ok {
				a.References["DalConnectionID"] = envoyx.Ref{
					ResourceType: types.DalConnectionResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.DalConnectionResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}
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
		rn.References["0"] = envoyx.Ref{
			ResourceType: a.ResourceType,
			Identifiers:  a.Identifiers,
			Scope:        scope,
		}

		for _, r := range rn.References {
			if r.Scope.IsEmpty() {
				continue
			}
			rn.Scope = r.Scope
			break
		}
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)
	out = append(out, rbacNodes...)

	return
}

// // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource dalSensitivityLevel
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalDalSensitivityLevelSeq unmarshals DalSensitivityLevel when provided as a sequence node
func (d *auxYamlDoc) unmarshalDalSensitivityLevelSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalDalSensitivityLevelNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalDalSensitivityLevelMap unmarshals DalSensitivityLevel when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalDalSensitivityLevelMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalDalSensitivityLevelNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalDalSensitivityLevelNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalDalSensitivityLevelNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.DalSensitivityLevel

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
	)
	_ = auxOut
	_ = refs

	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		var auxNodeValue any
		_ = auxNodeValue

		switch strings.ToLower(k.Value) {

		case "createdby":
			// Handle references
			err = y7s.DecodeScalar(n, "createdBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["CreatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

		case "deletedby":
			// Handle references
			err = y7s.DecodeScalar(n, "deletedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["DeletedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

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

		case "updatedby":
			// Handle references
			err = y7s.DecodeScalar(n, "updatedBy", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["UpdatedBy"] = envoyx.Ref{
				ResourceType: "corteza::system:user",
				Identifiers:  envoyx.MakeIdentifiers(auxNodeValue),
			}

			break

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

	// Apply the scope to all of the references of the same type
	for k, ref := range refs {
		if !strings.HasPrefix(ref.ResourceType, "corteza::compose") {
			continue
		}
		ref.Scope = scope
		refs[k] = ref
	}

	// Handle any resources that could be inserted under dalSensitivityLevel such as a module inside a namespace
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

			if _, ok := a.References["DalSensitivityLevelID"]; !ok {
				a.References["DalSensitivityLevelID"] = envoyx.Ref{
					ResourceType: types.DalSensitivityLevelResourceType,
					Identifiers:  ii,
					Scope:        scope,
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

	out = append(out, auxNestedNodes...)

	a := &envoyx.Node{
		Resource: r,

		ResourceType: types.DalSensitivityLevelResourceType,
		Identifiers:  ii,
		References:   refs,

		Config: envoyConfig,
	}

	// Put it all together...
	out = append(out, a)
	out = append(out, auxOut...)

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
