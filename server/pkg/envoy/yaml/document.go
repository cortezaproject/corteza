package yaml

import (
	"context"

	"github.com/cortezaproject/corteza/server/pkg/envoy"
	"github.com/cortezaproject/corteza/server/pkg/envoy/resource"
	"github.com/cortezaproject/corteza/server/pkg/y7s"
	"gopkg.in/yaml.v3"
)

type (
	// Document defines the supported yaml structure
	Document struct {
		compose    *compose
		automation *automation

		roles        roleSet
		users        userSet
		templates    templateSet
		applications applicationSet

		apiGateway apiGatewaySet
		reports    reportSet

		settings settingSet
		rbac     rbacRuleSet
		locale   resourceTranslationSet

		cfg *EncoderConfig
	}
)

func (doc *Document) UnmarshalYAML(n *yaml.Node) (err error) {
	if err = n.Decode(&doc.compose); err != nil {
		return
	}

	if err = n.Decode(&doc.automation); err != nil {
		return
	}

	if doc.rbac, err = decodeRbac(n); err != nil {
		return
	}

	if doc.locale, err = decodeLocale(n); err != nil {
		return
	}

	return y7s.EachMap(n, func(k, v *yaml.Node) error {
		switch k.Value {
		case "roles":
			return v.Decode(&doc.roles)

		case "users":
			return v.Decode(&doc.users)

		case "templates":
			return v.Decode(&doc.templates)

		case "apigateway", "apigw":
			return v.Decode(&doc.apiGateway)

		case "reports":
			return v.Decode(&doc.reports)

		case "applications":
			return v.Decode(&doc.applications)

		case "settings":
			return v.Decode(&doc.settings)

		}

		return nil
	})
}

func (doc *Document) MarshalYAML() (interface{}, error) {
	var err error
	dn, _ := makeMap()

	if doc.compose != nil {
		doc.compose.EncoderConfig = doc.cfg

		cn, err := encodeNode(doc.compose)
		if err != nil {
			return nil, err
		}

		dn, _ = inlineContent(dn, cn)
	}

	if doc.automation != nil {
		doc.automation.EncoderConfig = doc.cfg

		cn, err := encodeNode(doc.automation)
		if err != nil {
			return nil, err
		}

		dn, _ = inlineContent(dn, cn)
	}

	if doc.roles != nil {
		doc.roles.configureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "roles", doc.roles, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.users != nil && len(doc.users) > 0 {
		doc.users.configureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "users", doc.users, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.templates != nil && len(doc.templates) > 0 {
		doc.templates.configureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "templates", doc.templates, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.reports != nil && len(doc.reports) > 0 {
		doc.reports.configureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "reports", doc.reports, doc.cfg.MappedOutput, "handle")
		if err != nil {
			return nil, err
		}
	}

	if doc.apiGateway != nil && len(doc.apiGateway) > 0 {
		doc.apiGateway.configureEncoder(doc.cfg)

		// API GW don't support map representation
		// @todo use path+proto?
		dn, err = addMap(dn,
			"apigateway", doc.apiGateway,
		)
		if err != nil {
			return nil, err
		}
	}

	if doc.applications != nil && len(doc.applications) > 0 {
		doc.applications.configureEncoder(doc.cfg)

		// Applications don't support map representation
		dn, err = addMap(dn,
			"applications", doc.applications,
		)
		if err != nil {
			return nil, err
		}
	}

	if doc.settings != nil {
		dn, err = encodeResource(dn, "settings", doc.settings, doc.cfg.MappedOutput, "name")
		if err != nil {
			return nil, err
		}
	}

	if doc.locale != nil {
		doc.templates.configureEncoder(doc.cfg)

		dn, err = encodeResource(dn, "locale", doc.locale, doc.cfg.MappedOutput, "lang")
		if err != nil {
			return nil, err
		}
	}

	if doc.rbac != nil && len(doc.rbac) > 0 {
		// rbac doesn't support map representation
		m, err := encodeNode(doc.rbac)
		if err != nil {
			return nil, err
		}

		dn, err = inlineContent(dn, m)
		if err != nil {
			return nil, err
		}
	}

	return dn, nil
}

func (doc *Document) Decode(ctx context.Context) ([]resource.Interface, error) {
	nn := make([]resource.Interface, 0, 100)

	mm := make([]envoy.Marshaller, 0, 20)
	if doc.compose != nil {
		mm = append(mm, doc.compose)
	}
	if doc.roles != nil {
		mm = append(mm, doc.roles)
	}
	if doc.users != nil {
		mm = append(mm, doc.users)
	}
	if doc.templates != nil {
		mm = append(mm, doc.templates)
	}
	if doc.automation != nil {
		mm = append(mm, doc.automation)
	}
	if doc.applications != nil {
		mm = append(mm, doc.applications)
	}
	if doc.reports != nil {
		mm = append(mm, doc.reports)
	}
	if doc.apiGateway != nil {
		mm = append(mm, doc.apiGateway)
	}
	if doc.settings != nil {
		for _, s := range doc.settings {
			mm = append(mm, s)
		}
	}
	if doc.locale != nil {
		mm = append(mm, doc.locale)
	}
	if doc.rbac != nil {
		mm = append(mm, doc.rbac)
	}

	for _, m := range mm {
		if tmp, err := m.MarshalEnvoy(); err != nil {
			return nil, err
		} else {
			nn = append(nn, tmp...)
		}
	}

	return nn, nil
}

func (doc *Document) addComposeNamespace(n *composeNamespace) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		doc.compose.Namespaces = make(composeNamespaceSet, 0, 3)
	}

	doc.compose.Namespaces = append(doc.compose.Namespaces, n)
}

func (doc *Document) addComposeModule(m *composeModule) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Modules == nil {
		doc.compose.Modules = make(composeModuleSet, 0)
	}

	doc.compose.Modules = append(doc.compose.Modules, m)
}

func (doc *Document) addComposeRecord(r *composeRecord) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Records == nil {
		doc.compose.Records = make(composeRecordSet, 0)
	}

	doc.compose.Records = append(doc.compose.Records, r)
}

func (doc *Document) addComposePage(p *composePage) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Pages == nil {
		doc.compose.Pages = make(composePageSet, 0)
	}

	doc.compose.Pages = append(doc.compose.Pages, p)
}

func (doc *Document) addComposeChart(c *composeChart) {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Charts == nil {
		doc.compose.Charts = make(composeChartSet, 0)
	}

	doc.compose.Charts = append(doc.compose.Charts, c)
}

func (doc *Document) addAutomationWorkflow(m *automationWorkflow) {
	if doc.automation == nil {
		doc.automation = &automation{}
	}
	if doc.automation.Workflows == nil {
		doc.automation.Workflows = make(automationWorkflowSet, 0)
	}

	doc.automation.Workflows = append(doc.automation.Workflows, m)
}

func (doc *Document) addApiGateway(a *apiGateway) {
	doc.apiGateway = append(doc.apiGateway, a)
}

func (doc *Document) addReport(a *report) {
	doc.reports = append(doc.reports, a)
}

func (doc *Document) addRole(r *role) {
	if doc.roles == nil {
		doc.roles = make(roleSet, 0, 20)
	}

	doc.roles = append(doc.roles, r)
}

func (doc *Document) addUser(u *user) {
	if doc.users == nil {
		doc.users = make(userSet, 0, 20)
	}

	doc.users = append(doc.users, u)
}

func (doc *Document) addTemplate(u *template) {
	if doc.templates == nil {
		doc.templates = make(templateSet, 0, 20)
	}

	doc.templates = append(doc.templates, u)
}

func (doc *Document) addApplication(a *application) {
	if doc.applications == nil {
		doc.applications = make(applicationSet, 0, 20)
	}

	doc.applications = append(doc.applications, a)
}

func (doc *Document) addSetting(s *setting) {
	if doc.settings == nil {
		doc.settings = make([]*setting, 0, 100)
	}

	doc.settings = append(doc.settings, s)
}

func (doc *Document) addRbacRule(r *rbacRule) {
	if doc.rbac == nil {
		doc.rbac = make(rbacRuleSet, 0, 100)
	}

	doc.rbac = append(doc.rbac, r)
}

func (doc *Document) addResourceTranslation(l *resourceTranslation) {
	if doc.locale == nil {
		doc.locale = make(resourceTranslationSet, 0, 100)
	}

	doc.locale = append(doc.locale, l)
}

func (doc *Document) nestComposeModule(ns string, m *composeModule) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.addComposeModule(ns, m)
}

func (doc *Document) nestComposePage(ns string, p *composePage) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.addComposePage(ns, p)
}

func (doc *Document) nestComposePageChild(parent string, p *composePage) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(parent)
	}

	return doc.compose.Pages.addComposePage(parent, p)
}

func (doc *Document) nestComposeChart(ns string, c *composeChart) error {
	if doc.compose == nil {
		doc.compose = &compose{}
	}
	if doc.compose.Namespaces == nil {
		return composeNamespaceErrNotFound(ns)
	}

	return doc.compose.Namespaces.addComposeChart(ns, c)
}
