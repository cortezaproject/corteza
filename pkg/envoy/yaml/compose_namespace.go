package yaml

import (
	"fmt"
	"strconv"

	atypes "github.com/cortezaproject/corteza-server/automation/types"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
)

type (
	composeNamespace struct {
		res  *types.Namespace `yaml:",inline"`
		ts   *resource.Timestamps
		meta composeNamespaceMeta

		envoyConfig   *resource.EnvoyConfig
		encoderConfig *EncoderConfig

		// all known modules on a namespace
		modules composeModuleSet

		// all known charts on a namespace
		charts composeChartSet

		// all known records on a namespace
		records composeRecordSet

		// all known pages on a namespace
		pages composePageSet

		// module's RBAC rules
		rbac rbacRuleSet

		locale resourceTranslationSet
	}
	composeNamespaceSet []*composeNamespace

	composeNamespaceMeta types.NamespaceMeta
)

func composeNamespaceFromResource(r *resource.ComposeNamespace, cfg *EncoderConfig) *composeNamespace {
	return &composeNamespace{
		res:           r.Res,
		encoderConfig: cfg,
	}
}

// configureEncoder configures the composeNamespace encoding
func (nn composeNamespaceSet) configureEncoder(cfg *EncoderConfig) {
	for _, n := range nn {
		n.encoderConfig = cfg
	}
}

func (nn composeNamespaceSet) findComposeNamespace(ns string) *composeNamespace {
	for _, n := range nn {
		if resource.Check(ns, n.res.ID, n.res.Slug, n.res.Name) {
			return n
		}
	}
	return nil
}

func (nn composeNamespaceSet) addComposeModule(ref string, m *composeModule) error {
	ns := nn.findComposeNamespace(ref)

	if ns == nil {
		return composeNamespaceErrNotFound(ref)
	}
	if ns.modules == nil {
		ns.modules = make(composeModuleSet, 0, 1)
	}

	ns.modules = append(ns.modules, m)
	return nil
}

func (nn composeNamespaceSet) addComposePage(ref string, p *composePage) error {
	ns := nn.findComposeNamespace(ref)

	if ns == nil {
		return composeNamespaceErrNotFound(ref)
	}
	if ns.pages == nil {
		ns.pages = make(composePageSet, 0, 1)
	}

	ns.pages = append(ns.pages, p)
	return nil
}

func (nn composeNamespaceSet) addComposeChart(ref string, c *composeChart) error {
	ns := nn.findComposeNamespace(ref)

	if ns == nil {
		return composeNamespaceErrNotFound(ref)
	}
	if ns.charts == nil {
		ns.charts = make(composeChartSet, 0, 1)
	}

	ns.charts = append(ns.charts, c)
	return nil
}

func (nn composeNamespaceSet) addComposeRecord(ref string, r *composeRecord) error {
	ns := nn.findComposeNamespace(ref)

	if ns == nil {
		return composeNamespaceErrNotFound(ref)
	}
	if ns.records == nil {
		ns.records = make(composeRecordSet, 0, 1)
	}

	ns.records = append(ns.records, r)
	return nil
}

func (nn composeNamespaceSet) addRbacRule(ref string, r *rbacRule) error {
	ns := nn.findComposeNamespace(ref)

	if ns == nil {
		return composeNamespaceErrNotFound(ref)
	}
	if ns.rbac == nil {
		ns.rbac = make(rbacRuleSet, 0, 20)
	}

	ns.rbac = append(ns.rbac, r)
	return nil
}

func (c composeNamespaceMeta) empty() bool {
	if c == (composeNamespaceMeta{}) {
		return true
	}

	return false
}

func relNsToRef(ns *types.Namespace) string {
	return firstOkString(ns.Slug, ns.Name, strconv.FormatUint(ns.ID, 10))
}

func relModToRef(mod *types.Module) string {
	return firstOkString(mod.Handle, mod.Name, strconv.FormatUint(mod.ID, 10))
}

func relWfToRef(mod *atypes.Workflow) string {
	name := ""
	if mod.Meta != nil {
		name = mod.Meta.Name
	}
	return firstOkString(mod.Handle, name, strconv.FormatUint(mod.ID, 10))
}

func composeNamespaceErrNotFound(i string) error {
	return fmt.Errorf("namespace not found: %v", i)
}
