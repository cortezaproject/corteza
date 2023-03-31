package envoy

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	systemTypes "github.com/cortezaproject/corteza/server/system/types"
	sqlt "github.com/jmoiron/sqlx/types"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func (d *auxYamlDoc) unmarshalYAML(k string, n *yaml.Node) (out envoyx.NodeSet, err error) {
	switch k {
	case "settings", "setting":
		out, err = d.unmarshalSettingsNode(n)
	}

	return
}

func (d *auxYamlDoc) unmarshalSettingsNode(n *yaml.Node) (out envoyx.NodeSet, err error) {
	out = make(envoyx.NodeSet, 0, len(n.Content))

	err = y7s.Each(n, func(k, v *yaml.Node) (err error) {
		if v == nil {
			return y7s.NodeErr(n, "malformed setting definition")
		}

		s := &types.SettingValue{}
		switch v.Kind {
		case yaml.MappingNode:
			if err = v.Decode(&s); err != nil {
				return
			}

		default:
			if err = y7s.DecodeScalar(k, "setting name", &s.Name); err != nil {
				return err
			}

			if y7s.IsKind(v, yaml.SequenceNode) {
				aux := make([]interface{}, 0, 10)

				y7s.EachSeq(v, func(n *yaml.Node) error {
					var vx interface{}
					err := y7s.DecodeScalar(n, "setting value", &vx)
					if err != nil {
						return err
					}

					aux = append(aux, vx)
					return nil
				})

				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				s.Value = sqlt.JSONText(m)
			} else if y7s.IsKind(v, yaml.MappingNode) {
				aux := make(map[string]interface{})

				y7s.EachMap(v, func(k, v *yaml.Node) error {
					var vx interface{}
					err := y7s.DecodeScalar(v, "setting value", &vx)
					if err != nil {
						return err
					}

					aux[k.Value] = vx
					return nil
				})

				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				s.Value = sqlt.JSONText(m)
			} else {
				var aux interface{}
				err = y7s.DecodeScalar(v, "setting value", &aux)
				if err != nil {
					return err
				}
				m, err := json.Marshal(aux)
				if err != nil {
					return err
				}
				s.Value = sqlt.JSONText(m)
			}
		}

		out = append(out, &envoyx.Node{
			Resource:     s,
			ResourceType: types.SettingValueResourceType,
		})
		return
	})

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

func (d *auxYamlDoc) unmarshalFiltersExtendedNode(dctx documentContext, n *yaml.Node, meta ...*yaml.Node) (out envoyx.NodeSet, err error) {
	return d.unmarshalApigwFilterNode(dctx, n, meta...)
}

func (d *auxYamlDoc) unmarshalUserRolesNode(r *types.User, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref, len(n.Content))

	i := 0
	err = y7s.EachSeq(n, func(n *yaml.Node) error {
		refs[fmt.Sprintf("Roles.%d", i)] = envoyx.Ref{
			ResourceType: types.RoleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(n.Value),
		}

		return nil
	})

	return
}

func (d *auxYamlDoc) unmarshalAuthClientSecurityNode(r *types.AuthClient, n *yaml.Node) (refs map[string]envoyx.Ref, idents envoyx.Identifiers, err error) {
	refs = make(map[string]envoyx.Ref)

	err = y7s.EachMap(n, func(k, v *yaml.Node) (err error) {
		switch strings.ToLower(k.Value) {
		case "impersonateuser":
			var av string
			err = y7s.DecodeScalar(v, "Impersonate user", &av)

			refs["Security.ImpersonateUser"] = envoyx.Ref{
				ResourceType: types.UserResourceType,
				Identifiers:  envoyx.MakeIdentifiers(av),
			}
			break

		case "permittedroles":
			refs = envoyx.MergeRefs(refs, roleNodeSliceToRefs("Security.PermittedRoles", v))
			break

		case "prohibitedroles":
			refs = envoyx.MergeRefs(refs, roleNodeSliceToRefs("Security.ProhibitedRoles", v))
			break

		case "forcedroles":
			refs = envoyx.MergeRefs(refs, roleNodeSliceToRefs("Security.ForcedRoles", v))
			break

		}

		return nil
	})

	return
}

func roleNodeSliceToRefs(k string, n *yaml.Node) (refs map[string]envoyx.Ref) {
	refs = make(map[string]envoyx.Ref, len(n.Content))

	i := 0
	y7s.EachSeq(n, func(n *yaml.Node) error {
		refs[fmt.Sprintf("%s.%d", k, i)] = envoyx.Ref{
			ResourceType: types.RoleResourceType,
			Identifiers:  envoyx.MakeIdentifiers(n.Value),
		}
		i++
		return nil
	})

	return
}
