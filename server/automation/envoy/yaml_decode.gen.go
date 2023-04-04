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

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
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
		err = errors.Wrap(err, "automation yaml decoder: failed to decode document")
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
		// Decode all resources under the automation component
		case "workflow", "workflows":
			if y7s.IsMapping(v) {
				aux, err = d.unmarshalWorkflowMap(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}

			if y7s.IsSeq(v) {
				aux, err = d.unmarshalWorkflowSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: workflow")
			}
			return err

		case "trigger":
			// @note trigger doesn't support mapped inputs. This can be
			//       changed in the .cue definition under the
			//       .envoy.yaml.supportMappedInput field.
			if y7s.IsSeq(v) {
				aux, err = d.unmarshalTriggerSeq(dctx, v)
				d.nodes = append(d.nodes, aux...)
			}
			if err != nil {
				err = errors.Wrap(err, "failed to unmarshal node: trigger")
			}
			return err

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
} // // // // // // // // // // // // // // // // // // // // // // // // //
// Functions for resource workflow
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalWorkflowSeq unmarshals Workflow when provided as a sequence node
func (d *auxYamlDoc) unmarshalWorkflowSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalWorkflowNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalWorkflowMap unmarshals Workflow when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
func (d *auxYamlDoc) unmarshalWorkflowMap(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		aux, err = d.unmarshalWorkflowNode(dctx, n, k)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalTriggersExtendedSeq unmarshals Triggers when provided as a sequence node
func (d *auxYamlDoc) unmarshalExtendedTriggersSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalTriggersExtendedNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalWorkflowNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalWorkflowNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Workflow

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

		case "runas":
			// Handle references
			err = y7s.DecodeScalar(n, "runAs", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["RunAs"] = envoyx.Ref{
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

	// Handle any resources that could be inserted under workflow such as a module inside a namespace
	//
	// This operation is done in the second pass of the document so we have
	// the complete context of the current resource; such as the identifier,
	// references, and scope.
	var auxNestedNodes envoyx.NodeSet
	err = y7s.EachMap(n, func(k, n *yaml.Node) error {
		nestedNodes = nil

		switch strings.ToLower(k.Value) {

		case "triggers":
			if y7s.IsSeq(n) {
				nestedNodes, err = d.unmarshalExtendedTriggersSeq(dctx, n)
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

			if _, ok := a.References["WorkflowID"]; !ok {
				a.References["WorkflowID"] = envoyx.Ref{
					ResourceType: types.WorkflowResourceType,
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

		ResourceType: types.WorkflowResourceType,
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
// Functions for resource trigger
// // // // // // // // // // // // // // // // // // // // // // // // //

// unmarshalTriggerSeq unmarshals Trigger when provided as a sequence node
func (d *auxYamlDoc) unmarshalTriggerSeq(dctx documentContext, n *yaml.Node) (out envoyx.NodeSet, err error) {
	var aux envoyx.NodeSet
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		aux, err = d.unmarshalTriggerNode(dctx, n)
		if err != nil {
			return err
		}
		out = append(out, aux...)

		return nil
	})

	return
}

// unmarshalTriggerMap unmarshals Trigger when provided as a mapping node
//
// When map encoded, the map key is used as a preset identifier.
// The identifier is passed to the node function as a meta node
// @note this resource does not support map encoding.
//       Refer to the corresponding definition files to adjust if needed.

// unmarshalTriggerNode is a cookie-cutter function to unmarshal
// the yaml node into the corresponding Corteza type & Node
func (d *auxYamlDoc) unmarshalTriggerNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	var r *types.Trigger

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

		case "eventtype", "eventType", "event_type":
			// Handle field alias
			//
			// @todo consider adding an is empty check before overwriting
			err = y7s.DecodeScalar(n, "eventType", &r.EventType)
			if err != nil {
				return err
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

		case "resourcetype", "resourceType", "resource_type":
			// Handle field alias
			//
			// @todo consider adding an is empty check before overwriting
			err = y7s.DecodeScalar(n, "resourceType", &r.ResourceType)
			if err != nil {
				return err
			}

			break

		case "stepid", "stepID", "step_id":
			// Handle field alias
			//
			// @todo consider adding an is empty check before overwriting
			err = y7s.DecodeScalar(n, "stepID", &r.StepID)
			if err != nil {
				return err
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

		case "workflowid":
			// Handle references
			err = y7s.DecodeScalar(n, "workflowID", &auxNodeValue)
			if err != nil {
				return err
			}

			// Omit if not defined
			tmp := cast.ToString(auxNodeValue)
			if tmp == "0" || tmp == "" {
				break
			}
			refs["WorkflowID"] = envoyx.Ref{
				ResourceType: "corteza::automation:workflow",
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

	// Handle any resources that could be inserted under trigger such as a module inside a namespace
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

			if _, ok := a.References["TriggerID"]; !ok {
				a.References["TriggerID"] = envoyx.Ref{
					ResourceType: types.TriggerResourceType,
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

		ResourceType: types.TriggerResourceType,
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
