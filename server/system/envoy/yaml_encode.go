package envoy

import (
	"context"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"github.com/cortezaproject/corteza/server/system/types"
	"gopkg.in/yaml.v3"
)

type (
	// Wrapper to avoid constant type casting for role nodes
	rbacRuleWrap struct {
		node *envoyx.Node
		rule *rbac.Rule
	}

	// Wrapper to keep role and rules together
	rbacRuleRoleWrap struct {
		rules    []rbacRuleWrap
		roleNode *envoyx.Node
		role     *types.Role
	}
)

func (e YamlEncoder) encode(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, rt string, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	switch rt {
	case rbac.RuleResourceType:
		return e.encodeRbacRules(ctx, base, p, nodes, tt)
	case types.ResourceTranslationResourceType:
		return e.encodeResourceTranslations(ctx, base, p, nodes, tt)
	case types.SettingValueResourceType:
		return e.encodeSettingValues(ctx, base, p, nodes, tt)
	}

	return
}

func (e YamlEncoder) encodeResourceTranslations(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	err = e.resolveRulePathDeps(ctx, tt, nodes)
	if err != nil {
		return
	}

	byLang := make(map[string][]*envoyx.Node)

	for _, n := range nodes {
		byLang[n.Resource.(*types.ResourceTranslation).Lang.String()] = append(byLang[n.Resource.(*types.ResourceTranslation).Lang.String()], n)
	}

	out = base

	var aux *yaml.Node
	for lang, nodes := range byLang {
		aux, err = e.encodeResourceTranslationsByResource(p, nodes, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddMap(out, lang, aux)
		if err != nil {
			return
		}
	}

	return y7s.MakeMap("locale", out)
}

func (e YamlEncoder) encodeResourceTranslationsByResource(p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	byResource := make(map[string][]*envoyx.Node)

	for _, n := range nodes {
		byResource[n.Resource.(*types.ResourceTranslation).Resource] = append(byResource[n.Resource.(*types.ResourceTranslation).Resource], n)
	}

	var aux *yaml.Node
	for resource, nodes := range byResource {
		aux, err = e.encodeResourceTranslationsKv(p, nodes, tt)
		if err != nil {
			return
		}

		out, err = y7s.AddMap(out, resource, aux)
		if err != nil {
			return
		}
	}

	return
}

func (e YamlEncoder) encodeResourceTranslationsKv(p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	out, _ = y7s.MakeMap()

	for _, n := range nodes {
		rt := n.Resource.(*types.ResourceTranslation)
		out, err = y7s.AddMap(out, rt.K, rt.Message)
		if err != nil {
			return
		}
	}

	return
}

func (e YamlEncoder) encodeRbacRules(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	err = e.resolveRulePathDeps(ctx, tt, nodes)
	if err != nil {
		return
	}

	// Batch by access (allow, deny)
	allows := make([]rbacRuleWrap, 0, len(nodes))
	denies := make([]rbacRuleWrap, 0, len(nodes))
	for _, n := range nodes {
		r := n.Resource.(*rbac.Rule)
		if r.Access == rbac.Allow {
			allows = append(allows, rbacRuleWrap{
				node: n,
				rule: r,
			})
		} else {
			denies = append(denies, rbacRuleWrap{
				node: n,
				rule: r,
			})
		}
	}

	// Encode allows
	var aux *yaml.Node
	aux, err = e.encodeRbacRulesByRole(p, allows, tt)
	if err != nil {
		return
	}
	out = base
	out, err = y7s.AddMap(out, "allow", aux)
	if err != nil {
		return
	}

	// Encode denies
	aux, err = e.encodeRbacRulesByRole(p, denies, tt)
	if err != nil {
		return
	}
	out, err = y7s.AddMap(out, "deny", aux)
	if err != nil {
		return
	}

	return
}

func (e YamlEncoder) encodeSettingValues(ctx context.Context, base *yaml.Node, p envoyx.EncodeParams, nodes envoyx.NodeSet, tt envoyx.Traverser) (out *yaml.Node, err error) {
	// Setting values don't have any refs

	for _, n := range nodes {
		sv := n.Resource.(*types.SettingValue)
		base, err = y7s.AddMap(base, sv.Name, sv.Value)
		if err != nil {
			return
		}
	}

	return y7s.MakeMap("settings", base)
}

func (e YamlEncoder) encodeRbacRulesByRole(p envoyx.EncodeParams, rules []rbacRuleWrap, tt envoyx.Traverser) (out *yaml.Node, err error) {
	// Batch by role; make sure to keep the multiple identifier thing in mind
	byRole := make(map[string]rbacRuleRoleWrap)
	for _, r := range rules {
		role := r.node.References["RoleID"]
		aux := rbacRuleRoleWrap{}
		ok := false

		// Check if an identifier already exists
		for _, ident := range role.Identifiers.Slice {
			aux, ok = byRole[ident]
			if ok {
				break
			}
		}

		// If not, create a new entry and resolve the role
		if !ok {
			aux.roleNode = tt.ParentForRef(r.node, role)
			if aux.roleNode != nil {
				aux.role = aux.roleNode.Resource.(*types.Role)
			}
		}

		// Add the rule to the batch and update potentially missing identifiers
		aux.rules = append(aux.rules, r)
		for _, ident := range role.Identifiers.Slice {
			byRole[ident] = aux
		}
	}

	// Encode
	var aux *yaml.Node
	for _, r := range byRole {
		aux, err = e.encodeRbacRulesByResource(p, r.rules)
		if err != nil {
			return
		}

		out, err = y7s.AddMap(out,
			r.roleNode.Identifiers.FriendlyIdentifier(), aux,
		)
		if err != nil {
			return
		}
	}

	return
}

func (e YamlEncoder) encodeRbacRulesByResource(p envoyx.EncodeParams, rules []rbacRuleWrap) (out *yaml.Node, err error) {
	byResource := make(map[string][]rbacRuleWrap)

	// Batch
	for _, r := range rules {
		byResource[r.rule.Resource] = append(byResource[r.rule.Resource], r)
	}

	// Encode
	var aux *yaml.Node
	for resource, rules := range byResource {
		aux, err = e.encodeRbacRuleOperations(p, rules)
		if err != nil {
			return
		}

		out, err = y7s.AddMap(out, resource, aux)
		if err != nil {
			return
		}
	}

	return
}

func (e YamlEncoder) encodeRbacRuleOperations(p envoyx.EncodeParams, rules []rbacRuleWrap) (out *yaml.Node, err error) {
	ops, _ := y7s.MakeSeq()

	for _, r := range rules {
		ops, err = y7s.AddSeq(ops, r.rule.Operation)
		if err != nil {
			return
		}
	}

	return ops, nil
}

func (e YamlEncoder) resolveRulePathDeps(ctx context.Context, tt envoyx.Traverser, nodes envoyx.NodeSet) (err error) {
	for _, n := range nodes {

		var auxIdent string
		for fieldLabel, ref := range n.References {
			// Only resolve path refs; others are done later on
			if !strings.Contains(fieldLabel, "Path.") {
				continue
			}

			auxIdent = safeParentIdentifier(tt, n, ref)
			if auxIdent == "" || auxIdent == "0" {
				continue
			}

			err = n.Resource.SetValue(fieldLabel, 0, auxIdent)
			if err != nil {
				return
			}
		}
	}

	return
}

func (e YamlEncoder) encodeAuthClientSecurityC(ctx context.Context, p envoyx.EncodeParams, tt envoyx.Traverser, n *envoyx.Node, ac *types.AuthClient, sec *types.AuthClientSecurity) (_ any, err error) {
	sqPermittedRoles, err := e.encodeRoleSlice(n, tt, "Security.PermittedRoles", sec.PermittedRoles)
	if err != nil {
		return
	}

	sqProhibitedRoles, err := e.encodeRoleSlice(n, tt, "Security.ProhibitedRoles", sec.ProhibitedRoles)
	if err != nil {
		return
	}

	sqForcedRoles, err := e.encodeRoleSlice(n, tt, "Security.ForcedRoles", sec.ForcedRoles)
	if err != nil {
		return
	}

	var impersonateUser string
	if _, ok := n.References["Security.ImpersonateUser.UserID"]; ok {
		usrRef := n.References["Security.ImpersonateUser.UserID"]
		node := tt.ParentForRef(n, usrRef)
		if node == nil {
			err = fmt.Errorf("invalid user reference %v: user does not exist", usrRef)
			return
		}
		impersonateUser = n.Identifiers.FriendlyIdentifier()
	}

	return y7s.MakeMap(
		"impersonateUser", impersonateUser,
		"permittedRoles", sqPermittedRoles,
		"prohibitedRoles", sqProhibitedRoles,
		"forcedRoles", sqForcedRoles,
	)
}

func (e YamlEncoder) encodeRoleSlice(n *envoyx.Node, tt envoyx.Traverser, k string, rr []string) (out *yaml.Node, err error) {
	sq, _ := y7s.MakeSeq()

	for i := range rr {
		roleRef := n.References[fmt.Sprintf("%s.%d.RoleID", k, i)]
		node := tt.ParentForRef(n, roleRef)
		if node == nil {
			err = fmt.Errorf("invalid user reference %v: user does not exist", roleRef)
			return
		}

		sq, err = y7s.AddSeq(sq, node.Identifiers.FriendlyIdentifier())
		if err != nil {
			return
		}
	}

	return sq, nil
}
